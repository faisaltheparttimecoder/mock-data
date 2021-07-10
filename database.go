package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// Execute cleanup of demo database before loading if needed
func executeDemoDatabasePreCleanup() {
	// NOTE: This is a hidden feature used mainly by test cases, not used widely
	// so use this with caution, since it will cleanup all the tables in the said database
	shouldWeDrop := viper.GetBool("MOCK_DATA_TEST_RUNNER")
	if shouldWeDrop {
		Infof("Executing the cleanup of all objects in the database: %s", cmdOptions.Database)
		_, err := ExecuteDB(dropDemoDatabase())
		if err != nil {
			errMsg := fmt.Sprintf("%s", err)
			failureMsg := fmt.Sprintf("Failure in pre cleanup a demo tables in the database %s, err: %v",
				cmdOptions.Database, err)
			IgnoreError(errMsg, "does not exist", failureMsg)
		}
	}
}

// Execute demo database
func ExecuteDemoDatabase() {
	Infof("Create demo tables in the database: %s", cmdOptions.Database)

	// pre cleanup script
	executeDemoDatabasePreCleanup()

	// Execute the demo database dump
	var err error
	if GreenplumOrPostgres == "postgres" {
		_, err = ExecuteDB(demoDatabasePostgres())
	} else {
		_, err = ExecuteDB(demoDatabaseGreenplum())
	}
	if err != nil {
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
