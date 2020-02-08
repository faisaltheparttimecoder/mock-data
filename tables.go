package main

import (
	"fmt"
	"strings"
)

func CreateFakeTables() {
	Infof("Create %d fake tables on the database: %s", cmdOptions.Tab.TotalTables, cmdOptions.Database)

	// Start the progress bar
	bar := StartProgressBar("Creating Tables", cmdOptions.Tab.TotalTables)

	// Create tables
	for i := 1; i <= cmdOptions.Tab.TotalTables; i++ {
		createTable(i)
		bar.Add(1)
	}
	fmt.Println()

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
	howManyColumns, _ := RandomInt(2, cmdOptions.Tab.MaxColumns)

	// DataType
	dataTypes := SupportedDataTypes()
	totalDataTypes := len(dataTypes)
	tableColStatement := ""

	// Column + DataTypes
	for  i := 1; i <= howManyColumns; i++ {
		j, _ := RandomInt(1, totalDataTypes)
		if j <= totalDataTypes { // if the random number is greater than the array it would fail
			tableColStatement = tableColStatement + fmt.Sprintf(
				"%s_%d %s", columnName, i, dataTypes[j])
		}
	}

	// do we have a column
	if IsStringEmpty(tableColStatement) {
		Warnf("The table %s has no columns", tableName)
	}

	// Create table statement
	return createTableDDL + fmt.Sprintf("%s);", strings.TrimRight(tableColStatement, ",") )
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