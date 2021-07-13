package main

import (
	"fmt"
	"strings"
	"testing"
)

const (
	constraintRestorePrimaryKeyAndUniqueKeyTable = "pk_table"
	constraintRestorePrimaryKeyAndUniqueKeyIndex = "pk_index"
	constraintRestorePrimaryKeyTest              = "pk_table_test"
	constraintRestorePrimaryKeyColumn            = "id"
	constraintRestoreUniqueKeyColumn             = "email"
	constraintRestoreForeignKeyTable             = "fk_table"
)

func createFakeTablesFromConstraintRestoreTest() {
	setDatabaseConfigForTest()
	cmdOptions.Tab.SchemaName = "mock_data_table4"
	cmdOptions.DontPrompt = true
	cmdOptions.Rows = 30
	postgresOrGreenplum()
	savedConstraints = map[string][]constraint{} // reset
	sql := `
		DROP SCHEMA IF EXISTS %[1]s CASCADE;
		CREATE SCHEMA %[1]s;
		DROP TABLE IF EXISTS %[1]s.%[2]s;
		CREATE TABLE %[1]s.%[2]s (%[4]s int primary key, %[5]s varchar, gender varchar CHECK (gender IN ('M', 'F', 'O')));
		DROP INDEX IF EXISTS %[1]s.%[3]s;
		CREATE UNIQUE INDEX %[3]s ON %[1]s.%[2]s (%[5]s);
		DROP TABLE IF EXISTS %[1]s.%[6]s;
		CREATE TABLE %[1]s.%[6]s (%[4]s int primary key, country varchar(10));
		ALTER TABLE %[1]s.%[6]s ADD CONSTRAINT %[6]s_fk FOREIGN KEY (%[4]s) REFERENCES %[1]s.%[2]s(%[4]s);
		DROP TABLE IF EXISTS %[1]s.%[7]s;
		CREATE TABLE %[1]s.%[7]s (%[4]s int primary key, %[5]s varchar, gender varchar CHECK (gender IN ('M', 'F', 'O')));
	`
	sql = fmt.Sprintf(sql, cmdOptions.Tab.SchemaName, constraintRestorePrimaryKeyAndUniqueKeyTable,
		constraintRestorePrimaryKeyAndUniqueKeyIndex, constraintRestorePrimaryKeyColumn,
		constraintRestoreUniqueKeyColumn, constraintRestoreForeignKeyTable, constraintRestorePrimaryKeyTest)
	_, err := ExecuteDB(sql)
	if err != nil {
		Fatalf("createFakeTablesFromConstraintRestoreTest, error in executing the statement, err: %v", err)
	}
	tables := allTablesPostgres(fmt.Sprintf("AND n.nspname='%s'", cmdOptions.Tab.SchemaName))
	MockTable(tables)
}

// Test: FixConstraints
func TestFixConstraints(t *testing.T) {
	// Skipping ..., since it will calls in other functions and they are
	// tested individually
}

// Test: fixPKey, check if the primary keys are fixed
func TestFixPKey(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	for _, conStr := range []string{"PRIMARY", "UNIQUE"} {
		data, ok := savedConstraints[conStr]
		t.Run("check_for_saved_constraint_"+conStr, func(t *testing.T) {
			if !ok {
				t.Errorf("TestFixPKey, expected = %v, found none on the saved map", conStr)
			}
		})
		t.Run("check_for_any_available_constraint_"+conStr, func(t *testing.T) {
			if got := len(data); got <= 0 {
				t.Errorf("TestFixPKey = %v, want > 0 %v constriants", got, conStr)
			}
		})
		for _, ck := range data {
			fixPKey(ck)
			t.Run("check_data_is_fixed_and_constraints_is_available_for_constraints_"+conStr, func(t *testing.T) {
				tab := GenerateTableName(constraintRestorePrimaryKeyAndUniqueKeyTable, cmdOptions.Tab.SchemaName)
				if got := GetConstraintsPertab(tab); len(got) != 3 {
					t.Errorf("TestFixPKey = %v, want %v", len(got), 3)
				}
			})
		}
	}
}

