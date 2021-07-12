package main

import (
	"fmt"
	"testing"
)


// create fake tables for constraint backup
func createFakeTablesForConstraintBackup() {
	setDatabaseConfigForTest()
	postgresOrGreenplum()
	ExecuteDemoDatabase()
}

// Test: BackupDDL
func TestBackupDDL(t *testing.T) {
	// Skipping since this function basically calls other function, without doing anything
}

// Test: backupConstraints, check if the function backed up the constraint and wrote to a file
func TestBackupConstraints(t *testing.T) {
	createFakeTablesForConstraintBackup()
	CreateDirectory()
	backupConstraints()
	tests := []struct {
		name         string
		connKey      string
		shortConnKey string
		want         int
	}{
		{"checking_check_key_count", "CHECK", "c", 1},
		{"checking_foreign_key_count", "FOREIGN", "f", 18},
		{"checking_primary_key_count", "PRIMARY", "p", 15},
		{"checking_check_key_count", "UNIQUE", "u", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(savedConstraints[tt.connKey]); got != tt.want {
				t.Errorf("TestBackupConstraints in memory = %v, want %v", got, tt.want)
			}
			f := fmt.Sprintf("%s/%s_constraint_backup_%s.sql", Path, programName, tt.shortConnKey)
			content, _ := ReadFile(f)
			if got := len(content); got != tt.want {
				t.Errorf("TestBackupConstraints in file = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: backupIndexes, check if the function backed up the unique indexes and wrote to a file
func TestBackupIndexes(t *testing.T) {
	createFakeTablesForConstraintBackup()
	backupIndexes()
	want := 19
	t.Run("checking_unique_index_count", func(t *testing.T) {
		if got := len(savedConstraints["UNIQUE"]); got != want {
			t.Errorf("TestBackupIndexes in memory = %v, want %v", got, want)
		}
		f := fmt.Sprintf("%s/%s_constraint_backup_%s.sql", Path, programName, "u")
		content, _ := ReadFile(f)
		if got := len(content); got != want {
			t.Errorf("TestBackupIndexes in file = %v, want %v", got, want)
		}
	})
}

// Test: constraintFinder, find the postgres name for the shorthand constraint names
func TestConstraintFinder(t *testing.T) {
	tests := []struct {
		name          string
		shortConnName string
		want          string
	}{
		{"hunting_check_constraint", "c", "CHECK"},
		{"hunting_foreign_key_constraint", "f", "FOREIGN"},
		{"hunting_primary_key_constraint", "p", "PRIMARY"},
		{"hunting_check_key_constraint", "u", "UNIQUE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := constraintFinder(tt.shortConnName); got != tt.want {
				t.Errorf("TestConstraintFinder = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: RemoveConstraints, test if all the constraints are removed without any issues
func TestRemoveConstraints(t *testing.T) {
	tests := []struct {
		name  string
		table string
		want  int
	}{
		{"constraint_actor_table", "public.actor", 0},
		{"constraint_address_table", "public.address", 0},
		{"constraint_category_table", "public.category", 0},
		{"constraint_city_table", "public.city", 0},
		{"constraint_country_table", "public.country", 0},
		{"constraint_customer_table", "public.customer", 0},
		{"constraint_film_table", "public.film", 0},
		{"constraint_film_actor_table", "public.film_actor", 0},
		{"constraint_film_category_table", "public.film_category", 0},
		{"constraint_inventory_table", "public.inventory", 0},
		{"constraint_language_table", "public.language", 0},
		{"constraint_payment_table", "public.payment", 0},
		{"constraint_rental_table", "public.rental", 0},
		{"constraint_staff_table", "public.staff", 0},
		{"constraint_store_table", "public.store", 0},
	}
	setDatabaseConfigForTest()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveConstraints(tt.table)
			if got := len(GetConstraintsPertab(tt.table)); got != tt.want {
				t.Errorf("TestRemoveConstraints = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: saveConstraints, check if the contents are being saved in memory and to the file
func TestSaveConstraints(t *testing.T) {
	setDatabaseConfigForTest()
	fakeConstraint := "testsuite"
	fakeTable := "fakeTable"
	fakeColumn := "fakeColumn"
	f := fmt.Sprintf("%s/%s_constraint_backup_%s.sql", Path, programName, fakeConstraint)
	iteration := 10
	for i := 0; i < iteration; i++ {
		saveConstraints(f, "hello man", fakeConstraint, fakeTable, fakeColumn)
	}
	content, _ := ReadFile(f)
	t.Run("check_value_saved_in_memory", func(t *testing.T) {
		if got := savedConstraints[fakeConstraint]; len(got) != iteration {
			t.Errorf("TestSaveConstraints in memory = %v, want %v", len(got), iteration)
		}
	})
	t.Run("check_value_saved_in_file", func(t *testing.T) {
		if got := len(content); got != 1 { // since there is no new line
			t.Errorf("TestSaveConstraints in file = %v, want %v", got, 1)
		}
	})
}
