package main

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"testing"
)

// TODO: remove this at the end
func setDatabaseConfigForTest() {
	cmdOptions.Database = viper.GetString("PGDATABASE")
	cmdOptions.Password = viper.GetString("PGPASSWORD")
	cmdOptions.Username = viper.GetString("PGUSER")
	cmdOptions.Hostname = viper.GetString("PGHOST")
	cmdOptions.Port = viper.GetInt("PGPORT")
}

var (
	fakePkViolationTableName    = "pk_violation_table"
	fakePkColumnName            = "id"
	fakePkNonViolationTableName = "pk_non_violation_table"
	fakeSqlPKFormat             = `
  DROP TABLE IF EXISTS %[1]s;
  CREATE TABLE %[1]s ( %[2]s int, name varchar );
  INSERT INTO %[1]s VALUES (1, 'john');
  INSERT INTO %[1]s VALUES (1, 'jack');
  INSERT INTO %[1]s VALUES (2, 'jill');
  INSERT INTO %[1]s VALUES (2, 'james');
  INSERT INTO %[1]s VALUES (3, 'jerry');

  DROP TABLE IF EXISTS %[3]s;
  CREATE TABLE %[3]s ( %[2]s int, name varchar );
  INSERT INTO %[3]s VALUES (1, 'john');
  INSERT INTO %[3]s VALUES (2, 'jill');
  INSERT INTO %[3]s VALUES (3, 'jerry');
`
	fakePKSql                   = fmt.Sprintf(fakeSqlPKFormat, fakePkViolationTableName,
		fakePkColumnName, fakePkNonViolationTableName)
	fakeFkViolationTableName    = "fk_violation_table"
	fakeFkColumnName            = "id"
	fakeFkNonViolationTableName = "fk_non_violation_table"
	fakeSqlFKFomat              = `
  DROP TABLE IF EXISTS %[1]s;
  CREATE TABLE %[1]s ( %[2]s int, country varchar );
  INSERT INTO %[1]s VALUES (50, 'us');
  INSERT INTO %[1]s VALUES (60, 'india');
  INSERT INTO %[1]s VALUES (70, 'spain');
  INSERT INTO %[1]s VALUES (80, 'new zealand');
  INSERT INTO %[1]s VALUES (90, 'south africa');

  DROP TABLE IF EXISTS %[3]s;
  CREATE TABLE %[3]s ( %[2]s int, country varchar );
  INSERT INTO %[3]s VALUES (1, 'us');
  INSERT INTO %[3]s VALUES (2, 'india');
  INSERT INTO %[3]s VALUES (3, 'spain');
`
	fakeFKSql                     = fmt.Sprintf(fakeSqlFKFomat, fakeFkViolationTableName,
		fakeFkColumnName, fakeFkNonViolationTableName)
	fakeSql                       = fakePKSql + fakeFKSql
)

// Test: dbVersion, check if the command provides the database version
func TestDbVersion(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("check_database_version", func(t *testing.T) {
		// If its successful then it prints message on the screen, which means its a success else it craps out
		dbVersion()
	})
}

// Test: postgresOrGreenplum, check if this is a postgres or a greenplum database
func TestPostgresOrGreenplum(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("check_which_flavour_of_postgres_is_this", func(t *testing.T) {
		postgresOrGreenplum()
		if GreenplumOrPostgres != "postgres" { // since all tests are designed for a postgres containers
			t.Errorf("TestPostgresOrGreenplum = %v, want %v", GreenplumOrPostgres, "postgres")
		}
	})
}

// Test: allTablesPostgres, should produce all the tables
func TestAllTablesPostgres(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("should_extract_all_tables", func(t *testing.T) {
		postgresOrGreenplum()
		ExecuteDemoDatabase()
		if got := allTablesPostgres(""); len(got) != 15 {
			t.Errorf("TestAllTablesPostgres = %v, want %d tables", len(got), 15)
		}
	})
}

// Test: allTablesGPDB, should produce all the tables for GPDB
func TestAllTablesGPDB(t *testing.T) {
	// Not implemented since we dont have a GPDB container to run test
}

