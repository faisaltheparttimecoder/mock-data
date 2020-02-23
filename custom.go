package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type Skeleton struct {
	Custom []TableModel `yaml:"Custom"`
}

type TableModel struct {
	Schema string        `yaml:"Schema"`
	Table  string        `yaml:"Table"`
	Column []ColumnModel `yaml:"Column"`
}

type ColumnModel struct {
	Name   string   `yaml:"Name"`
	Type   string   `yaml:"Type"`
	Random bool     `yaml:"Random"`
	Values []string `yaml:"Values"`
}

// Generate a YAML of the mock plan related to this table
func GenerateMockPlan() {
	Infof("Generating a skeleton YAML for the list of table provided")

	// Check if the tables provided exists
	whereClause := generateWhereClause()
	tableList := dbExtractTables(whereClause)

	// If there is any then extract the column and data type
	if len(tableList) > 0 {
		columns := columnExtractor(tableList)
		s := BuildSkeletonYaml(columns)
		RegisterSkeletonYamlToFile(s)
	} else {
		Warn("No table available to generate the table skeleton, closing the program")
	}
}

// Build the skeleton Yaml
func BuildSkeletonYaml(column []TableCollection) Skeleton {
	Debugf("Build skeleton yaml")
	var s Skeleton
	for _, v := range column {
		var t = &TableModel{
			Schema: v.Schema,
			Table:  v.Table,
		}
		for _, u := range v.Columns {
			var c = &ColumnModel{
				Name:   u.Column,
				Type:   u.Datatype,
				Random: true,
			}
			t.Column = append(t.Column, *c)
		}
		s.Custom = append(s.Custom, *t)
	}
	return s
}

// Create a yaml file and write the contents onto the file
func RegisterSkeletonYamlToFile(skeleton Skeleton) {
	Debugf("Save the yaml skeleton to file")
	f := fmt.Sprintf("%s_skeleton_%s.yaml", programName, ExecutionTimestamp)

	// Marshal the content
	d, err := yaml.Marshal(&skeleton)
	if err != nil {
		Fatalf("Error when marshalling the yaml file, err: %v", err)
	}

	// Save to file
	err = WriteToFile(f, string(d))
	if err != nil {
		Fatalf("Error when saving the yaml to file %s, err: %v", f, err)
	}

	Infof("The YAML is saved to file: %s/%s", CurrentDir(), f)
}

// Load the custom configuration and mock data based on that configuration
func MockCustoms() {
	Infof("Loading the table using custom configuration")
	c := Skeleton{}

	// Read the configuration
	c.ReadConfiguration()

	// Start Loading the data
	c.LoadDataByConfiguration()

	// If the program skipped the tables lets the users know
	skipTablesWarning()
}

// Read Configuration and load it to skeleton struct
func (c *Skeleton) ReadConfiguration() {
	Debugf("Reading the configuration file: %s", cmdOptions.File)

	// Using viper reading the configuration
	viper.SetConfigFile(cmdOptions.File)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		Fatalf("Error reading the YAML file, err: %v", err)
	}

	// Load the configuration onto struct
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

// Load the data based on custom configuration
func (c *Skeleton) LoadDataByConfiguration() {
	Infof("Loading data to the table based on what is defined by file %s", cmdOptions.File)

	// Open db connection
	db := ConnectDB()
	defer db.Close()

	for _, s := range c.Custom {
		// Initialize the mocking process
		tab := GenerateTableName(s.Table, s.Schema)
		msg := fmt.Sprintf("Mocking Table %s", tab)
		bar := StartProgressBar(msg, cmdOptions.Rows)

		// Name the for loop to break when we encounter error
	DataTypePickerLoop:
		// Loop through the row count and start loading the data
		for i := 0; i < cmdOptions.Rows; i++ {
			var data []string
			var col []string

			// Column info
			for _, v := range s.Column {
				var d interface{}
				var err error
				if v.Random { // If the user said for this column choose anything
					d, err = BuildData(v.Type)
					if err != nil {
						if strings.HasPrefix(fmt.Sprint(err), "unsupported datatypes found") {
							Debugf("Table %s skipped: %v", tab, err)
							skippedTab = append(skippedTab, tab)
							bar.Add(cmdOptions.Rows)
							break DataTypePickerLoop
						} else {
							Fatalf("Error when building data for table %s: %v", tab, err)
						}
					}
				} else { // User asked to use only the one that is provided
					if len(v.Values) > 0 {
						d = RandomPickerFromArray(v.Values)
					} else {
						Fatalf("Random is set to false for table %s column %s, but "+
							"no value is provided to custom fit, please check configuration file %s",
							tab, v.Name, cmdOptions.File)
					}
				}
				col = append(col, v.Name)
				data = append(data, fmt.Sprintf("%v", d))
			}

			// Copy the data to the table
			CopyData(tab, col, data, db)
			bar.Add(1)
		}
		fmt.Println()
	}
}
