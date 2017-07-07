package main

import (
	"database/sql"
	"fmt"
	"strings"

	"./core"
	"github.com/lib/pq"
)

var (
	db   *sql.DB
	stmt *sql.Stmt
)

// Postgres 9 and above
func PGColumnQry1(table string) string {
	return "SELECT   a.attname, " +
		"        pg_catalog.Format_type(a.atttypid, a.atttypmod), " +
		"	 COALESCE((SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128) " +
		"        FROM pg_catalog.pg_attrdef d " +
		"        WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef), '') " +
		"FROM     pg_catalog.pg_attribute a " +
		"WHERE    a.attrelid = '" + table + "'::regclass " +
		"AND      a.attnum > 0 " +
		"AND      NOT a.attisdropped " +
		"ORDER BY a.attnum "
}

// Postgres 8.3, GPDB, HDB
func PGColumnQry2(table string) string {
	return "SELECT         a.attname, " +
		"               pg_catalog.Format_type(a.atttypid, a.atttypmod), " +
		"	        COALESCE((SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128) " +
		"                FROM pg_catalog.pg_attrdef d " +
		"                WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef), '') " +
		"FROM            pg_catalog.pg_attribute a " +
		"LEFT OUTER JOIN pg_catalog.pg_attribute_encoding e " +
		"ON              e.attrelid = a .attrelid " +
		"AND             e.attnum = a.attnum " +
		"WHERE           a.attrelid = '" + table + "'::regclass " +
		"AND             a.attnum > 0 " +
		"AND             NOT a.attisdropped " +
		"ORDER BY        a.attnum"
}

// Postgres version
func PGVersion() string {
	return "select version()"
}

// Database connection
func dbConn() error {
	dbconn, err := sql.Open(DBEngine, fmt.Sprintf("user=%v dbname=%v sslmode=disable", Connector.Username, Connector.Db))
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
	rows, err := db.Query(PGVersion())
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
			rows, err = db.Query(PGColumnQry1(v))
		} else { // Use greenplum, hdb query to extract the columns
			rows, err = db.Query(PGColumnQry2(v))
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

	log.Info("Separating the input to tables, columns & datatypes and attempting to mock data to the table")
	for _, v := range table_info {
		err := splitter(v.columns, v.tabname)
		if err != nil {
			return err
		}
	}

	return nil
}

// Segregate tables and columns to load data
func splitter(columns map[string]string, tabname string) error {

	var schema string
	var colkey, coldatatypes []string

	// Collect the column and datatypes
	for key, dt := range columns {
		colkey = append(colkey, key)
		coldatatypes = append(coldatatypes, dt)
	}

	// Split the table into schema and tablename
	tab := strings.Split(tabname, ".")
	if len(tab) == 1 {
		schema = "public"
	} else {
		schema = tab[0]
		tabname = tab[1]
	}

	// Start the progress bar
	core.ProgressBar(Connector.RowCount, schema+"."+tabname)

	// Commit the data to the database
	err := commitData(schema, tabname, colkey, coldatatypes)
	if err != nil {
		return err
	}

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

	// Iterate through till the row the users want
	for i := 0; i < Connector.RowCount; i++ {

		// data collector
		var data []interface{}

		// Generate data based on the columns datatypes
		for _, v := range dtkeys {
			dataoutput, err := core.BuildData(v)
			if err != nil {
				return err
			}
			data = append(data, dataoutput)
		}

		// Execute the statement
		_, err = stmt.Exec(data...)
		if err, ok := err.(*pq.Error); ok {
			return fmt.Errorf("Error in executing the transaction statement: %v", err)
		}

		// IncrementBar
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

	// Establishing a connection to the database
	log.Infof("Attempting to establish a connection to the %s database", DBEngine)

	// Make a connection to the database
	err := dbConn()
	if err != nil {
		return err
	}

	// Check if we can query the database and get the version of the database in the meantime
	err = dbVersion()
	if err != nil {
		return err
	}

	// Get the columns and datatypes from the table
	table, err := dbColDataType()
	if err != nil {
		return err
	}

	// Build & commit data
	if len(table) > 0 {
		err = extractor(table)
		if err != nil {
			return err
		}
	} else {
		log.Warning("No table's available to load the mock data, closing the program")
	}

	// Close the database connection
	defer db.Close()

	return nil
}
