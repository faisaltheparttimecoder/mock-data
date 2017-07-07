package main

import (
	"os"

	"github.com/ielizaga/mockd/core"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
)

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
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	// Parse the arguments that has been passed on to the OS
	log.Info("Parsing all the command line arguments")
	ArgPaser()

	// What is the database engine that needs to be used
	if core.StringContains(DBEngine, engines) {
		err := MockPostgres()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	} else { // Unsupported database engine
		log.Errorf("mockd application doesn't support the database: %s", DBEngine)
		os.Exit(1)
	}

	log.Info("mockd program has successfully completed")

}
