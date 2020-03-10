package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"strings"
)

type TableCollection struct {
	DBTables
	Columns []DBColumns
}

var (
	skippedTab []string
	delimiter  = "$"
)

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
			columns = columnExtractorPostgres(fmt.Sprintf("\"%s\"", t.Schema), t.Table)
		} else {
			columns = columnExtractorGPDB(fmt.Sprintf("\"%s\"", t.Schema), t.Table)
		}

		// There are instance where the table can have one column and datatype serial
		// take a look at the issue: https://github.com/pivotal-gss/mock-data/issues/29
		// lets fix this
		if len(columns) == 1 {
			checkAndAddDataIfItsASerialDatatype(t, columns)
		}

		// Loops through the columns and make a collection of tables
		// & column, we ignore sequence since they are auto injected also
		for _, c := range columns {
			if !isItSerialDatatype(c) {
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
	skipTablesWarning()

	Infof("Completed loading mock data to %d tables", totalTables)
}

// Start Committing data to the database
func CommitData(t TableCollection) {
	// Start committing data
	tab := GenerateTableName(t.Table, t.Schema)
	msg := fmt.Sprintf("Mocking Table %s", tab)
	bar := StartProgressBar(msg, cmdOptions.Rows)

	// Open db connection
	db := ConnectDB()
	defer db.Close()

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
					Debugf("Table %s skipped, since the column %s, had unknown data type %s: %v",
						tab, c.Column, c.Datatype, err)
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

		// Copy the data to the table
		CopyData(tab, col, data, db)
		bar.Add(1)
	}
	fmt.Println()
}

// Copy the data to the database table
func CopyData(tab string, col []string, data []string, db *pg.DB) {
	// Copy Statement and start loading
	copyStatment := fmt.Sprintf(`COPY %s("%s") FROM STDIN WITH CSV DELIMITER '%s' QUOTE e'\x01'`,
		tab, strings.Join(col, "\",\""), delimiter)
	_, err := db.CopyFrom(strings.NewReader(strings.Join(data, delimiter)), copyStatment)

	// Handle Error
	if err != nil {
		fmt.Println()
		Debugf("Table: %s", tab)
		Debugf("Copy Statement: %s", copyStatment)
		Debugf("Data: %s", strings.Join(data, delimiter))
		Fatalf("Error during committing data: %v", err)
	}
}

// Insert data to the table if its only a single column with serial data type
func checkAndAddDataIfItsASerialDatatype(t DBTables, c []DBColumns) {
	Debugf("Check if the table %s.%s which has only a single column is of serial data type", t.Schema, t.Table)
	column := c[0] // we know its only one , because we did a check on the parent function
	total := 0
	if isItSerialDatatype(column) {
		for total < cmdOptions.Rows {
			query := "INSERT INTO \"%s\".\"%s\" default values;"
			query = fmt.Sprintf(query, t.Schema, t.Table)
			_, err := ExecuteDB(query)
			if err != nil {
				Fatalf("Error when loading the serial datatype for table %s.%s, err: %v",
					t.Schema, t.Table, err)
			}
			total++
		}
	}
}

// Is it serial
func isItSerialDatatype(c DBColumns) bool {
	if strings.HasPrefix(c.Sequence, "nextval") {
		return true
	}
	return false
}

// Generate table name
func GenerateTableName(tab, schema string) string {
	return fmt.Sprintf("\"%s\".\"%s\"", schema, tab)
}

// Throw warning if there is skipped tables
func skipTablesWarning() {
	if len(skippedTab) > 0 {
		Warnf("These tables are skipped since these data types are not supported by %s: %s",
			programName, strings.Join(skippedTab, ","))
	}
}