// Test: fixPKViolator, test if the function actually cleans up the pk violation
func TestFixPKViolator(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	for _, conStr := range []string{"PRIMARY", "UNIQUE"} {
		var col string = constraintRestorePrimaryKeyColumn
		var dtype string = "int"
		if conStr == "UNIQUE" {
			col = constraintRestoreUniqueKeyColumn
			dtype = "varchar"
		}
		t.Run("should_return_all_rows_fixed_for_constraint_"+conStr, func(t *testing.T) {
			tab := GenerateTableName(constraintRestorePrimaryKeyAndUniqueKeyTable, cmdOptions.Tab.SchemaName)
			fixPKViolator(tab, col, dtype)
			if got := getTotalPKViolator(tab, col); got > 0 {
				t.Errorf("TestFixPKViolator = %v, want = 0", got)
			}
		})
	}
}

// Test: fixFKey, check if the fk's are fixed and the constraint restored
func TestFixFKey(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	data, ok := savedConstraints["FOREIGN"]
	t.Run("check_for_saved_constraint_FOREIGN", func(t *testing.T) {
		if !ok {
			t.Errorf("TestFixFKey found none on the saved map")
		}
	})
	t.Run("check_for_any_available_constraint_FOREIGN", func(t *testing.T) {
		if got := len(data); got <= 0 {
			t.Errorf("TestFixFKey = %v, want > 0 %v constriants", got, "FOREIGN")
		}
	})
	for _, ck := range data {
		fixFKey(ck)
		t.Run("check_data_is_fixed_and_constraints_is_available_for_constraints_FOREIGN", func(t *testing.T) {
			tab := GenerateTableName(constraintRestoreForeignKeyTable, cmdOptions.Tab.SchemaName)
			if got := GetConstraintsPertab(tab); len(got) != 3 {
				t.Errorf("TestFixFKey = %v, want %v", len(got), 3)
			}
		})
	}
}

// Test: getForeignKeyColumns, should provide us with value fkcolumn, fktable, ref table, ref column
func TestGetForeignKeyColumns(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	data, _ := savedConstraints["FOREIGN"]
	var found = false
	for _, ck := range data {
		keys := getForeignKeyColumns(ck)
		tab := GenerateTableName(constraintRestoreForeignKeyTable, cmdOptions.Tab.SchemaName)
		refT := fmt.Sprintf("%s.%s", cmdOptions.Tab.SchemaName, constraintRestorePrimaryKeyAndUniqueKeyTable)
		if keys.Table == tab {
			found = true
			t.Run("check_if_all_key_matches", func(t *testing.T) {
				if keys.Table != tab {
					t.Errorf("TestGetForeignKeyColumns(Table Check) = %v, want %v", keys.Table,
						constraintRestoreForeignKeyTable)
				}
				if keys.Column != constraintRestorePrimaryKeyColumn {
					t.Errorf("TestGetForeignKeyColumns(Column Check) = %v, want %v", keys.Column,
						constraintRestorePrimaryKeyColumn)
				}
				if keys.Reftable != refT {
					t.Errorf("TestGetForeignKeyColumns(Ref Tab Check) = %v, want %v", keys.Reftable, refT)
				}
				if keys.Refcolumn != constraintRestorePrimaryKeyColumn {
					t.Errorf("TestGetForeignKeyColumns(Ref Column Check) = %v, want %v", keys.Refcolumn,
						constraintRestorePrimaryKeyColumn)
				}
			})
		}
	}
	if !found {
		t.Errorf("TestGetForeignKeyColumns, unable to find the fk table")
	}
}

