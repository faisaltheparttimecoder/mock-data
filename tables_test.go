package main

import (
	"fmt"
	"testing"
)

// Test: CreateFakeTables, test the creation of fake tables
func TestCreateFakeTables(t *testing.T) {
	setDatabaseConfigForTest()
	cmdOptions.Tab.SchemaName = "mock_data_table"
	cmdOptions.Tab.TotalTables = 100
	createSchemaStmt := fmt.Sprintf("CREATE SCHEMA %s;", cmdOptions.Tab.SchemaName)
	_, err := ExecuteDB(createSchemaStmt)
	if err != nil {
		t.Errorf("TestCreateFakeTables, failed creating schema, err: %v", err)
	}
	CreateFakeTables()
	t.Run("should_provide_the_total_table_created", func(t *testing.T) {
		if got := allTablesPostgres(fmt.Sprintf("AND n.nspname='%s'", cmdOptions.Tab.SchemaName));
			len(got) != cmdOptions.Tab.TotalTables {
			t.Errorf("TestCreateFakeTables = %v, want %v", len(got), cmdOptions.Tab.TotalTables)
		}
	})
}

// Test: tableNameGenerator, test name should match the regex
func TestTableNameGenerator(t *testing.T) {
	setDatabaseConfigForTest()
	randomNumber := 456
	cmdOptions.Tab.TableNamePrefix = "mock_data_table_random"
	t.Run("should_provide_a_random_table_name", func(t *testing.T) {
		re := fmt.Sprintf("%s_[a-zA-Z0-9]*_%d", cmdOptions.Tab.TableNamePrefix, randomNumber)
		if got := tableNameGenerator(randomNumber); !doesDataMatchDataType(got, re) {
			t.Errorf("TestTableNameGenerator = %v, want a valid match", got)
		}
	})
}

// Test: createTableStatementGenerator, should generate different table name based on parameter passed
func TestCreateTableStatementGenerator(t *testing.T) {
	setDatabaseConfigForTest()
	cmdOptions.Tab.SchemaName = "mock_data_table2"
	cmdOptions.Tab.ColumnNamePrefix = "mock_data_column"
	randomNumber := 9877
	createSchemaStmt := fmt.Sprintf("CREATE SCHEMA %s;", cmdOptions.Tab.SchemaName)
	_, err := ExecuteDB(createSchemaStmt)
	if err != nil {
		t.Errorf("TestCreateTableStatementGenerator, failed creating schema, err: %v", err)
	}

	tests := []struct {
		name          string
		caseSensitive bool
		re            string
	}{
		{"generate_case_sensitive_table_name", true, "%s_[a-zA-Z0-9]*_%d"},
		{"generate_case_insensitive_table_name", false, "%s_[a-z0-9]*_%d"},
	}

	for _, tt := range tests {
		cmdOptions.Tab.CaseSensitive = tt.caseSensitive
		t.Run(tt.name, func(t *testing.T) {
			stmt := createTableStatementGenerator(randomNumber)
			if _, err := ExecuteDB(stmt); err != nil {
				t.Errorf("TestCreateTableStatementGenerator, stmt execution failed, err: %v", err)
			}
		})

		tableFromDb := allTablesPostgres(fmt.Sprintf("AND n.nspname='%s'", cmdOptions.Tab.SchemaName))
		if len(tableFromDb) > 0 {
			table := tableFromDb[0].Table
			schema := tableFromDb[0].Schema
			tab := GenerateTableName(table, schema)

			t.Run(tt.name+"_table_exists_check", func(t *testing.T) { // Checking for the name casesensitivity
				re := fmt.Sprintf(tt.re, cmdOptions.Tab.TableNamePrefix, randomNumber)
				if !doesDataMatchDataType(table, re) {
					t.Errorf("TestCreateTableStatementGenerator: can't find the table: %v in the database", tab)
				}
			})
			t.Run(tt.name+"_column_check", func(t *testing.T) { // column should be atleast > 2 and < 1600
				got := columnExtractorPostgres(schema, table)
				if len(got) < 2 && len(got) > 1599 {
					t.Errorf("TestCreateTableStatementGenerator: have no columns or the columns are way to high")
				}
			})
			t.Run(tt.name+"_datatype_check", func(t *testing.T) { // datatype check
				got := getDatatype(tab, []string{fmt.Sprintf("%s_%d", cmdOptions.Tab.ColumnNamePrefix, 1)})
				dt := got[0].Dtype
				if !StringContains(dt+",", SupportedDataTypes()) {
					t.Errorf("TestCreateTableStatementGenerator = %v, isn't part of the supported list", dt)
				}
			})
			dropTableStmt := fmt.Sprintf("DROP TABLE %s;", tab)
			_, err := ExecuteDB(dropTableStmt)
			if err != nil {
				t.Errorf("TestCreateTableStatementGenerator, failed in dropping table, err: %v", err)
			}

		} else {
			t.Errorf("TestCreateTableStatementGenerator, failed to extract tables from the database")
		}
	}
}

// Test: createTable
func TestCreateTable(t *testing.T) {
	// Already validated this with TestCreateFakeTables & TestCreateTableStatementGenerator,
	// there is nothing different here, skip...
}

// Test: MockTables
func TestMockTables(t *testing.T) {
	// It basically calls all the other functions and those functions are already validated
	// there is nothing different here, skip...
}

// Test: generateWhereClause, should generated appropriate where clause
func TestGenerateWhereClause(t *testing.T) {
	cmdOptions.Tab.SchemaName = "public" // the default
	tests := []struct {
		name     string
		tableRow string
		want     string
	}{
		{"single_table", "tab1", "ANDnnspnamecrelnameINpublictab1"},
		{"multiple_table", "tab1,tab2", "ANDnnspnamecrelnameINpublictab1publictab2"},
		{"single_table_schema", "schema1.tab1", "ANDnnspnamecrelnameINschema1tab1"},
		{"multiple_table_schema", "schema1.tab1,schema1.tab2",
			"ANDnnspnamecrelnameINschema1tab1schema1tab2"},
		{"multiple_schema_table_schema", "schema1.tab1,schema1.tab2,schema2.tab1,schema2.tab2",
			"ANDnnspnamecrelnameINschema1tab1schema1tab2schema2tab1schema2tab2"},
		{"multiple_schema_table_schema_and_no_schema_tab",
			"schema1.tab1,schema1.tab2,schema2.tab3,schema2.tab4,tab5",
			"ANDnnspnamecrelnameINschema1tab1schema1tab2schema2tab3schema2tab4publictab5"},
	}
	for _, tt := range tests {
		cmdOptions.Tab.FakeTablesRows = tt.tableRow
		if got := RemoveSpecialCharacters(generateWhereClause()); got != tt.want {
			t.Errorf("TestGenerateWhereClause = %v, want %v", got, tt.want)
		}
	}
}
