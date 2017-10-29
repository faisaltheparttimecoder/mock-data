package postgres

import (
	"fmt"
	"strings"
	"database/sql"
	"github.com/pivotal/mock-data/core"
)

// Fix Foreign Key
func fixFKey(db *sql.DB, fk constraint, debug bool) error {

	// The objects involved in this foriegn key clause
	fkeyObjects, err := getFKeyObjects(fk)
	if err != nil {
		return fmt.Errorf("Unable to scan the total primary key violators: %v", err)
	}

	// Time to fix the foriegn key issues
	err = UpdateFKViolationRecord(db, fkeyObjects, debug)
	if err != nil {
		return err
	}
	return nil
}

// Get Foriegn key objects
type foreignKey struct {
	table, column, reftable, refcolumn string
}


// Functions to extract the columns names from the provided
// command output.
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
func UpdateFKViolationRecord(db *sql.DB, fkObjects foreignKey, debug bool) error {

	var TotalViolators int = 1

	// Get total number of records on the table
	totalRow, err := totalRows(db, fkObjects.reftable)
	if err != nil {
		return err
	}

	if debug {
		log.Debugf("Checking / Fixing FOREIGN KEY Violation table: %s, column: %s, reference: %s(%s)", fkObjects.table, fkObjects.table, fkObjects.reftable, fkObjects.refcolumn)

	}

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
