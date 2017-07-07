package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ielizaga/mockd/core"
)

type connector struct {
	Db, Username, Password, Table string
	HostName, RegionName          string
	Port, RowCount                int
}

var engines = []string{"postgres", "greenplum", "hawq"}
var Connector connector

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

	// GemFire Command Parser
	gemfireFlag := flag.NewFlagSet("gemfire", flag.ExitOnError)
	gemfireHostNameFlag := gemfireFlag.String("h", "localhost", "The server host name.")
	gemfireRegionNameFlag := gemfireFlag.String("r", "replicatedRegion", "The region where the mock entries will be inserted.")
	gemfirePortFlag := gemfireFlag.Int("p", 8080, "The port that is used by the REST API endpoint.")
	gemfireTotalRowsFlag := gemfireFlag.Int("n", 1, "The total number of mocked objects that is needed.")
	flag.Parse()

	// If there is a command keyword provided then check to what is it and then parse the appropriate options
	var engine_arg = os.Args[1]

	switch {
	case core.StringContains(engine_arg, engines): // Postgres command parser
		postgresFlag.Parse(os.Args[2:])

	case engine_arg == "gemfire": // GemFire Command Parser
		gemfireFlag.Parse(os.Args[2:])

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

	if gemfireFlag.Parsed() {
		DBEngine = "gemfire"
		Connector.HostName = *gemfireHostNameFlag
		Connector.RegionName = *gemfireRegionNameFlag
		Connector.Port = *gemfirePortFlag
		Connector.RowCount = *gemfireTotalRowsFlag
	}
}
