# Mock Data [![go version](https://img.shields.io/github/go-mod/go-version/pivotal-gss/mock-data?filename=go.mod&style=flat&logo=go&label=GoLang)](https://golang.org/dl/) [![CI](https://img.shields.io/github/workflow/status/pivotal-gss/mock-data/CI?logo=github&style=flat&label=CI%20Workflow)](https://github.com/pivotal-gss/mock-data/actions/workflows/ci.yml) [![CI](https://img.shields.io/github/workflow/status/pivotal-gss/mock-data/Tests?logo=github&style=flat&label=Tests%20Workflow)](https://github.com/pivotal-gss/mock-data/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/pivotal-gss/mock-data/branch/master/graph/badge.svg)](https://codecov.io/gh/pivotal-gss/mock-data)

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

## Table of Contents

   * [Important & Disclaimer](#important--disclaimer)
   * [Supported database engines &amp; data types](#supported-database-engines--data-types)
        * [Database Engine](#database-engine)
        * [Data types](#data-types)
   * [How it works](#how-it-works)
   * [Usage](#usage)
   * [Installation](#installation)
   * [Examples](#examples)
   * [Known Issues](#known-issues)
   * [Developers / Collaboration](#developers--collaboration)
   * [Contributors](#Contributors)
   * [License](#license)

## Important & Disclaimer

Mock-data idea is to generate fake data in new test cluster, and it is **NOT TO BE USED IN PRODUCTION ENVIRONMENTS**. Please ensure you have a backup of your database before running Mock-data in an environment you can't afford losing.

## Supported database engines & data types

### Database Engine
+ PostgresSQL
+ Greenplum Database

### Data types

+ All datatypes that are listed on the [postgres datatype](https://www.postgresql.org/docs/9.6/static/datatype.html) website are supported
+ As Greenplum are both base from postgres, the supported postgres datatype also apply in their case

## How it works

+ PARSES the CLI arguments
+ CHECKS if the database connection can be established
+ BASED on sub commands i.e either database , table or schema it pull / verifies the tables
+ CREATES a backup of all constraints (PK, UK, CK, FK ) and unique indexes (due to cascade nature of the drop constraints)
+ STORES this constraint/unique index information in memory and also saves it to the file under `$HOME/mock`
+ REMOVES all the constraints on the table
+ STARTS loading random data based on the columns datatype
+ READS all the constraints information from memory
+ FIXES PK and UK initially
+ FIXES FK
+ CHECK constraints are ignored (coming soon?)
+ LOADS constraints that it had backed up (Mock-data can fail at this stage if its not able to fix the constraint violations)

## Usage

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
      --uri string        Postgres connection URI, eg. postgres://user:pass@host:=port/db?sslmode=disable
  -u, --username string   Username to connect to the database
  -v, --verbose           Enable verbose or debug logging
      --version           version for mock

Use "mock [command] --help" for more information about a command.
```

## Installation

### Using Binary
[Download](https://github.com/pivotal/mock-data/releases/latest) the latest release for your OS & Architecture and you're ready to go!

**[Optional]** You can copy the mock program to the PATH folder, so that you can use the mock from anywhere in the terminal, for eg.s

    cp mock-darwin-amd64-v2.0 /usr/local/bin/mock
    chmod +x /usr/local/bin/mock

provided `/usr/local/bin` is part of the $PATH environment variable.

### Via docker
+ Pull the image & you are all set
    ```
    docker pull ghcr.io/pivotal-gss/mock-data:latest
    ```
+ **[OPTIONAL]** add a tag for easy acess
    ```
    docker image tag ghcr.io/pivotal-gss/mock-data mock
    ```
+ For mac users to connect to the host database you can run the below command
    ```
    docker run mock -a host.docker.internal <flags...>
    ```
  eg
    ```
    docker run mock database -f -a host.docker.internal -u postgres -d demodb
    ```

## Examples

Here is a simple demo of how the tool works, provide us your table and we will load the data for you

![demo-table-loading](https://github.com/pivotal-legacy/mock-data/blob/image/images/faking-data-to-single-table.gif)

For more examples how to use the tool, please check out the [wiki](https://github.com/pivotal-legacy/mock-data/wiki) page for categories like

* Look here on how the [database connection](https://github.com/pivotal-legacy/mock-data/wiki/Connecting-to-Database) works
* Read this section on how the subcommand [custom](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Custom) works
* Read this section on how the subcommand [database](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Database) works
* Read this section on how the subcommand [schema](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Schema) works
* Read this section on how the subcommand [tables](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Tables) works
 

## Known Issues

1. We do struggle when recreating constraints, even though we do try to fix the primary key , foreign key, unique key. So there is no guarantee that the tool will fix all the constraints and manual intervention is needed in some cases.
2. If you have a composite unique index where one column is part of foreign key column then there are chances the constraint creation would fail.
3. Fixing CHECK constraints isn't supported due to complexity, so recreating check constraints would fail, use `custom` subcommand to control the data being inserted
4. On Greenplum Database partition tables are not supported (due to check constraint issues defined above), so use the `custom` sub command to define the data to be inserted to the column with check constraints
5. Custom data types are not supported, use `custom` sub command to control the data for that custom data types

## Developers / Collaboration

You can sumbit issues or pull request via [github](https://github.com/pivotal/mock-data) and we will try our best to fix them.

To customize this repository, follow the steps

1. Clone the git repository
2. Export the GOPATH
    ```
    export GOPATH=<path to the clone repository>
    ```
3. Install all the dependencies. 
    ```
    go mod vendor
    ```
4. Make sure you have a demo postgres database to connect or if you are using mac, you can use 
    ```
    make install_postgres
    make start_postgres
    make stop_postgres
    make uninstall_postgres
    ```
5. You are all set, you can run it locally using
    ```
    go run . <commands> <flags.........>
    ```
6. To run test, use 
    ```
    # Edit the database environment variables on the "Makefile"
    make unit_tests
    make integration_tests
    make tests # Runs the above two test simultaneously 
    ```
7. To build the package use
    ```
    make build
    ```
**[Optional]** For formatting and checking your code you can use linters. We have [.golangci.yml](https://github.com/pivotal-gss/mock-data/blob/master/.golangci.yml) available with the repository, check out this [blog](https://betterprogramming.pub/how-to-improve-code-quality-with-an-automatic-check-in-go-d18a5eb85f09) on how to set it up.

**--- HAPPY HACKING ---**

## Contributors

## License

The Project is licensed under [MIT](https://github.com/pivotal-legacy/mock-data/blob/master/LICENSE)