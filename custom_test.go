package main

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

const (
	customTable = "custom_table"
)

var YamlDataForTest = `Custom:
- Schema: mock_data_table5
  Table: custom_table
  Column:
  - Name: id
    Type: integer
    UserData: []
    Realistic: ""
  - Name: name
    Type: character varying(100)
    UserData: []
    Realistic: "NameFullName"
  - Name: gender
    Type: character varying
    UserData:
      - "M"
      - "F"
      - "O"
    Realistic: ""
  - Name: salary
    Type: money
    UserData: []
    Realistic: ""
`

// create fake tables form custom.go
func createFakeTablesFromCustom() {
	setDatabaseConfigForTest()
	cmdOptions.Tab.SchemaName = "mock_data_table5"
	cmdOptions.DontPrompt = true
	cmdOptions.Rows = 30
	cmdOptions.Tab.FakeTablesRows = customTable
	postgresOrGreenplum()
	sql := `
		DROP SCHEMA IF EXISTS %[1]s CASCADE;
		CREATE SCHEMA %[1]s;
		DROP TABLE IF EXISTS %[1]s.%[2]s;
		CREATE TABLE %[1]s.%[2]s (id int primary key, name varchar(100), gender varchar, salary money);
	`
	sql = fmt.Sprintf(sql, cmdOptions.Tab.SchemaName, customTable)
	_, err := ExecuteDB(sql)
	if err != nil {
		Fatalf("createFakeTablesFromCustom, error in executing the statement, err: %v", err)
	}
}

// Load and return the struct
func (c *Skeleton) LoadYamlForCustomTest(f string) error {
	viper.SetConfigFile(f)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading the YAML file: %v", err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %v", err)
	}
	return nil
}

// Test: GenerateMockPlan, check if the function generates a plan
func TestGenerateMockPlan(t *testing.T) {
	createFakeTablesFromCustom()
	f := fmt.Sprintf("%s_skeleton_%s.yaml", programName, ExecutionTimestamp)
	t.Run("should_generate_a_plan_file", func(t *testing.T) {
		GenerateMockPlan()
		if _, err := ReadFile(f); err != nil {
			t.Errorf("TestGenerateMockPlan, should have generated a yaml, but got err: %v", err)
		}
	})
	t.Run("should_generate_a_valid_file", func(t *testing.T) {
		s := new(Skeleton)
		err := s.LoadYamlForCustomTest(f)
		if err != nil {
			t.Errorf("TestGenerateMockPlan: %v", err)
		}
	})
}

// Test: BuildSkeletonYaml, check if the skeleton generated matches
func TestBuildSkeletonYaml(t *testing.T) {
	createFakeTablesFromCustom()
	notBlankString := func(k, s string) {
		t.Run("should_have_some_value_for_the_key_"+k, func(t *testing.T) {
			if IsStringEmpty(s) {
				t.Errorf("TestBuildSkeletonYaml, Key \"%s\" should contain a value, got it empty", k)
			}
		})
	}
	whereClause := generateWhereClause()
	tableList := dbExtractTables(whereClause)
	columns := columnExtractor(tableList)
	s := BuildSkeletonYaml(columns)
	for _, v := range s.Custom {
		notBlankString("Table", v.Table)
		notBlankString("Schema", v.Table)
		for _, c := range v.Column {
			notBlankString("Column.Name", c.Name)
			notBlankString("Column.Type", c.Type)
			t.Run("should_be_zero_for_the_key_Column.UserData", func(t *testing.T) {
				if len(c.UserData) != 0 {
					t.Errorf("TestBuildSkeletonYaml, Key \"%s\" should be empty", "Column.UserData")
				}
			})
			t.Run("should_be_empty_for_the_key_Column.Realistic", func(t *testing.T) {
				if !IsStringEmpty(c.Realistic) {
					t.Errorf("TestBuildSkeletonYaml, Key \"%s\" should not have a value", "Column.Realistic")
				}
			})
		}
	}
}

// Test: RegisterSkeletonYamlToFile, should successfully record the data to file
func TestRegisterSkeletonYamlToFile(t *testing.T) {
	// I believe this can be skipped since this has been tested with
	// TestGenerateMockPlan, skipped
}

// Test: MockCustoms, check if the function successfully loads data based on Yaml
func TestMockCustoms(t *testing.T) {
	// Its basically calls other function and those functions will be tested
	// individually , skipping ...
}

// Test: ReadConfiguration, check if its a valid struct
func TestReadConfiguration(t *testing.T) {
	// Already tested this via TestGenerateMockPlan, skipping
}

// Test: TestRealisticDataBuilder, check if all the fake key maps, return some data
func TestRealisticDataBuilder(t *testing.T) {
	// Get all the keys
	keys := reflect.ValueOf(fakeRealisticMaps()).MapKeys()
	strKeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strKeys[i] = keys[i].String()
	}
	for _, tt := range strKeys {
		t.Run("check_data_for_key_"+tt, func(t *testing.T) {
			if got := RealisticDataBuilder(tt); IsStringEmpty(fmt.Sprintf("%v", got)) {
				t.Errorf("TestRealisticDataBuilder = %v is detected returning empty, need some data", tt)
			}
		})
	}
}

// Test: LoadDataByConfiguration, check if the data is loaded based on the configuration file
func TestLoadDataByConfiguration(t *testing.T) {
	createFakeTablesFromCustom()
	s := new(Skeleton)
	cmdOptions.File = fmt.Sprintf("%s_skeleton_%s_test.yaml", programName, ExecutionTimestamp)
	err := WriteToFile(cmdOptions.File, YamlDataForTest)
	if err != nil {
		t.Errorf("TestLoadDataByConfiguration, failed to write to file, err: %v", err)
	}
	s.ReadConfiguration()
	s.LoadDataByConfiguration()
	t.Run("should_show_rows_in_the_table", func(t *testing.T) {
		tab := GenerateTableName(customTable, cmdOptions.Tab.SchemaName)
		if got := TotalRows(tab); got != cmdOptions.Rows {
			t.Errorf("TestLoadDataByConfiguration = %v, want %v", got, cmdOptions.Rows)
		}
	})
}
