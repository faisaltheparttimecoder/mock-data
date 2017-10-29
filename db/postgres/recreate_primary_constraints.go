package postgres

import (
	"strings"
	"fmt"
	"database/sql"
	"github.com/pivotal/mock-data/core"
)

// Fix Primary Key
func fixPKey(db *sql.DB, pk constraint, fixtype string, debug bool) error {

	// Start with 1 violation to begin the loop
	var TotalViolators int = 1

	// Extract the columns from the list that was collected during backup
	keys, err := core.ColExtractor(pk.column, `\(.*?\)`)
	if err != nil {
		return fmt.Errorf("Unable to determine PK violators key columns: %v", err)
	}
	cols := strings.Trim(keys, "()")

	// If logging is tuned on then paste this message on the screen.
	if debug {
		log.Debugf("Checking / Fixing %s KEY Violation table: %s, column: %s", fixtype, pk.table, cols)
	}

	// Loop till we get a 0 value (i.e 0 violation )
	for TotalViolators > 0 {

		// How many violations are we having, if zero then loop breaks
		TotalViolators, err = getTotalViolation(db, GetTotalPKViolator(pk.table, cols))
		if err != nil {
			return err
		}

		// Perform the below action only if the violators is greater than 0
		if TotalViolators > 0 {

			// If there are two or more columns forming a PK or UK
			// lets only fix column by column.
			totalColumns := strings.Split(cols, ",")

			// Get datatype assosicated with the datatypes
			dtypes, err := obtainDataType(db, GetDatatype(pk.table, totalColumns))
			if err != nil {
				return err
			}

			// Fix the primary constraints by picking the columns from the
			// array, i.e we update the column one by one.
			for _, v := range dtypes {
				column := strings.Split(v, ":")[0]
				dtype := strings.Split(v, ":")[1]
				err = fixPKViolator(db, pk.table, column, dtype)
				if err != nil {
						return err
				}
			}
		}
	}

	return nil
}

// Get datatype of the column that is based on the
// column provided.
func obtainDataType(db *sql.DB, query string) ([]string, error) {
	var colname, dtype string
	var dtypes []string

	rows, err := db.Query(query)
	if err != nil {
		return dtypes, fmt.Errorf("Unable to get the datatype of key violators: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&colname, &dtype)
		if err != nil {
			return dtypes, fmt.Errorf("Unable to scan the datatype of key violators: %v", err)
		}
		dtypes = append(dtypes, colname + ":" + dtype)
	}

	return dtypes, nil
}


// Get total violations of the columns that are part of the PK
func getTotalViolation(db *sql.DB, query string) (int, error) {

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


// Fix Primary Ket string violators.
func fixPKViolator(db *sql.DB, tab, col, dttype string) error {

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

		// Update the column
		var duplicateRow string
		err = rows.Scan(&duplicateRow)
		if err != nil {
			return fmt.Errorf("Error in scanning rows of PK violators: %v", err)
		}
		_, err = db.Exec(UpdatePKey(tab, col, duplicateRow, fmt.Sprintf("%v", newdata)))
		if err != nil {
			return fmt.Errorf("Error in fixing the rows of PK violators: %v", err)
		}
	}

	return nil
}