// Test: IgnoreErrorString, check if this ignore these error messages
func TestIgnoreErrorString(t *testing.T) {
	tests := []struct {
		name   string
		errMsg string
		want   bool
	}{
		{"error_multiple_primary_keys_for_table", "ERROR #42P16 multiple primary keys for table \"x\" are not allowed",
			true},
		{"error_table_already_exists", "table x already exists", true},
		{"error_table_does_not_exists", "table x does not exists", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IgnoreErrorString(tt.errMsg); got != tt.want {
				t.Errorf("TestIgnoreErrorString = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: recreateAllConstraints, test if all recreating is done
func TestRecreateAllConstraints(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	failedConstraintsFile := fmt.Sprintf("%s/failed_constraint_creations.sql", Path)
	recreateAllConstraints()
	for _, conStr := range []string{"PRIMARY", "UNIQUE", "FOREIGN"} {
		t.Run("recreate_constraint_"+conStr, func(t *testing.T) {
			tab := GenerateTableName(constraintRestorePrimaryKeyAndUniqueKeyTable, cmdOptions.Tab.SchemaName)
			if conStr == "FOREIGN" {
				tab = GenerateTableName(constraintRestoreForeignKeyTable, cmdOptions.Tab.SchemaName)
			}
			if got := GetConstraintsPertab(tab); len(got) != 3 {
				t.Errorf("TestRecreateAllConstraints = %v, want %v", len(got), 3)
			}
		})
	}
	t.Run("should_have_at_least_one_constraint_reported_on_the_file", func(t *testing.T) {
		c, _ := ReadFile(failedConstraintsFile)
		if got := len(c); got <= 0 {
			t.Errorf("TestRecreateAllConstraints = %v, want > 0", c)
		}
	})
}

// Test: deleteViolatingPkOrUkConstraints, should delete violating rows
func TestDeleteViolatingPkOrUkConstraints(t *testing.T) {
	createFakeTablesFromConstraintRestoreTest()
	executeStatement := func(stmt string) error {
		_, err := ExecuteDB(stmt)
		return err
	}
	t.Run("should_delete_all_violating_rows", func(t *testing.T) {
		col := "id"
		tab := GenerateTableName(constraintRestorePrimaryKeyTest, cmdOptions.Tab.SchemaName)
		dropPkStatement := fmt.Sprintf(`ALTER TABLE %[1]s DROP CONSTRAINT %[2]s_pkey;`,
			tab, constraintRestorePrimaryKeyTest)
		err := executeStatement(dropPkStatement)
		if err != nil {
			t.Errorf("TestDeleteViolatingPkOrUkConstraints, drop failed with err: %v", err)
		}
		loadDuplicateDate := fmt.Sprintf(`INSERT INTO %[1]s SELECT %[2]s FROM %[1]s`, tab, col)
		err = executeStatement(loadDuplicateDate)
		if err != nil {
			t.Errorf("TestDeleteViolatingPkOrUkConstraints, insert failed with err: %v", err)
		}
		err = deleteViolatingConstraintKeys(tab, "id")
		if err != nil {
			t.Errorf("TestDeleteViolatingPkOrUkConstraints, delete failed with err: %v", err)
		}
		addPkStatement := fmt.Sprintf(`ALTER TABLE %[1]s ADD CONSTRAINT %[2]s_pkey PRIMARY KEY(%[2]s);`,
			tab, col, constraintRestorePrimaryKeyTest)
		err = executeStatement(addPkStatement)
		if err != nil {
			t.Errorf("TestDeleteViolatingPkOrUkConstraints, duplicate entry clean up failed with err: %v", err)
		}
	})
}

// Test: ExtractTableNColumnName, extract table and column from given PK or UK constraint DDL
func TestExtractTableNColumnName(t *testing.T) {
	removeWhiteSpaces := func(s string) string {
		return strings.Trim(s, " ")
	}
	tests := []struct {
		name       string
		s          string
		wantTable  string
		wantColumn string
	}{
		{"column_extraction_alter_statement",
			"ALTER TABLE testme ADD CONSTRAINT something UNIQUE (dist_id, zipcode);",
			"testme", "dist_id, zipcode"},
		{"column_extraction_on_using_clause",
			"CREATE UNIQUE INDEX film_fulltext_idx ON public.film USING gist (fulltext);",
			"public.film", "fulltext"},
		{"column_extraction_brackets_with_brackets",
			"CREATE UNIQUE INDEX ir_translation_code_unique ON public.ir_translation USING " +
				"btree (type, lang, md5(src)) WHERE ((type)::text = 'code'::text);",
			"public.ir_translation", "type, lang, md5(src)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTable, gotColumn := ExtractTableNColumnName(tt.s);
				removeWhiteSpaces(gotTable) != removeWhiteSpaces(tt.wantTable) ||
					removeWhiteSpaces(gotColumn) != removeWhiteSpaces(tt.wantColumn) {
				t.Errorf("TestExtractTableNColumnName() = %v / %v, want %v / %v",
					gotTable, gotColumn, tt.wantTable, tt.wantColumn)
			}
		})
	}
}
