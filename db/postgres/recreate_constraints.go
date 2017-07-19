package postgres

import (
	"fmt"
	"database/sql"
	"github.com/ielizaga/mockd/core"
	"strings"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("mockd")
)

// fix the data loaded so that we can reenable the constraints
func FixConstraints(db *sql.DB, timestamp string) error {

	// Fix the constraints in this order
	//var constr = []string{"PRIMARY", "UNIQUE", "CHECK", "FOREIGN"}
	var constr = []string{"PRIMARY", "UNIQUE", "FOREIGN"}
	for _, v := range constr {
		log.Infof("Checking for any %s KEYS, fixing them if there is any violations", v)
		for _, con := range savedConstraints[v] {
			switch  {
			case v == "PRIMARY":
				err := fixPKey(db, con, v)
				if err != nil {
					return err
				}
			case v == "UNIQUE": // Run the same logic as primary key
				err := fixPKey(db, con, v)
				if err != nil {
					return err
				}
			case v == "CHECK":
				err := fixCheck(db, con)
				if err != nil {
					return err
				}
			case v == "FOREIGN":
				err := fixFKey(db, con)
				if err != nil {
					return err
				}
			}
		}
	}

	// Recreate constraints
	err := recreateAllConstraints(db, timestamp)
	if err != nil {
		return err
	}

	return nil
}

