package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"testing"
)

const (
	workerTestConstraintTable           = "constraint_table"
	workerTestSingleColumnTable         = "single_column"
	workerTestContainSerialColumn       = "contain_serial_column"
	workerTestUnsupportedDataTypeColumn = "unsupported_datatype"
	workerTestCopyData                  = "copy_table"
)

func setDatabaseConfigForTest() {
	cmdOptions.Database = viper.GetString("PGDATABASE")
	cmdOptions.Password = viper.GetString("PGPASSWORD")
	cmdOptions.Username = viper.GetString("PGUSER")
	cmdOptions.Hostname = viper.GetString("PGHOST")
	cmdOptions.Port = viper.GetInt("PGPORT")
}

// Create mocking tables for worker_test
func createFakeTablesForWorkerTest() []DBTables {
	setDatabaseConfigForTest()
	postgresOrGreenplum()

	// Fake tables
	cmdOptions.Tab.SchemaName = "mock_data_table3"
	cmdOptions.Rows = 67
	cmdOptions.DontPrompt = true
	sql := `
		DROP SCHEMA IF EXISTS %[1]s CASCADE; 
		CREATE SCHEMA %[1]s;
		DROP TABLE IF EXISTS %[1]s.%[2]s; 
		CREATE TABLE %[1]s.%[2]s (id int primary key, name varchar);
		DROP TABLE IF EXISTS %[1]s.%[3]s; 
		CREATE TABLE %[1]s.%[3]s (id serial);
		DROP TABLE IF EXISTS %[1]s.%[4]s; 
		CREATE TABLE %[1]s.%[4]s (id serial, name varchar);
		DROP TYPE IF EXISTS foo;
		CREATE TYPE foo AS (f1 int, f2 text); 
		DROP TABLE IF EXISTS %[1]s.%[5]s; 
		CREATE TABLE %[1]s.%[5]s (id serial, name foo);
		DROP TABLE IF EXISTS %[1]s.%[6]s; 
		CREATE TABLE %[1]s.%[6]s (name varchar(10), salary money);
	`
	sql = fmt.Sprintf(sql, cmdOptions.Tab.SchemaName, workerTestConstraintTable,
		workerTestSingleColumnTable, workerTestContainSerialColumn,
		workerTestUnsupportedDataTypeColumn, workerTestCopyData)
	_, err := ExecuteDB(sql)
	if err != nil {
		Fatalf("createFakeTablesForTest, failed to create sql, err %v", err)
	}
	CreateFakeTables()
	return allTablesPostgres(fmt.Sprintf("AND n.nspname = '%s'", cmdOptions.Tab.SchemaName))
}

// Test: MockTable, fake the tables given
func TestMockTable(t *testing.T) {
	tableList := createFakeTablesForWorkerTest()
	tests := []struct {
		name  string
		cc    bool
		total int
	}{
		{"check_row_count_with_constraint_for_tab_%s", false, 2},
		{"check_row_count_without_constraint_for_tab_%s", true, 0},
	}
	for _, test := range tests {
		cmdOptions.IgnoreConstraint = test.cc
		MockTable(tableList)

		// Should produce some rows
		for _, tt := range tableList {
			if tt.Table == workerTestUnsupportedDataTypeColumn { // skip this table
				continue
			}
			tab := GenerateTableName(tt.Table, tt.Schema)
			t.Run(fmt.Sprintf("check_rows_count_for_table_%s", tab), func(t *testing.T) {
				if got := TotalRows(tab); got <= 0 {
					t.Errorf("TestMockTable = %v, want > 0", got)
				}
			})
		}

		// check if the unsupported datatype is registered in the slick
		t.Run("check_if_supported_datatype_is_found", func(t *testing.T) {
			if got := len(skippedTab); got <= 0 {
				t.Errorf("TestMockTable = %v, want > 0", got)
			}
		})

		// constraint should be available and not available based on the flag
		t.Run(fmt.Sprintf("the_constriaint_should_be_%v", !test.cc), func(t *testing.T) {
			tab := GenerateTableName(workerTestConstraintTable, cmdOptions.Tab.SchemaName)
			if got := GetConstraintsPertab(tab); len(got) != test.total {
				t.Errorf("TestMockTable = %v, want %v", got, test.total)
			}
		})
	}
}

// Test: tableMocker, Execute the table mocker
func TestTableMocker(t *testing.T) {
	// Already verified using the test TestMockTable, skipped ...
}

