package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pivotal/mock-data/core"
	"github.com/pivotal/mock-data/db/postgres"
)

// Global Variables
var (
	skippedTab []string
	db         *sql.DB
	stmt       *sql.Stmt
)

// Progress Database connection
func dbConn() error {
	dbconn, err := sql.Open(DBEngine, fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", Connector.Username, Connector.Password, Connector.Host, Connector.Port, Connector.Db))
	if err != nil {
		return fmt.Errorf("Cannot establish a database connection: %v\n", err)
	}
	db = dbconn
	return nil
}

// Check if we can run the query and extract the version of the database
func dbVersion() error {

	log.Info("Obtaining the version of the database")
	var version string

	// Obtain the version of the database
	rows, err := db.Query(postgres.PGVersion())
	if err != nil {
		return fmt.Errorf("Cannot extracting version, error from the database: %v", err)
	}

	// Store the information of the version onto a variable
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return fmt.Errorf("Error scanning the rows from the version query: %v", err)
		}
	}

	// Print the version of the database on the logs
	log.Infof("Version of the database: %v", version)

	return nil

}

// Extract all the tables in the database
func dbExtractTables() ([]string, error) {

	log.Info("Extracting all the tables in the database")
	var tableString []string
	var rows *sql.Rows
	var err error

	// Obtain all the tables in the database
	if Connector.Engine == "postgres" { // Use postgres specific query
		rows, err = db.Query(postgres.PGAllTablesQry1())
	} else { // Use greenplum, hdb query to extract the columns
		rows, err = db.Query(postgres.PGAllTablesQry2())
	}

	if err != nil {
		return tableString, fmt.Errorf("Cannot extract all the tables, error from the database: %v", err)
	}

	// Loop through the rows and store the table names.
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return tableString, fmt.Errorf("Error extracting the rows of the list of tables: %v", err)
		}
		tableString = append(tableString, table)
	}

	return tableString, nil
}

// Get all the columns and its datatype from the query
func dbColDataType() ([]Table, error) {

	log.Info("Checking for the existence of the table provided to the application, if exist extract all the column and datatype information")
	var table []Table
	var rows *sql.Rows
	var err error

	// Loop through the table list provided and collect the columns and datatypes
	for _, v := range strings.Split(Connector.Table, ",") {
		var tab Table
		if DBEngine == "postgres" { // Use postgres specific query
			rows, err = db.Query(postgres.PGColumnQry1(v))
		} else { // Use greenplum, hdb query to extract the columns
			rows, err = db.Query(postgres.PGColumnQry2(v))
		}
		if err != nil {
			return table, fmt.Errorf("Cannot extracting the column info, error from the database: %v", err)
		}
		for rows.Next() {

			var col string
			var datatype string
			var seqCol string = ""

			// Scan and store the rows
			err = rows.Scan(&col, &datatype, &seqCol)
			if err != nil {
				return table, fmt.Errorf("Error extracting the rows of the list of columns: %v", err)
			}

			// Ignore columns with sequence, since its auto loaded no need to randomize
			if !strings.HasPrefix(seqCol, "nextval") {
				tab.tabname = v
				if tab.columns == nil {
					tab.columns = make(map[string]string)
				}
				tab.columns[col] = datatype
			}
		}

		// If there is no columns, then ignore that table
		if len(tab.columns) > 0 {
			table = append(table, tab)
		}

	}

	return table, nil
}

// Extract the table & columns and request to load data
func extractor(table_info []Table) error {

	// Before we begin lets take a backup of all the PK, UK, FK, CK ( unless user says to ignore it )
	// constraints since we are not sure when we send cascade to constraints
	// what all constraints are dropped. so its easy to take a backup of all
	// constraints and then execute this DDL script at the end after we fix all the
	// constraint issues.
	// THEORY: already exists would fail and not available would be created.
	if !Connector.IgnoreConstraints {
		log.Info("Backup up all the constraint in the database")
		err := postgres.BackupDDL(db, ExecutionTimestamp)
		if err != nil {
			return err
		}
	}

	// Loop through all the tables available and start to load data
	// based on columns datatypes
	log.Info("Separating the input to tables, columns & datatypes and attempting to mock data to the table")
	for _, v := range table_info {
		err := splitter(v.columns, v.tabname)
		if err != nil {
			return err
		}
	}

	return nil
}

