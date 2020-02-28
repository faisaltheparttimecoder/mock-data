package main

import (
	"fmt"
	"strings"
)

var (
	ignoreErr = []string{
		"ERROR #42P16 multiple primary keys for table",
		"already exists"}
)

// Get Foriegn key objects
type ForeignKey struct {
	Table, Column, Reftable, Refcolumn string
}

// Ty to recreate all the constraints where ever we can
func FixConstraints() {
	//Fix the constraints in this order
	//var constr = []string{"PRIMARY", "UNIQUE", "CHECK", "FOREIGN"}
	var constr = []string{"PRIMARY", "UNIQUE", "FOREIGN"}
	for _, v := range constr {
		totalViolations := len(savedConstraints[v])
		Infof("Found %v violation of %s KEYS, attempting to fix them", totalViolations, v)
		bar := StartProgressBar(fmt.Sprintf("Fixing %s KEYS violation", v), totalViolations)
		for _, con := range savedConstraints[v] {
			switch {
			case v == "PRIMARY":
				fixPKey(con)
			case v == "UNIQUE": // Run the same logic as primary key
				fixPKey(con)
				//case v == "CHECK": // TODO: Its hard to predict the check constraint ATM
				//  fixCheck(db, con)
			case v == "FOREIGN":
				fixFKey(con)
			}
			bar.Add(1)
		}
		fmt.Println()
	}

	// Recreate constraints
	recreateAllConstraints()
}

// Fix the primary key
func fixPKey(pk constraint) {
	Debugf("Fixing the Primary / Unique Key")
	totalViolators := 1
	// Extract the columns from the list that was collected during backup
	keys, err := ColExtractor(pk.column, `\(.*?\)`)
	if err != nil {
		Fatalf("unable to determine PK violators key columns: %v", err)
	}
	cols := strings.Trim(keys, "()")

	for totalViolators > 0 { // Loop till we get a 0 value (i.e 0 violation )
		// How many violations are we having, if zero then loop breaks
		totalViolators = getTotalPKViolator(pk.table, cols)
		if totalViolators > 0 { // Found violation, time to fix it

			// If there are two or more columns forming a PK or UK
			// lets only fix column by column.
			totalColumns := strings.Split(cols, ",")

			// Get data type associated with the data types
			dTypes := getDatatype(pk.table, totalColumns)

			//Fix the primary constraints by picking the columns from the
			//array, i.e we update the column one by one.
			for _, v := range dTypes {
				fixPKViolator(pk.table, v.Colname, v.Dtype)
			}
		}
	}
}

// Fix Primary Key string violators.
func fixPKViolator(tab, col, dttype string) {
	// Get all the strings that violates the primary key constraints
	pkViolators := GetPKViolators(tab, col)

	for _, v := range pkViolators {
		// Get a new random data based on data type
		newdata, err := BuildData(dttype)
		if err != nil {
			Fatalf("Error when generating new data for PK Violation: %v", err)
		}

		// Replace the old data with new data
		UpdatePKKey(tab, col, v.Row, fmt.Sprintf("%v", newdata))
	}
}

// Fix the Foreign Keys
func fixFKey(con constraint) {
	Debugf("Fixing the Primary / Unique Key")
	totalViolators := 1

	// The objects involved in this foriegn key clause
	fkeyObjects := getForeignKeyColumns(con)

	// Time to fix the foreign key issues
	// Get total number of records on the table
	totalRow := TotalRows(fkeyObjects.Reftable)

	Debugf("Checking / Fixing FOREIGN KEY Violation table: %s, column: %s, reference: %s(%s)",
		fkeyObjects.Table, fkeyObjects.Column, fkeyObjects.Reftable, fkeyObjects.Refcolumn)

	// Loop till we reach the the end of the loop
	for totalViolators > 0 {

		// Total foreign key violators
		totalViolators = GetTotalFKViolators(*fkeyObjects)

		// Run only if there is a violations else no
		if totalViolators > 0 {
			violators := GetFKViolators(*fkeyObjects)
			for _, v := range violators {
				UpdateFKeys(*fkeyObjects, totalRow, v.Row)
			}
		}
	}

}

