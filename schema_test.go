package main

import (
	"fmt"
	"testing"
)

// Test: MockSchema, check if the data is loaded with schema keyword
func TestMockSchema(t *testing.T) {
	setDatabaseConfigForTest()
	postgresOrGreenplum()
	fakeSchema := "mock_data_schema"
	fakeTable := "fake_table"
	createSchemaStmt := fmt.Sprintf("CREATE SCHEMA %s;", fakeSchema)
	createFakeSchemaTableStmt := `
	CREATE TABLE %[1]s.%[2]s1 (id serial, name text, active bool);
	CREATE TABLE %[1]s.%[2]s2 (country varchar, active_date date, last_login timestamp);
	CREATE TABLE %[1]s.%[2]s3 (rating int, price money, balance numeric(4,2));
	CREATE TABLE %[1]s.%[2]s4 (inactive bool, gender char, address text);
	CREATE TABLE %[1]s.%[2]s5 (category varchar, comment text, feedback varchar(500));
`
	_, err := ExecuteDB(createSchemaStmt)
	if err != nil {
		t.Errorf("TestMockSchema, failed creating schema, err: %v", err)
	}
	_, err = ExecuteDB(fmt.Sprintf(createFakeSchemaTableStmt, fakeSchema, fakeTable))
	if err != nil {
		t.Errorf("TestMockSchema, failed creating schema tables, err: %v", err)
	}
	cmdOptions.SchemaName = fakeSchema
	cmdOptions.Rows = 100
	MockSchema()
	for i := 0; i < 5; i++ {
		tableNumber := i + 1
		tabName := fmt.Sprintf("%s%d", fakeTable, tableNumber)
		tableNameWithSchema := GenerateTableName(tabName, fakeSchema)
		t.Run(fmt.Sprintf("checking_rows_count_of_table_%s", tabName), func(t *testing.T) {
			if got := TotalRows(tableNameWithSchema); got != cmdOptions.Rows {
				t.Errorf("TestMockSchema = %v, want %v", got, cmdOptions.Rows)
			}
		})
	}
}

