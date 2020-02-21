package main

import (
	"fmt"
)

// Execute demo database
func ExecuteDemoDatabase() {
	Infof("Create demo tables in the database: %s", cmdOptions.Database)

	// Execute the demo database dump
	_, err := ExecuteDB(demoDatabase())
	if err != nil {
		Fatalf("%v", err)
		errMsg := fmt.Sprintf("%s", err)
		failureMsg := fmt.Sprintf("Failure in creating a demo tables in the database %s, err: %v",
			cmdOptions.Database, err)
		IgnoreError(errMsg, "does not exist", failureMsg)
	}

	Infof("Completed creating demo tables in the database: %s", cmdOptions.Database)
}

// Mock the whole database
func MockDatabase() {
	// Get the table list that we have to mock the data
	Infof("Starting the program to mock full database")
	tableList := dbExtractTables("")
	MockTable(tableList)
}

// Extract all the tables in the database
func dbExtractTables(whereClause string) []DBTables {
	Infof("Extracting the tables in the database: %s", cmdOptions.Database)

	// Obtain all the tables in the database
	var tables []DBTables
	if GreenplumOrPostgres == "postgres" { // Use postgres specific query
		tables = allTablesPostgres(whereClause)
	} else { // Greenplum flavor postgres database
		tables = allTablesGPDB(whereClause)
	}

	return tables
}
