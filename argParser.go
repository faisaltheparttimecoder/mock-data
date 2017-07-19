package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/ielizaga/mockd/core"
	"strings"
)

// Connector struct
type connector struct {
	Engine string
	Db, Username, Password, Host, Table string
	Port, RowCount  int
	AllTables, IgnoreConstraints bool
}

// The connector
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

	// Postgres/Greenplum/Hawq(HDB) Command Parser
	postgresFlag := flag.NewFlagSet("postgres", flag.ExitOnError)
	postgresPortFlag := postgresFlag.Int("p", 5432, "The port that is used by the database engine")
	postgresDBFlag := postgresFlag.String("d", "postgres", "The database name where the table resides")
	postgresUsernameFlag := postgresFlag.String("u", *postgresDBFlag, "The username that can be used to connect to the database")
	postgresPasswordFlag := postgresFlag.String("w", "", "The password for the user that can be used to connect to the database")
	postgresHostFlag := postgresFlag.String("h", "localhost", "The hostname that can be used to connect to the database")
	postgresTotalRowsFlag := postgresFlag.Int("n", 1, "The total number of mocked rows that is needed")
	postgresTableFlag := postgresFlag.String("t", "", "The table name to be filled in with mock data")
	postgresAllDBFlag := postgresFlag.Bool("x", false, "Mock all the tables in the database")
	postgresIgnoreConstrFlag := postgresFlag.Bool("i", false, "Ignore checking and fixing constraint issues")
	flag.Parse()

	// Greenplum , HDB is built on top of postgres, so they will use the same Mock logic
	var engineArgs = os.Args[1]
	// Postgres
	var postgresEngines = []string{"postgres", "greenplum", "hawq"}

	// If there is a command keyword provided then check to what is it and then parse the appropriate options
	switch {
		// Postgres command parser
		case core.StringContains(engineArgs, postgresEngines):
			postgresFlag.Parse(os.Args[2:])
		// Help Menu
		case engineArgs == "help":
			ShowHelp()
		// If not of the list of supported engines, error out
		default:
			log.Errorf("%q is not valid command.", os.Args[1])
			ShowHelp()
	}

	// Parse the command line argument
	// Postgres database engine
	if postgresFlag.Parsed() {

		// Store all connector information
		DBEngine = "postgres"
		Connector.Engine = engineArgs
		Connector.Db = *postgresDBFlag
		Connector.Username = *postgresUsernameFlag
		Connector.Password = *postgresPasswordFlag
		Connector.Table = *postgresTableFlag
		Connector.Port = *postgresPortFlag
		Connector.Host = *postgresHostFlag
		Connector.RowCount = *postgresTotalRowsFlag
		Connector.AllTables = *postgresAllDBFlag
		Connector.IgnoreConstraints = *postgresIgnoreConstrFlag

		// If both -t and -x are provided, error out
		if Connector.AllTables && strings.TrimSpace(Connector.Table) != "" {
			log.Error("Cannot have both table (-t) and all tables (-x) flag working together, choose one.\n")
			fmt.Printf("Usage of engine: %s\n", Connector.Engine)
			postgresFlag.PrintDefaults()
			os.Exit(1)
		} else if !Connector.AllTables  && strings.TrimSpace(Connector.Table) == "" { // if -t is empty
			log.Error("Provide the list of tables (-t) to mock or -x for all database.\n")
			fmt.Printf("Usage of engine: %s\n", Connector.Engine)
			postgresFlag.PrintDefaults()
			os.Exit(1)
		}
	}

}
