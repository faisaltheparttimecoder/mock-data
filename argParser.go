package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ielizaga/mockd/core"
)

type PostgresConnector struct {
	Db, Username, Password, Table string
	Port, RowCount                int
}

var engines = []string{"postgres", "greenplum", "hawq"}
var Connector PostgresConnector

// Program Usage.
func ShowHelp() {
	fmt.Print(`
USAGE: mockd <DATABASE ENGINE> <OPTIONS>
DATABASE ENGINE:
	postgres        Postgres database
	greenplum       Greenplum database
	hdb             Hawq Database
	help            Show help
OPTIONS:
	Execute "mockd <database engine> -h" for all the database specific options

`)
	os.Exit(0)
}

// Parse the arguments passed via the OS command line
func ArgPaser() {

	// Download Command Parser
	postgresFlag := flag.NewFlagSet("postgres", flag.ExitOnError)
	postgresPortFlag := postgresFlag.Int("p", 5432, "The port that is used by the database engine")
	postgresDBFlag := postgresFlag.String("d", "postgres", "The database name where the table resides")
	postgresUsernameFlag := postgresFlag.String("u", *postgresDBFlag, "The username that can be used to connect to the database")
	postgresPasswordFlag := postgresFlag.String("w", "", "The password for the user that can be used to connect to the database")
	postgresTotalRowsFlag := postgresFlag.Int("n", 1, "The total number of mocked rows that is needed")
	postgresTableFlag := postgresFlag.String("t", "", "The table name to be filled in with mock data")
	flag.Parse()

	// If there is a command keyword provided then check to what is it and then parse the appropriate options
	var engine_arg = os.Args[1]

	switch {
	case core.StringContains(engine_arg, engines): // Postgres command parser
		postgresFlag.Parse(os.Args[2:])
	case engine_arg == "help": // Help Menu
		ShowHelp()
	default: // Error if command is invalid
		fmt.Printf("ERROR: %q is not valid command.\n", os.Args[1])
		ShowHelp()
	}

	// Parse the command line argument
	// Postgres database engine
	if postgresFlag.Parsed() {
		DBEngine = "postgres"
		Connector.Db = *postgresDBFlag
		Connector.Username = *postgresUsernameFlag
		Connector.Password = *postgresPasswordFlag
		Connector.Table = *postgresTableFlag
		Connector.Port = *postgresPortFlag
		Connector.RowCount = *postgresTotalRowsFlag
	}
}
