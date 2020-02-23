# Mock-data [![Go Version](https://img.shields.io/badge/go-v1.13.4-green.svg?style=flat-square)](https://golang.org/dl/)

    Here are my tables
    Load them [with data] for me
    I don't care how
    
Mock-data is the result of a Pivotal internal hackathon in July 2017. The idea behind it is to allow users to test database queries with sets of fake data in any pre-defined table.

With Mock-data users can have 

+ Their own tables defined with any particular (supported) data types. It's only needed to provide the target table(s) and the number of rows of randomly generated data to insert.
+ Create a demo database
+ Create `n` number of table with `n` number of column
+ Custom fit data into the table

An ideal environment to make Mock-data work without any errors would be 

+ Tables with no constraints
+ No custom data types

However, please **DO MAKE SURE TO TAKE A BACKUP** of your database before you mock data in it as it has not been tested extensively.

Check on the "Known Issues" section below for more information about current identified bugs.

# Important information and disclaimer

Mock-data idea is to generate fake data in new test cluster and it is **NOT TO BE USED IN PRODUCTION ENVIRONMENTS**. Please ensure you have a backup of your database before running Mock-data in an environment you can't afford losing.

# Supported database engines & data types

### Database Engine
+ PostgresSQL
+ Greenplum Database

### Data types

+ All datatypes that are listed on the [postgres datatype](https://www.postgresql.org/docs/9.6/static/datatype.html) website are supported
+ As Greenplum / HAWQ are both base from postgres, the supported postgres datatype also apply in their case

# How it works.

+ PARSES the CLI arguments
+ CHECKS if the database connection can be established
+ IF all database flag is set, then extract all the tables in the database
+ ELSE IF tables are specified then uses only target tables
+ CREATES a backup of all constraints (PK, UK, CK, FK ) and unique indexes (due to cascade nature of the drop constraints)
+ STORES this constraint/unique index information in memory and also saves it to the file under `$HOME/mock`
+ REMOVES all the constraints on the table
+ STARTS loading random data based on the columns datatype
+ READS all the constraints information from memory
+ FIXES PK and UK initially
+ FIXES FK
+ CHECK constraints are ignored (coming soon?)
+ LOADS constraints that it had backed up (Mock-data can fail at this stage if its not able to fix the constraint violations)

# Usage
```
$ mock --help
This program generates fake data into a postgres database cluster. 
PLEASE DO NOT run on a mission critical databases

Usage:
  mock [flags]
  mock [command]

Available Commands:
  custom      Controlled mocking of tables
  database    Mock at database level
  help        Help about any command
  schema      Mock at schema level
  tables      Mock at table level

Flags:
  -a, --address string    Hostname where the postgres database lives
  -d, --database string   Database to mock the data
  -q, --dont-prompt       Run without asking for confirmation
  -h, --help              help for mock
  -i, --ignore            Ignore checking and fixing constraints
  -w, --password string   Password for the user to connect to database
  -p, --port int          Port number of the postgres database
  -r, --rows int          Total rows to be faked or mocked (default 10)
  -u, --username string   Username to connect to the database
  -v, --verbose           Enable verbose or debug logging
      --version           version for mock

Use "mock [command] --help" for more information about a command.
```

# Installation

[Download](https://github.com/pivotal/mock-data/releases/latest) the latest release and you're ready to go!

**Optional:**

You can copy the mock program to the PATH folder, so that you can use the mock from anywhere in the terminal, for eg.s
```
cp mock-darwin-amd64-v2.0 /usr/local/bin/mock
chmod +x /usr/local/bin/mock
```

provided `/usr/local/bin` is part of the $PATH environment variable.

# Known Issues

1. If you have a unique index on a foreign key column then there are chance the constraint creation would fail, since mockd doesn't pick up unique value for foriegn key value it picks up random values from the reference table.
2. Still having issues with Check constraint, only check that works is "COLUMN > 0"
3. On Greenplum Datbase/HAWQ partition tables are not supported (due to check constraint issues defined above), so use the `custom` sub command to define the data to be inserted to the column with check constraints
4. Custom data types are not supported, use `custom` sub command to control the data for that custom data types

# Developers / Collaboration

You can sumbit issues or pull request via [github](https://github.com/pivotal/mock-data) and we will try our best to fix them.

To customize this repository, follow the steps

1. Clone the git repository

2. Export the GOPATH

    ```
    export GOPATH=<path to the clone repository>
    ```

3. Install all the dependencies. If you don't have dep installed, follow the instruction from [here](https://github.com/golang/dep)

    ```
    dep ensure
    ```

4. Make sure you have a demo postgres database to connect.
5. You are all set, you can run it locally using

    ```
    go run *.go <commands> .....
    ```

6. To build the package use

    ```
    /bin/sh build.sh
    ```

# License

MIT

# Authors

[![Ignacio](https://img.shields.io/badge/github-Ignacio_Elizaga-green.svg?style=social)](https://github.com/ielizaga) [![Aitor](https://img.shields.io/badge/github-Aitor_Cedres-green.svg?style=social)](https://github.com/Zerpet) [![Juan](https://img.shields.io/badge/github-Juan_Ramos-green.svg?style=social)](https://github.com/jujoramos) [![Faisal](https://img.shields.io/badge/github-Faisal_Ali-green.svg?style=social)](https://github.com/faisaltheparttimecoder) [![Adam](https://img.shields.io/badge/github-Adam_Clevy-green.svg?style=social)](https://github.com/adamclevy)
