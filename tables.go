package main

import (
	"fmt"
	"strings"
)

// CreateFakeTables creates fakes tables has requested by the user
func CreateFakeTables() {
	Infof("Create %d fake tables on the database: %s", cmdOptions.Tab.TotalTables, cmdOptions.Database)

	// Start the progress bar
	bar := StartProgressBar("Creating Tables", cmdOptions.Tab.TotalTables)

	// Create tables
	for i := 1; i <= cmdOptions.Tab.TotalTables; i++ {
		createTable(i)
		_ = bar.Add(1)
	}

	Infof("Completed creating %d fake tables on the database: %s", cmdOptions.Tab.TotalTables, cmdOptions.Database)
}

// Generate table names
func tableNameGenerator(n int) string {
	return fmt.Sprintf("%s_%s_%d",
		RemoveSpecialCharacters(cmdOptions.Tab.TableNamePrefix), RandomString(6), n)
}

// Generate create table statement
func createTableStatementGenerator(n int) string {
	// Table
	createTableDDL := "CREATE TABLE "
	tableName := tableNameGenerator(n)
	if cmdOptions.Tab.CaseSensitive {
		createTableDDL = createTableDDL + fmt.Sprintf("\"%s\".\"%s\" (", cmdOptions.Tab.SchemaName, tableName)
	} else {
		createTableDDL = createTableDDL + fmt.Sprintf("\"%s\".%s (", cmdOptions.Tab.SchemaName, tableName)
	}

	// Column
	columnName := RemoveSpecialCharacters(cmdOptions.Tab.ColumnNamePrefix)
	howManyColumns := RandomInt(2, cmdOptions.Tab.MaxColumns)

	// DataType
	tableColStatement := ""

	// Column + DataTypes
	for i := 1; i <= howManyColumns; i++ {
		tableColStatement = tableColStatement + fmt.Sprintf("%s_%d %s", columnName,
			i, RandomPickerFromArray(SupportedDataTypes()))
	}

	// do we have a column
	if IsStringEmpty(tableColStatement) {
		Warnf("The table %s has no columns", tableName)
	}

	// Create table statement
	return createTableDDL + fmt.Sprintf("%s);", strings.TrimRight(tableColStatement, ","))
}

// Execute the create table statement generated above in the database
func createTable(n int) {
	createTableDDL := createTableStatementGenerator(n)
	Debugf("Create tables: %s", createTableDDL)

	// Execute the statement
	_, err := ExecuteDB(createTableDDL)
	if err != nil {
		Fatalf("Failure in creating the tables in the database %s, err: %v", cmdOptions.Database, err)
	}
}

// MockTables mocks the provided tables
func MockTables() {
	Infof("Starting mocking of table: %s", cmdOptions.Tab.FakeTablesRows)
	whereClause := generateWhereClause()
	tableList := dbExtractTables(whereClause)
	MockTable(tableList)
}

// Generate the where clause from the table list provided
func generateWhereClause() string {
	Debug("Generating the where condition for the table list")

	// where condition syntax
	whereClause := "AND (n.nspname    || '.' || c.relname) IN (%s)"

	// Loop and generate the where condition
	var w []string
	t := strings.Split(cmdOptions.Tab.FakeTablesRows, ",")
	for _, table := range t {

		// if there is no schema then add in public the default schema
		if !strings.Contains(table, ".") {
			table = fmt.Sprintf("%s.%s", cmdOptions.Tab.SchemaName, table)
		}

		// Separate the table name and schema name
		s := strings.Split(table, ".")

		// If no schema name is found then error out
		if len(s) < 2 {
			Fatalf("Table options should be of the format <schema>.<table>, table option: %s", table)
		}

		// generate the in clause
		w = append(w, fmt.Sprintf("'%s.%s'", s[0], s[1]))
	}
	return fmt.Sprintf(whereClause, strings.Join(w, ","))
}
