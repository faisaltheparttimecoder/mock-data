package main

import (
	"fmt"
	"strings"
)

type constraint struct {
	table, column string
}

var (
	savedConstraints = map[string][]constraint{"PRIMARY": {}, "CHECK": {}, "UNIQUE": {}, "FOREIGN": {}}
	constraints      = []string{"p", "f", "u", "c"}
)

// Backup DDL of objects which are going to drop to
// allow faster and smooth transition of inputting data.
func BackupDDL() {
	// Create the necessary directory of this run
	Infof("Saving all the backup files to the path: %s", Path)
	CreateDirectory()

	// Backup Constraints and Indexes
	backupConstraints()
	backupIndexes()
}

// Backup all the constraints
func backupConstraints() {
	Debugf("Backing up all the constraints from the database: %s", cmdOptions.Database)
	for _, constr := range constraints {
		filename := fmt.Sprintf("%s/%s_constraint_backup_%s.sql", Path, programName, constr)
		constraintInfo := GetPGConstraintDDL(constr)
		for _, c := range constraintInfo {
			// DDL
			constraintDDL := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;\n",
				c.Tablename, c.Constraintname, c.Constraintkey)

			// Before dropping the constraints ensure we have the information
			// saved about the state of what constraints this table had
			ctype, err := constraintFinder(constr)
			if err != nil {
				Fatalf("Error during backup of constraints: %v", err)
			}
			saveConstraints(filename, constraintDDL, ctype, c.Tablename, c.Constraintkey)
		}
	}
}

// Backup all the unique index
func backupIndexes() {
	Debugf("Backing up all the unique indexes from the database: %s", cmdOptions.Database)
	filename := fmt.Sprintf("%s/mockd_constraint_backup_u.sql", Path)
	indexes := GetPGIndexDDL()
	for _, i := range indexes {
		indexDDL := fmt.Sprintf("%s;\n", i.Indexdef)
		saveConstraints(filename, indexDDL, "UNIQUE", i.Tablename, i.Indexdef)
	}
}

// Save constraints on a file and on buffer
func saveConstraints(filename, message, ctype, table, column string) {
	err := WriteToFile(filename, message)
	if err != nil {
		Fatalf("Error in saving the index DDL to the file: %v", err)
	}
	savedConstraints[ctype] = append(savedConstraints[ctype],
		constraint{table: table, column: column})
}

// What is the type of constraints is it
func constraintFinder(contype string) (string, error) {
	switch {
	// Check constraint
	case strings.Contains(contype, "c"):
		return "CHECK", nil
		// Primary constraint
	case strings.Contains(contype, "p"):
		return "PRIMARY", nil
		// Foreign constraint
	case strings.Contains(contype, "f"):
		return "FOREIGN", nil
		// Unique constraint
	case strings.Contains(contype, "u"):
		return "UNIQUE", nil
	default:
		return "", fmt.Errorf("cannot understand the type of constraints")
	}
	return "", nil
}

// Remove the constraints before loading to ease the pain of any
// failure due to constraint errors
func RemoveConstraints(table string) {
	Debugf("Removing constraints for table: %s", table)
	var statement string

	// Obtain all the constraints on the table that we going to load data to
	constraints := GetConstraintsPertab(table)

	// scan through the rows and generate the drop command
	for _, c := range constraints {
		// Generate the DROP DDL command
		if c.Constrainttype == "index" { // if the constraint is a index
			statement = fmt.Sprintf("DROP INDEX %s CASCADE;", c.Constraintname)
		} else { // if the constraint is a constraint
			statement = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s CASCADE;", table, c.Constraintname)
		}

		// Execute the statement
		_, err := ExecuteDB(statement)
		if err != nil {
			// Ignore does not exist error eg.s the primary key is dropped
			// then the index also goes along with it , so no need to panic here
			errMsg := fmt.Sprintf("%s", err)
			failureMsg := fmt.Sprintf("Encountered error when removing constraints for table %s, err: %v", table, err)
			IgnoreError(errMsg, "does not exist", failureMsg)
		}
	}
}
