package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"github.com/pivotal/mock-data/core"
)

// Version of Mock-data
var version = "1.3"

// All global variables
var (
	DBEngine string
)

// Define the logging format, used in the project
var (
	log    = logging.MustGetLogger("mockd")
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000}:%{level:s} > %{color:reset}%{message}`,
	)
)

// file timestamp
var ExecutionTimestamp = core.TimeNow()

// An Engine is an implementation of a database
// engine like PostgreSQL, MySQL or Greenplum
type Engine struct {
	name, version string
	port          int
}

// A Table is an implementation of a database with a set of columns and datatypes
type Table struct {
	tabname string
	columns map[string]string
}

// Main block
func main() {

	// Logger for go-logging package
	// create backend for os.Stderr, set the format and update the logger to what logger to be used
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	// Parse the arguments that has been passed on to the OS
	ArgPaser()

	// This execution timestamp
	log.Infof("Timestamp of this mockd execution: %s", ExecutionTimestamp)

	// What is the database engine that needs to be used
	// call the appropriate program that is specific to database engine
	if DBEngine == "postgres" {
		err := MockPostgres()
		if err != nil {
			log.Error(err)
			log.Info("mockd program has completed with errors")
			os.Exit(1)
		}
	} else { // Unsupported database engine.
		log.Errorf("mockd application doesn't support the database: %s", DBEngine)
		os.Exit(1)
	}

	log.Info("mockd program has successfully completed")

}
