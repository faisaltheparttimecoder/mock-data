package main

import (
	"fmt"
	"strings"
)

type TableCollection struct {
	DBTables
	Columns []DBColumns
}

var skippedTab []string

func MockTable(tables []DBTables) {
	// Check if there is any rows on the table list, if yes then start
	// the loading process
	totalTables := len(tables)
	if totalTables > 0 {
		Debugf("Total number of tables to mock: %d", totalTables)
		tableMocker(tables)
		if !cmdOptions.IgnoreConstraint {
			FixConstraints()
		}
	} else { // no tables found, explain that to the user and exit
		Warn("No table available to mock the data, closing the program")
	}
}

// Extract the column & Start the table mocking process
func tableMocker(tables []DBTables) {
	Info("Beginning the mocking process for the tables")

	// Before beginning the process, recheck with the user
	// they still want to continue
	if !cmdOptions.DontPrompt {
		_ = YesOrNoConfirmation()
	}

	// User confirmed to continue, first extract the column
	// and its data types
	columns := columnExtractor(tables)

	// If there is some tables in the list, then go through the
	// next step, else print warning for the users
	if len(columns) > 0 {
		BackupConstraintsAndStartDataLoading(columns)
	} else { // no tables
		Warn("No columns available to mock the data, closing the program")
	}
}

// Extract the column and its datatypes of the table
func columnExtractor(tables []DBTables) []TableCollection {
	Info("Extracting the columns and datatype information")
	var columns []DBColumns
	var collection []TableCollection

	// Start Progress bar
	bar := StartProgressBar("Extracting column information from tables", len(tables))

	for _, t := range tables {
		var tempColumns []DBColumns
		if GreenplumOrPostgres == "postgres" {
			columns = columnExtractorPostgres(t.Schema, t.Table)
		} else {
			columns = columnExtractorGPDB(t.Schema, t.Table);
		}

		// Loops through the columns and make a collection of tables
		// & column, we ignore sequence since they are auto injected also
		for _, c := range columns {
			if !strings.HasPrefix(c.Sequence, "nextval") {
				tempColumns = append(tempColumns, c)
			}
		}

		// ignore the table, that doesn't have columns
		if len(tempColumns) > 0 {
			collection = append(collection, TableCollection{t, tempColumns})
		}
		bar.Add(1)
	}
	fmt.Println()
	return collection
}

// Backup and start the loading process
func BackupConstraintsAndStartDataLoading(tables []TableCollection) {
	// Backup the DDL first
	BackupDDL()
	// Loop through the tables, splits the tables in schema
	// & table and start loading
	totalTables := len(tables)
	Infof("Total numbers of tables to mock: %d", totalTables)
	for _, t := range tables {
		// Remove Constraints
		table := fmt.Sprintf("\"%s\".\"%s\"", t.Schema, t.Table)
		RemoveConstraints(table)

		// Start the committing data to the table
		CommitData(t)
	}

	// If the program skipped the tables lets the users know
	if len(skippedTab) > 0 {
		Warnf("These tables are skipped since these datatypes are not supported by %s:%s",
			programName, strings.Join(skippedTab, ","))
	}
	Infof("Completed loading mock data to %d tables", totalTables)
}

// Start Committing data to the database
func CommitData(t TableCollection) {
	// Start committing data
	tab := fmt.Sprintf("\"%s\".\"%s\"", t.Schema, t.Table)
	msg := fmt.Sprintf("Mocking Table %s", tab)
	bar := StartProgressBar(msg, cmdOptions.Rows)

	// db connection
	db := ConnectDB()
	defer db.Close()

	// Delimiter
	delimiter := "$"

	// Name the for loop to break when we encounter error
DataTypePickerLoop:
	// Loop through the row count and start loading the data
	for i := 0; i < cmdOptions.Rows; i++ {
		var data []string
		var col []string

		// Column info
		for _, c := range t.Columns {
			d, err := BuildData(c.Datatype)
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
			col = append(col, c.Column)
			data = append(data, fmt.Sprintf("%v", d))
		}

		// Copy Statement and start loading
		copyStatment := fmt.Sprintf(`COPY %s(%s) FROM STDIN WITH CSV DELIMITER '%s' QUOTE e'\x01'`,
			tab, strings.Join(col, ","), delimiter)
		_, err := db.CopyFrom(strings.NewReader(strings.Join(data, delimiter)), copyStatment)

		// Handle Error
		if err != nil {
			fmt.Println()
			Debugf("Table: %s", tab)
			Debugf("Copy Statement: %s", copyStatment)
			Debugf("Data: %s", strings.Join(data, delimiter))
			Fatalf("Error during committing data: %v", err)
		}
		bar.Add(1)
	}
	fmt.Println()
}