// Fix Primary Key
func fixPKey(db *sql.DB, pk constraint, fixtype string) error {

	var TotalViolators int = 1

	// Extract the columns from the list
	keys, err := core.ColExtractor(pk.column, `\(.*?\)`)
	if err != nil {
		return fmt.Errorf("Unable to determine PK violators key columns: %v", err)
	}
	cols := strings.Trim(keys, "()")

	log.Debugf("Checking / Fixing %s KEY Violation table: %s, column: %s", fixtype, pk.table, cols)

	// Loop till we get a 0 value
	for TotalViolators > 0 {

		// How many violations are we having
		TotalViolators, err = getTotalViolation(db, GetTotalPKViolator(pk.table, cols))
		if err != nil {
			return err
		}

		// Perform the below action only if the violators is greater than 0
		if TotalViolators > 0 {

			// If there are two or more columns forming a PK or UK
			// lets only fix one column, not need to fix all of it
			totalColumns := strings.Split(cols, ",")
			if len(totalColumns) > 1 {
				cols = totalColumns[0]
			}

			// Get datatype
			dtype, err := obtainDataType(db, GetDatatype(pk.table, cols))
			if err != nil {
				return err
			}

			// Fix the constraints ( on test only primary key is violated by int )
			testdata, err := core.BuildData(dtype)
			if err != nil {
				return err
			}

			// Take action based on datatype
			if core.IsIntorString(fmt.Sprintf("%v", testdata)) { // Integer
				_, err := db.Exec(UpdateIntPKey(pk.table, cols, dtype))
				if err != nil {
					return fmt.Errorf("Error in cleaning up primary key violators (int): %v", err)
				}
			} else { // a string ( most of the string are unique, anyways confirm it)
				err = fixPKStringViolator(db, pk.table, cols, dtype)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

// Get total violations
func getTotalViolation(db *sql.DB, query string) (int, error ) {

	var TotalViolators int

	// Check if this table is violating any primary constraints
	rows, err := db.Query(query)
	if err != nil {
		return TotalViolators, fmt.Errorf("Unable to get the total key violators: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&TotalViolators)
		if err != nil {
			return TotalViolators, fmt.Errorf("Unable to scan the total key violators: %v", err)
		}

	}

	return TotalViolators, nil
}

// Get datatype of the column
func obtainDataType(db *sql.DB, query string) (string, error) {
	var dtype string

	rows, err := db.Query(query)
	if err != nil {
		return dtype, fmt.Errorf("Unable to get the datatype of key violators: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&dtype)
		if err != nil {
			return dtype, fmt.Errorf("Unable to scan the datatype of key violators: %v", err)
		}
	}

	return dtype, nil
}

// Fix Primary Ket string violators.
func fixPKStringViolator(db *sql.DB, tab, col, dttype string) error {

	// Get all the strings that violates the primary key constraints
	rows, err := db.Query(GetPKViolator(tab, col))
	if err != nil {
		return fmt.Errorf("Error in getting rows of PK violators: %v", err)
	}
	for rows.Next() {

		// Get a random data based on datatype
		newdata, err := core.BuildData(dttype)
		if err != nil {
			return err
		}

		// Update the string
		var duplicateRow string
		err = rows.Scan(&duplicateRow)
		if err != nil {
			return fmt.Errorf("Error in scanning rows of PK violators: %v", err)
		}
		_, err = db.Exec(UpdateStringPKey(tab, col, duplicateRow, fmt.Sprintf("%v", newdata)))
		if err != nil {
			return fmt.Errorf("Error in fixing the rows of PK violators: %v", err)
		}
	}

	return nil
}

// fix Check constraints
func fixCheck(db *sql.DB, ck constraint) error {

	var TotalViolators int

	// Extract the key columns
	keys, err := core.ColExtractor(ck.column, `\(.*?\)`)
	if err != nil {
		return fmt.Errorf("Unable to determine CK violators key columns: %v", err)
	}
	cols := strings.Trim(keys, "()")

	// Extract the column name from it
	colname, err := core.ColExtractor(cols, `^(.*?)\s`)
	if err != nil {
		return fmt.Errorf("Unable to determine column name from keys: %v", err)
	}

	// Check if this table is violating any check constraints
	rows, err := db.Query(GetTotalCKViolator(ck.table, colname, cols))
	if err != nil {
		return fmt.Errorf("Unable to get the total primary key violators: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&TotalViolators)
		if err != nil {
			return fmt.Errorf("Unable to scan the total primary key violators: %v", err)
		}

	}

	log.Info(TotalViolators)

	return nil
}

// Fix Foreign Key
func fixFKey(db *sql.DB, fk constraint) error {

	// The objects involved in this foriegn key clause
	fkeyObjects, err := getFKeyObjects(fk)
	if err != nil {
		return fmt.Errorf("Unable to scan the total primary key violators: %v", err)
	}

	// Time to fix the foriegn key issues
	err = UpdateFKViolationRecord(db, fkeyObjects)
	if err != nil {
		return err
	}
	return nil
}

// Get Foriegn key objects
type foreignKey struct {
	table, column, reftable, refcolumn string
}
func getFKeyObjects(fk constraint) (foreignKey, error) {

	var foriegnClause foreignKey

	// Extract reference clause from the value
	refClause, err := core.ColExtractor(fk.column, `REFERENCES[ \t]*([^\n\r]*\))`)
	if err != nil {
		return foriegnClause, fmt.Errorf("Unable to extract reference key clause from fk clause: %v", err)
	}

	// Extract the fk column from the clause
	fkCol, err := core.ColExtractor(strings.Replace(fk.column, refClause, "", -1), `\(.*?\)`)
	if err != nil {
		return foriegnClause, fmt.Errorf("Unable to extract foreign key coulmn from fk clause: %v", err)
	}
	fkCol = strings.Trim(fkCol, "()")

	// Extract the reference column from the clause
	refCol, err := core.ColExtractor(refClause, `\(.*?\)`)
	if err != nil {
		return foriegnClause, fmt.Errorf("Unable to extract reference key coulmn from fk clause: %v", err)
	}

	// Extract reference table from the clause
	refTab := strings.Replace(refClause, refCol, "", -1)
	refTab = strings.Replace(refTab, "REFERENCES ", "", -1)
	refCol = strings.Trim(refCol, "()")

	foriegnClause = foreignKey{fk.table, fkCol, refTab, refCol}

	return foriegnClause, nil
}
// Update the foriegn key violation tables.
func UpdateFKViolationRecord(db *sql.DB, fkObjects foreignKey) error {

	var TotalViolators int = 1

	// Get total number of records on the table
	totalRow, err := totalRows(db, fkObjects.reftable)
	if err != nil {
		return err
	}

	log.Debugf("Checking / Fixing FOREIGN KEY Violation table: %s, column: %s, reference: %s(%s)", fkObjects.table, fkObjects.table, fkObjects.reftable, fkObjects.refcolumn)

	// Loop till we reach the the end of the loop
	for TotalViolators > 0 {

		// Total foreign key violators
		TotalViolators, err = totalFKViolators(db, fkObjects)
		if err != nil {
			return err
		}

		// Run only if there is a violations else no
		if TotalViolators > 0 {
			err := updateFKViolators(db, fkObjects, totalRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Total rows of the table
func totalRows(db *sql.DB, table string) (string, error) {

	var TotalRow string

	// Total rows in the reference table
	rows, err := db.Query(TotalRows(table))
	if err != nil {
		return TotalRow, fmt.Errorf("Error in getting total rows from the table: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&TotalRow)
		if err != nil {
			return TotalRow, fmt.Errorf("Error in scanning total rows from the table: %v", err)
		}
	}

	return TotalRow, nil
}

// Total Foreign Key violators
func totalFKViolators(db *sql.DB, fkObjects foreignKey) (int, error) {

	var TotalViolators int

	// Get the total rows that are violating the FK constraint
	rows, err := db.Query(GetTotalFKViolators(fkObjects.table, fkObjects.column, fkObjects.reftable, fkObjects.refcolumn))
	if err != nil {
		return TotalViolators, fmt.Errorf("Error in getting total violator of foriegn keys: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&TotalViolators)
		if err != nil {
			return TotalViolators, fmt.Errorf("Error in scanning total violator of foriegn keys: %v", err)
		}
	}

	return TotalViolators, nil
}

// Update foreign key violators
func updateFKViolators(db *sql.DB, fkObjects foreignKey, totalRows string) error {

	var violatorKey string

	// Get all the rows that are violating the FK constraint
	rows, err := db.Query(GetFKViolators(fkObjects.table, fkObjects.column, fkObjects.reftable, fkObjects.refcolumn))
	if err != nil {
		return fmt.Errorf("Error in retreving total violator values of foriegn keys: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&violatorKey)
		if err != nil {
			return fmt.Errorf("Error in scanning total violator of foriegn keys: %v", err)
		}
		_, err := db.Exec(UpdateFKeys(fkObjects.table, fkObjects.column, fkObjects.reftable, fkObjects.refcolumn, violatorKey, totalRows))
		if err != nil {
			return fmt.Errorf("Error in update violator of foriegn keys: %v", err)
		}
	}

	return nil
}

// Recreate all the constraints of the database ( in case we have dropped any )
func recreateAllConstraints(db *sql.DB,timestamp string) error {

	// list the backup files collected.
	for _, con := range constraints {
		backupFile, err := core.ListFile(".", "*_"+con+"_"+timestamp+".sql")
		if err != nil {
			return fmt.Errorf("Error in getting the list of backup files: %v", err)
		}

		// run it only if we do find the backup file
		if len(backupFile) > 0 {
			contents, err := core.ReadFile(backupFile[0])
			if err != nil {
				return fmt.Errorf("Error in reading the backup files: %v", err)
			}

			for _, content := range contents {
				_, err = db.Exec(content)
				if err != nil && !core.IgnoreErrorString(fmt.Sprintf("%s", err), ignoreErr) {
					return fmt.Errorf("Error in recreating constraints: %v", err)
				}
			}
		}
	}

	return nil
}