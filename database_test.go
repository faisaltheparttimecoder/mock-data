package main

import (
	"github.com/spf13/viper"
	"testing"
)

func setDatabaseConfigForTest() {
	cmdOptions.Database = viper.GetString("PGDATABASE")
	cmdOptions.Password = viper.GetString("PGPASSWORD")
	cmdOptions.Username = viper.GetString("PGUSER")
	cmdOptions.Hostname = viper.GetString("PGHOST")
	cmdOptions.Port = viper.GetInt("PGPORT")
}

// Test: ExecuteDemoDatabase, check if all the tables are created
func TestExecuteDemoDatabase(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("should_extract_all_tables", func(t *testing.T) {
		postgresOrGreenplum()
		ExecuteDemoDatabase()
		if got := allTablesPostgres(""); len(got) != 15 {
			t.Errorf("TestExecuteDemoDatabase = %v, want %d tables", len(got), 15)
		}
	})
}

// Test: MockDatabase, Mock the entire database
func TestMockDatabase(t *testing.T) {
	setDatabaseConfigForTest()
	cmdOptions.Rows = 100
	MockDatabase()
	tests := []struct {
		name  string
		table string
	}{
		{"mock_database_check_actor_table", "public.actor"},
		{"mock_database_check_address_table", "public.address"},
		{"mock_database_check_category_table", "public.category"},
		{"mock_database_check_city_table", "public.city"},
		{"mock_database_check_country_table", "public.country"},
		{"mock_database_check_customer_table", "public.customer"},
		{"mock_database_check_film_table", "public.film"},
		{"mock_database_check_film_actor_table", "public.film_actor"},
		{"mock_database_check_film_category_table", "public.film_category"},
		{"mock_database_check_inventory_table", "public.inventory"},
		{"mock_database_check_language_table", "public.language"},
		{"mock_database_check_payment_table", "public.payment"},
		{"mock_database_check_rental_table", "public.rental"},
		{"mock_database_check_staff_table", "public.staff"},
		{"mock_database_check_store_table", "public.store"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Why > 0, since we do delete some rows with deleteViolatingConstraintKeys, so we dont know
			// how much we deleted, so lets just say it should have greater than 0 rows
			if got := TotalRows(tt.table); got <= 0 {
				t.Errorf("TestMockDatabase = %v, want %v", got, "> 0")
			}
		})
	}
}

// Test: dbExtractTables, check if all the tables are returned
func TestDbExtractTables(t *testing.T) {
	setDatabaseConfigForTest()
	t.Run("should_return_all_tables", func(t *testing.T) {
		if got := dbExtractTables(""); len(got) != 15 {
			t.Errorf("TestExecuteDemoDatabase = %v, want %d tables", len(got), 15)
		}
	})
}