// Get Foreign Keys column and reference column
func getForeignKeyColumns(con constraint) *ForeignKey {
	// Extract reference clause from the value
	refClause, err := ColExtractor(con.column, `REFERENCES[ \t]*([^\n\r]*\))`)
	if err != nil {
		Fatalf("Unable to extract reference key clause from fk clause: %v", err)
	}

	// Extract the fk column from the clause
	fkCol, err := ColExtractor(strings.Replace(con.column, refClause, "", -1), `\(.*?\)`)
	if err != nil {
		Fatalf("Unable to extract foreign key column from fk clause: %v", err)
	}
	fkCol = strings.Trim(fkCol, "()")

	// Extract the reference column from the clause
	refCol, err := ColExtractor(refClause, `\(.*?\)`)
	if err != nil {
		Fatalf("Unable to extract reference key column from fk clause: %v", err)
	}

	// Extract reference table from the clause
	refTab := strings.Replace(refClause, refCol, "", -1)
	refTab = strings.Replace(refTab, "REFERENCES ", "", -1)
	refCol = strings.Trim(refCol, "()")

	return &ForeignKey{con.table, fkCol, refTab, refCol}
}

// Ignore Error strings matches
func IgnoreErrorString(errmsg string) bool {
	for _, ignore := range ignoreErr {
		if strings.HasSuffix(errmsg, ignore) || strings.HasPrefix(errmsg, ignore) {
			return true
		}
	}
	return false
}

// Recreate all the constraints of the database ( in case we have dropped any )
func recreateAllConstraints() {
	Infof("Attempting to recreating all the constraints")
	failedConstraintsFile := fmt.Sprintf("%s/failed_constraint_creations.sql", Path)
	var AnyError bool = false

	// list the backup files collected.
	for _, con := range constraints {
		backupFile, err := ListFile(Path, fmt.Sprintf("%s_constraint_backup_%s.sql", programName, con))
		if err != nil {
			Fatalf("Error when listing all the backup files from the directory %s, err: %v", Path, err)
		}

		// run it only if we do find the backup file
		if len(backupFile) > 0 {
			b := backupFile[0]
			contents, err := ReadFile(b)
			if err != nil {
				Fatalf("Error in reading the backup file %s: %v", b, err)
			}

			// Start the progress bar
			bar := StartProgressBar(fmt.Sprintf("Recreated the Constraint Type \"%s\"", con), len(contents))

			// Recreate all constraints one by one, if we can't create it then display the message
			// on the screen and continue with the rest, since we don't want it to fail if we cannot
			// recreate constraint of a single table.
			for _, content := range contents {
				_, err := ExecuteDB(content)
				if err != nil && !IgnoreErrorString(fmt.Sprintf("%s", err)) {
					Debugf("Error creating constraint %s, err: %v", content, err)
					// Try an attempt to recreate constraint again after deleting the
					// violating row
					successOrFailure := deleteViolatingPkOrUkConstriants(content)
					if !successOrFailure { // didn't succeed, ask the user to fix it manually
						err = WriteToFile(failedConstraintsFile, content+"\n")
						if err != nil {
							Fatalf("Error when saving the failed restore to file %s, err %v",
								failedConstraintsFile, err)
						}
						AnyError = true
					}
				}
				bar.Add(1)
			}
			fmt.Println()
		}
	}

	if AnyError {
		Warnf("There have been issue creating few constraints and would need manual cleanup at your end, "+
			"all the constraints that failed has been saved on to file: %s", failedConstraintsFile)
	}
}

// we tried to fix the primary key violation, but due to the nature
// of how we fix the constraints like PK (or UK) followed by FK , there
// are chances that we might inject duplicate keys again, for eg.s if
// there is a PK ( or UK ) on a FK reference table. so the aim here
// is, we will delete the rows that violates it and hoping that it will help in
// recreating the constraints. Yes we will loose that row at least that help to
// recreate constraints ( fingers crossed :) )
func deleteViolatingPkOrUkConstriants(con string) bool {
	Debugf("Attempting to run the constraint command %s second time, after deleting violating rows", con)
	// does the DDL contain PK or UK keyword then do the following
	// rest send them back for user to fix it.
	if isSubStringAvailableOnString(con, "ADD CONSTRAINT.*PRIMARY KEY|ADD CONSTRAINT.*UNIQUE") {
		column, _ := ColExtractor(con, `\(.*?\)`)
		table, _ := ColExtractor(con, `ALTER TABLE(.*)ADD CONSTRAINT`)
		table = strings.Trim(strings.Trim(table, "ALTER TABLE"), "ADD CONSTRAINT")
		column = strings.Trim(column, "()")
		err := deleteViolatingConstraintKeys(table, column)
		if err != nil { // we failed to delete the the constraint violation rows
			Debugf("Error when deleting rows from the constraint violation rows: %v", err)
			return false
		}
		_, err = ExecuteDB(con) // retry to create the constraint again
		if err != nil { // we failed to recreate the constraint
			Debugf("Error when 2nd attempt to recreate constraint: %v", err)
			return false
		}
		// successfully cleaned it up
		return true
	}
	return false
}