// Test: checkIfOneColumnIsASerialDatatype check for the single column function
func TestCheckIfOneColumnIsASerialDatatype(t *testing.T) {
	_ = createFakeTablesForWorkerTest()
	checkIfOneColumnIsASerialDatatype(DBTables{Table: workerTestContainSerialColumn,
		Schema: cmdOptions.Tab.SchemaName}, []DBColumns{
		DBColumns{"id", "integer", "nextval('id_seq'::regclass)"},
	})
	t.Run("check_if_the_table_is_registered_correctly", func(t *testing.T) {
		if got := len(oneColumnTable); got <= 0 {
			t.Errorf("TestCheckIfOneColumnIsASerialDatatype = %v, want > 0", got)
		}
	})
}

// Test: columnExtractor, check if its able to extract column
func TestColumnExtractor(t *testing.T) {
	tables := createFakeTablesForWorkerTest()
	allColumns := columnExtractor(tables)
	t.Run("did_it_find_any_column", func(t *testing.T) {
		if got := len(allColumns); got <= 0 {
			t.Errorf("TestColumnExtractor = %d, want > 0", got)
		}
	})
	t.Run("did_it_find_single_column_table_with_serial_datatype", func(t *testing.T) {
		if got := len(oneColumnTable); got <= 0 {
			t.Errorf("TestColumnExtractor = %d, want > 0", got)
		}
	})
	t.Run("did_it_place_any_serial_columns_in_the_list", func(t *testing.T) {
		for _, cols := range allColumns {
			for _, c := range cols.Columns {
				if strings.HasPrefix(c.Sequence, "nextval") {
					t.Errorf("TestColumnExtractor = %v, it shouldn't be part of the collection", c.Sequence)
					break
				}
			}
		}
	})
}

// Test: BackupConstraintsAndStartDataLoading
func TestBackupConstraintsAndStartDataLoading(t *testing.T) {
	// All the function involved inside has been tested individually, skipped ....
}

// Test: CommitData
func TestCommitData(t *testing.T) {
	// This has been verified by others function, nothing to add specifically, skipped...
}

// Test: CopyData
func TestCopyData(t *testing.T) {
	createFakeTablesForWorkerTest()
	var result struct {
		Count int
	}
	tab := GenerateTableName(workerTestCopyData, cmdOptions.Tab.SchemaName)
	col := []string{"name", "salary"}
	data := []string{"john", "10000"}
	db := ConnectDB()
	CopyData(tab, col, data, db)
	sql := fmt.Sprintf("SELECT COUNT(*) AS count FROM %s WHERE name = '%s';", tab, "john")
	_, err := db.Query(&result, sql)
	if err != nil {
		t.Errorf("TestCopyData error when getting data from the database, err: %v", err)
	}
	t.Run("should_return_one_successful_row", func(t *testing.T) {
		if result.Count != 1 {
			t.Errorf("TestCopyData = %v, want %v", result.Count, 1)
		}
	})
}

// Test: addDataIfItsASerialDatatype, if it has only one column and that is serial
// it should generate rows for it
func TestAddDataIfItsASerialDatatype(t *testing.T) {
	_ = createFakeTablesForWorkerTest()
	addDataIfItsASerialDatatype()
	t.Run("single_column_table_should_return_rows", func(t *testing.T) {
		if got := TotalRows(GenerateTableName(workerTestSingleColumnTable, cmdOptions.Tab.SchemaName)); got <= 0 {
			t.Errorf("TestAddDataIfItsASerialDatatype = %v, want > 0", got)
		}
	})
}

// Test: isItSerialDatatype, is the datatype of the column serial
func TestIsItSerialDatatype(t *testing.T) {
	tests := []struct {
		name   string
		column DBColumns
		want   bool
	}{
		{"serial_column",
			DBColumns{"actor_id", "integer", "nextval('actor_actor_id_seq'::regclass)"}, true},
		{"non_serial_column",
			DBColumns{"first_name", "character varying(45) ", " now()"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isItSerialDatatype(tt.column); got != tt.want {
				t.Errorf("TestIsItSerialDatatype = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: GenerateTableName, check if it create valid table name with quotes
func TestGenerateTableName(t *testing.T) {
	t.Run("should_generate_valid_table_name", func(t *testing.T) {
		tab := "mock_data_test_table"
		schema := "mock_data_test_schema"
		want := fmt.Sprintf("\"%s\".\"%s\"", schema, tab)
		if got := GenerateTableName(tab, schema); got != want {
			t.Errorf("TestGenerateTableName = %v, want %v", got, want)
		}
	})
}