// Segregate tables, columns & datatypes to load data
func splitter(columns map[string]string, tabname string) error {

	var schema string
	var colkey, coldatatypes []string

	// Collect the column and datatypes
	for key, dt := range columns {
		colkey = append(colkey, key)
		coldatatypes = append(coldatatypes, dt)
	}

	// Ensure all the constriants are removed from the table
	// and also store them to ensure all the constraints conditions
	// are met when we re-enable them
	err := postgres.RemoveConstraints(db, tabname)
	if err != nil {
		return err
	}

	// Split the table into schema and tablename
	tab := strings.Split(tabname, ".")
	if len(tab) == 1 { // if no schema provide then use the default postgres schema "public"
		schema = "public"
	} else { // else what is provided by the user
		schema = tab[0]
		tabname = tab[1]
	}

	// Start the progress bar
	progressMsg := "(Mocking Table: " + schema + "." + tabname + ")"
	core.ProgressBar(Connector.RowCount, progressMsg)

	// Commit the data to the database
	err = commitData(schema, tabname, colkey, coldatatypes)
	if err != nil {
		return err
	}

	// Close the Progress bar
	core.CloseProgressBar()

	return nil
}

// Start a transaction block and commit the data
func commitData(schema, tabname string, colkey, dtkeys []string) error {

	// Start a transaction
	txn, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error in starting a transaction: %v", err)
	}

	// Prepare the copy statement
	stmt, err = txn.Prepare(pq.CopyInSchema(schema, tabname, colkey...))
	if err != nil {
		return fmt.Errorf("Error in preparing the transaction statement: %v", err)
	}

	// Iterate through connector row count and build data for each datatype
DataTypePickerLoop: // Label the loop to break, if there is a datatype that we don't support
	for i := 0; i < Connector.RowCount; i++ {

		// data collector
		var data []interface{}

		// Generate data based on the columns datatype
		for _, v := range dtkeys {
			dataoutput, err := core.BuildData(v)
			if err != nil {
				if strings.HasPrefix(fmt.Sprint(err), "Unsupported datatypes found") {
					log.Errorf("Skipping table \"%s\" due to error \"%v\"", tabname, err)
					skippedTab = append(skippedTab, tabname)
					break DataTypePickerLoop // break the loop
				} else {
					return err
				}

			}
			data = append(data, dataoutput)
		}

		// Execute the statement
		_, err = stmt.Exec(data...)
		if err != nil {
			return err
		}

		// Increment progress bar
		core.IncrementBar()
	}

	// Close the statement
	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("Error in closing the transaction statement: %v", err)
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("Error in committing the transaction statement: %v", err)
	}

	return nil

}

// Main postgres data mocker
func MockPostgres() error {

	var table []Table
	log.Infof("Attempting to establish a connection to the %s database", DBEngine)

	// Establishing a connection to the database
	err := dbConn()
	if err != nil {
		return err
	}

	// Check if we can query the database and get the version of the database in the meantime
	err = dbVersion()
	if err != nil {
		return err
	}

	// If the request is to load all table then, extract all tables
	// and pass to the connector table argument.
	if Connector.AllTables {
		tableList, err := dbExtractTables()
		if err != nil {
			return err
		}
		Connector.Table = strings.Join(tableList, ",")
	}

	// Extract the columns and datatypes from the table defined on the connector table.
	if Connector.Table != "" { // if there are only tables in the connector table variables
		table, err = dbColDataType()
		if err != nil {
			return err
		}
	}

	// Build data for all the column and datatypes & then commit data
	if len(table) > 0 { // if there are tables found, then proceed
		err = extractor(table)
		if err != nil {
			// TODO: need to fix constraints here as well.
			log.Error("Unexpected error encountered by MockD..")
			return err
		}

		// Recreate all the constraints of the table unless user wants to ignore it
		if !Connector.IgnoreConstraints {
			err = postgres.FixConstraints(db, ExecutionTimestamp)
			if err != nil {
				backupFiles, _ := core.ListFile(".", "*_"+ExecutionTimestamp+".sql")
				log.Criticalf("Constraints creation failed, all the DDL are saved in the files: \n%v", strings.Join(backupFiles, "\n"))
				log.Criticalf("Will need your intervention to fix constraints")
				return err
			}
		}

	} else { // We didn't obtain any table from the database ( eg.s fresh DB's or User gave a view name etc )
		log.Warning("No table's available to load the mock data, closing the program")
	}

	// If there is tables that are skipped, report to the user.
	if len(skippedTab) > 0 {
		log.Warning("These tables (below) are skipped, since it contain unsupported datatypes")
		log.Warningf("%s", strings.Join(skippedTab, ","))
	}

	// Close the database connection
	defer db.Close()

	return nil
}