// Test: columnExtractorPostgres, should produce all the columns and datatype for the given table
func TestColumnExtractorPostgres(t *testing.T) {
	tests := []struct {
		name   string
		schema string
		table  string
		want   int
	}{
		{"column_actor_table", "public", "actor", 4},
		{"column_address_table", "public", "address", 8},
		{"column_category_table", "public", "category", 3},
		{"column_city_table", "public", "city", 4},
		{"column_country_table", "public", "country", 3},
		{"column_customer_table", "public", "customer", 10},
		{"column_film_table", "public", "film", 14},
		{"column_film_actor_table", "public", "film_actor", 3},
		{"column_film_category_table", "public", "film_category", 3},
		{"column_inventory_table", "public", "inventory", 4},
		{"column_language_table", "public", "language", 3},
		{"column_payment_table", "public", "payment", 6},
		{"column_rental_table", "public", "rental", 7},
		{"column_staff_table", "public", "staff", 11},
		{"column_store_table", "public", "store", 4},
	}
	setDatabaseConfigForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := columnExtractorPostgres(tt.schema, tt.table); len(got) != tt.want {
				t.Errorf("TestColumnExtractorPostgres = %v, wanted %v", len(got), tt.want)
			}
		})
	}
}

// Test: columnExtractorGPDB, should produce all the column for a given table for GPDB
func TestColumnExtractorGPDB(t *testing.T) {
	// Not implemented since we dont have a GPDB container to run test
}

// Test: GetPGConstraintDDL, should produce DDL for constraints
func TestGetPGConstraintDDL(t *testing.T) {
	tests := []struct {
		name     string
		connType string
		want     int
	}{
		{"constraint_primary_key_ddl", "p", 15},
		{"constraint_foreign_key_ddl", "f", 18},
		{"constraint_check_key_ddl", "c", 1},
		{"constraint_unique_key_ddl", "u", 1},
	}
	setDatabaseConfigForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPGConstraintDDL(tt.connType); len(got) != tt.want {
				t.Errorf("TestGetPGConstraintDDL = %v, wanted %v", len(got), tt.want)
			}
		})
	}
}

// Test: GetPGIndexDDL, should produce DDL for all indexes
func TestGetPGIndexDDL(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("get_index_ddl", func(t *testing.T) {
		if got := GetPGIndexDDL(); len(got) != 18 {
			t.Errorf("TestGetPGIndexDDL = %v, wanted %v", len(got), 18)
		}
	})
}

// Test: GetConstraintsPertab, should produce all constraints based on table name
func TestGetConstraintsPertab(t *testing.T) {
	tests := []struct {
		name  string
		table string
		want  int
	}{
		{"constraint_actor_table", "public.actor", 2},
		{"constraint_address_table", "public.address", 2},
		{"constraint_category_table", "public.category", 2},
		{"constraint_city_table", "public.city", 2},
		{"constraint_country_table", "public.country", 1},
		{"constraint_customer_table", "public.customer", 2},
		{"constraint_film_table", "public.film", 2},
		{"constraint_film_actor_table", "public.film_actor", 3},
		{"constraint_film_category_table", "public.film_category", 3},
		{"constraint_inventory_table", "public.inventory", 2},
		{"constraint_language_table", "public.language", 1},
		{"constraint_payment_table", "public.payment", 4},
		{"constraint_rental_table", "public.rental", 4},
		{"constraint_staff_table", "public.staff", 2},
		{"constraint_store_table", "public.store", 3},
	}
	setDatabaseConfigForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConstraintsPertab(tt.table); len(got) != tt.want {
				t.Errorf("GetConstraintsPertab = %v, wanted %v", len(got), tt.want)
			}
		})
	}
}

// Test: getDatatype, should produce all the datatypes of the column for the table
func TestGetDatatype(t *testing.T) {
	tests := []struct {
		name   string
		table  string
		column []string
		want   string
	}{
		{"datatype_actor_table", "public.actor", []string{"first_name"}, "character varying(45)"},
		{"datatype_address_table", "public.address", []string{"address_id"}, "integer"},
		{"datatype_category_table", "public.category", []string{"last_update"}, "timestamp without time zone"},
		{"datatype_city_table", "public.city", []string{"last_update"}, "timestamp without time zone"},
		{"datatype_country_table", "public.country", []string{"last_update"}, "timestamp without time zone"},
		{"datatype_customer_table", "public.customer", []string{"activebool"}, "boolean"},
		{"datatype_film_table", "public.film", []string{"user_rating"}, "rating"},
		{"datatype_film_actor_table", "public.film_actor", []string{"last_update"}, "timestamp without time zone"},
		{"datatype_film_category_table", "public.film_category", []string{"last_update"}, "timestamp without time zone"},
		{"datatype_inventory_table", "public.inventory", []string{"film_id"}, "smallint"},
		{"datatype_language_table", "public.language", []string{"name"}, "character(20)"},
		{"datatype_payment_table", "public.payment", []string{"amount"}, "numeric(5,2)"},
		{"datatype_rental_table", "public.rental", []string{"return_date"}, "timestamp without time zone"},
		{"datatype_staff_table", "public.staff", []string{"picture"}, "bytea"},
		{"datatype_store_table", "public.store", []string{"address_id"}, "smallint"},
	}
	setDatabaseConfigForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDatatype(tt.table, tt.column)
			if len(got) > 0 {
				if got[0].Dtype != tt.want {
					t.Errorf("TestGetDatatype = %v, wanted %v", got[0].Dtype, tt.want)
				}
			} else {
				t.Errorf("TestGetDatatype = %v, want %v", len(got), len(tt.column))
			}
		})
	}
}

