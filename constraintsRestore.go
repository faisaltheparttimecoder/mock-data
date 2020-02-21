package main

import (
	"fmt"
	"strings"
)

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
			//	fixCheck(db, con)
			//case v == "FOREIGN":
			//	fixFKey(con)
			}
			bar.Add(1)
		}
		fmt.Println()
	}

	//// Recreate constraints
	//failureDetected, err := recreateAllConstraints(db, timestamp, debug)
	//if failureDetected || err != nil {
	//	return err
	//}
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
		for totalViolators > 0 { // Found violation, time to fix it

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

// Fix Primary Ket string violators.
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
