package main

import "fmt"

// MockSchema: Extract all the table from schema and start mocking
func MockSchema() {
	Infof("Starting the program to mock all the tables under the schema %s in the database: %s",
		cmdOptions.SchemaName, cmdOptions.Database)

	// Extract the table
	var tables []DBTables
	whereClause := fmt.Sprintf("AND n.nspname = '%s'", cmdOptions.SchemaName)
	if GreenplumOrPostgres == "postgres" { // Use postgres specific query
		tables = allTablesPostgres(whereClause)
	} else { // Greenplum flavor postgres database
		tables = allTablesGPDB(whereClause)
	}

	// Start the mocking process
	MockTable(tables)
}