// Test: getTotalPKViolator, check the table reports any PK violation
func TestGetTotalPKViolator(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakePKSql)
	if err != nil {
		Fatalf("TestGetTotalPKViolator, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name  string
		table string
		want  int
	}{
		{"pk_violation_test_total", fakePkViolationTableName, 2},
		{"pk_non_violation_test_total", fakePkNonViolationTableName, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTotalPKViolator(tt.table, fakePkColumnName); got != tt.want {
				t.Errorf("TestGetTotalPKViolator = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: getPKViolator, check the sql generated for PK violation is valid
func TestGetPKViolator(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakePKSql)
	if err != nil {
		Fatalf("TestGetPKViolator, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name  string
		table string
		want  string
	}{
		{"pk_violation_test_sql", fakePkViolationTableName,
			"SELECT id FROM pk_violation_table GROUP BY id HAVING COUNT(*) > 1"},
		{"pk_non_violation_test_sql", fakePkNonViolationTableName,
			"SELECT id FROM pk_non_violation_table GROUP BY id HAVING COUNT(*) > 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPKViolator(tt.table, fakePkColumnName); got != tt.want {
				t.Errorf("TestGetPKViolator = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: GetPKViolators, check the data of PK violation
func TestGetPKViolators(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakePKSql)
	if err != nil {
		Fatalf("TestGetPKViolators, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name  string
		table string
		want  int
	}{
		{"pk_violation_test_rows", fakePkViolationTableName, 2},
		{"pk_non_violation_test_rows", fakePkNonViolationTableName, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPKViolators(tt.table, fakePkColumnName); len(got) != tt.want {
				t.Errorf("TestGetPKViolators = %v, want %v", len(got), tt.want)
			}
		})
	}
}

// Test: UpdatePKKey, make sure all the PK is fixed
func TestUpdatePKKey(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakePKSql)
	if err != nil {
		Fatalf("TestUpdatePKKey, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name     string
		table    string
		oldValue string
		newValue string
	}{
		{"pk_violation_test_rows_with_value_2", fakePkViolationTableName, "2", "10"},
		{"pk_violation_test_rows_with_value_1", fakePkViolationTableName, "1", "30"},
	}
	for _, tt := range tests {
		UpdatePKKey(tt.table, fakePkColumnName, tt.oldValue, tt.newValue)
	}
	t.Run("pk_violation_clearance", func(t *testing.T) {
		if got := getTotalPKViolator(fakePkViolationTableName, fakePkColumnName); got != 0 {
			t.Errorf("TestUpdatePKKey = %v, want %v", got, 0)
		}
	})
}

// Test: getFKViolator, make sure we get a valid fk check sql
func TestGetFKViolator(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakeSql)
	if err != nil {
		Fatalf("TestGetFKViolator, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name string
		fk   ForeignKey
		want string
	}{
		{"fk_violation_test_sql", ForeignKey{fakeFkViolationTableName, fakeFkColumnName,
			fakePkViolationTableName, fakePkColumnName},
			"SELECTidFROMfk_violation_tableWHEREidNOTIN(SELECTidFROMpk_violation_table)"},
		{"fk_non_violation_test_sql", ForeignKey{fakeFkNonViolationTableName, fakeFkColumnName,
			fakePkNonViolationTableName, fakePkColumnName},
			"SELECTidFROMfk_non_violation_tableWHEREidNOTIN(SELECTidFROMpk_non_violation_table)"},
	}
	format := regexp.MustCompile(`[\n\s]*`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format.ReplaceAllString(getFKViolator(tt.fk), "")
			if got != tt.want {
				t.Errorf("TestGetFKViolator = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: GetTotalFKViolators, get the total FK violators
func TestGetTotalFKViolators(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakeSql)
	if err != nil {
		Fatalf("TestGetTotalFKViolators, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name string
		fk   ForeignKey
		want int
	}{
		{"fk_violation_test_total", ForeignKey{fakeFkViolationTableName, fakeFkColumnName,
			fakePkViolationTableName, fakePkColumnName}, 5},
		{"fk_non_violation_test_total", ForeignKey{fakeFkNonViolationTableName, fakeFkColumnName,
			fakePkNonViolationTableName, fakePkColumnName}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTotalFKViolators(tt.fk); got != tt.want {
				t.Errorf("TestGetTotalFKViolators = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: TotalRows, Get total number of rows in the table
func TestTotalRows(t *testing.T) {
	tests := []struct {
		name  string
		table string
		want  int
	}{
		{"fk_violation_table", fakeFkViolationTableName, 5},
		{"fk_non_violation_table", fakeFkNonViolationTableName, 3},
		{"pk_violation_table", fakePkViolationTableName, 5},
		{"pk_non_violation_table", fakePkNonViolationTableName, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TotalRows(tt.table); got != tt.want {
				t.Errorf("TestTotalRows = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: GetFKViolators, get the total FK violators
func TestGetFKViolators(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakeSql)
	if err != nil {
		Fatalf("TestGetFKViolators, error in executing the statement, err: %v", err)
	}
	tests := []struct {
		name string
		fk   ForeignKey
		want int
	}{
		{"fk_violation_test_row", ForeignKey{fakeFkViolationTableName, fakeFkColumnName,
			fakePkViolationTableName, fakePkColumnName}, 5},
		{"fk_non_violation_test_row", ForeignKey{fakeFkNonViolationTableName, fakeFkColumnName,
			fakePkNonViolationTableName, fakePkColumnName}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFKViolators(tt.fk); len(got) != tt.want {
				t.Errorf("TestGetFKViolators = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: UpdateFKeys, fix the foreign key violation
func TestUpdateFKeys(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakeSql)
	if err != nil {
		Fatalf("TestUpdateFKeys, error in executing the statement, err: %v", err)
	}
	fk := ForeignKey{fakeFkViolationTableName, fakeFkColumnName,
		fakePkViolationTableName, fakePkColumnName}
	totalRowsOnPkTable := TotalRows(fk.Reftable)
	for GetTotalFKViolators(fk) > 0 {
		ViolationRows := GetFKViolators(fk)
		for _, v := range ViolationRows {
			UpdateFKeys(fk, totalRowsOnPkTable, v.Row)
		}
	}
	t.Run("validate_all_fk_rows_are_fixed", func(t *testing.T) {
		if got := GetTotalFKViolators(fk); got != 0 {
			t.Errorf("TestUpdateFKeys = %v, want %v", got, 0)
		}
	})
}

// Test: deleteViolatingConstraintKeys, delete the pk, uk column that violates the constraint even after the fix
func TestDeleteViolatingConstraintKeys(t *testing.T) {
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fakePKSql)
	if err != nil {
		Fatalf("TestDeleteViolatingConstraintKeys, error in executing the statement, err: %v", err)
	}
	check := func(want int) {
		if got := getTotalPKViolator(fakePkViolationTableName, fakePkColumnName); got != want {
			t.Errorf("TestDeleteViolatingConstraintKeys = %v, want %v", got, want)
		}
	}
	t.Run("check_for_pk_violation_rows", func(t *testing.T) {
		check(2)
	})
	t.Run("cleanup_violation_rows_by_deleting_it", func(t *testing.T) {
		err := deleteViolatingConstraintKeys(fakePkViolationTableName, fakePkColumnName)
		if err != nil {
			t.Errorf("TestDeleteViolatingConstraintKeys received error from database, err: %v", err)
		}
		check(0)
	})
}

// Test: checkEnumDatatype, get the Emun datatype rows
func TestCheckEnumDatatype(t *testing.T) {
	t.Run("enum_datatype_data_validation", func(t *testing.T) {
		if got := checkEnumDatatype("rating"); len(got) != 3 {
			t.Errorf("TestCheckEnumDatatype = %v, want %v", len(got), 3)
		}
	})
}
