# Mock Data [![go version](https://img.shields.io/github/go-mod/go-version/pivotal-gss/mock-data?filename=go.mod&style=flat&logo=go&label=Go)](https://golang.org/dl/) [![CI](https://img.shields.io/github/workflow/status/pivotal-gss/mock-data/CI?logo=github&style=flat&label=CI)](https://github.com/pivotal-gss/mock-data/actions/workflows/ci.yml) [![CI](https://img.shields.io/github/workflow/status/pivotal-gss/mock-data/Tests?logo=github&style=flat&label=Tests)](https://github.com/pivotal-gss/mock-data/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/pivotal-gss/mock-data/branch/master/graph/badge.svg)](https://codecov.io/gh/pivotal-gss/mock-data) [![Go Report Card](https://goreportcard.com/badge/github.com/pivotal-gss/mock-data?logo=go)](https://goreportcard.com/report/github.com/pivotal-gss/mock-data) [![Github Releases Stats of mock-data](https://img.shields.io/github/downloads/pivotal-gss/mock-data/total.svg?logo=github&label=Downloads)](https://somsubhra.github.io/github-release-stats/?username=pivotal-gss&repository=mock-data)
    Here are my tables
    Load them [with data] for me
    I don't care how
    
Mock-data is the result of a Pivotal internal hackathon in July 2017. The idea behind it is to allow users to test database queries with sets of fake data in any pre-defined table.

With Mock-data users can have

+ Their own tables defined with any particular (supported) data types. It's only needed to provide the target table(s), and the number of rows of randomly generated data to insert.
+ Create a demo database
+ Create `n` number of table with `n` number of column
+ Custom fit data into the table
+ Option to select `realistic` data to be loaded onto the table

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

### Using Docker
+ Pull the image & you are all set
    ```
    docker pull ghcr.io/faisaltheparttimecoder/mock-data:latest
    ```
+ **[OPTIONAL]** add a tag for easy acess
    ```
    docker image tag ghcr.io/faisaltheparttimecoder/mock-data mock
    ```
+ Create a local directory on the host to mount has a volume inside the container, needed to store files (eg.s constraints list) or to send in configuration files to the mock data tool (like custom subcommand)
    ```
    mkdir /tmp/mock
    ```
+ Now run the docker command
    ```
    docker run -v /tmp/mock:/home/mock [docker-image-tag] [subcommand] <flags...>
    ```  
    eg.s
    ```
    docker run -v /tmp/mock:/home/mock mock database -f -u postgres -d demodb
    ```
+ For mac users to connect to the host(or local host) database you can use the address `host.docker.internal` as shown in the below command
    ```
    docker run -v /tmp/mock:/home/mock [docker-image-tag] [subcommand] -a host.docker.internal <flags...>
    ```
    eg.s
    ```
    docker run -v /tmp/mock:/home/mock mock database -f -a host.docker.internal -u postgres -d demodb
    ```
 + **[Optional]** You can also make an alias of the above command, for eg.s alias with `.zshrc` 
    ```
    echo alias mock=\"docker run -it -v /tmp/mock:/home/mock ghcr.io/faisaltheparttimecoder/mock-data:latest\" >> ~/.zshrc
    source ~/.zshrc
    mock tables -t "public.gardens" --uri="postgres://pg_user:mypassword@myhost:5432/database_name?sslmode=disable"
    ```

## Examples

Here is a simple demo of how the tool works, provide us your table and we will load the data for you

![demo-table-loading](https://github.com/pivotal-legacy/mock-data/blob/image/images/faking-data-to-single-table.gif)

For more examples how to use the tool, please check out the [wiki](https://github.com/pivotal-legacy/mock-data/wiki) page for categories like

* Look here on how the [database connection](https://github.com/pivotal-legacy/mock-data/wiki/Connecting-to-Database) works
* For realistic & controlled data, read this section on how the subcommand [custom](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Custom) works
* For mocking the whole database or creating a demo database, read this section on how the subcommand [database](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Database) works
* For mocking the whole tables of the schema, read this section on how the subcommand [schema](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Schema) works
* For creating fake tables and mocking selected tables, read this section on how the subcommand [tables](https://github.com/pivotal-legacy/mock-data/wiki/Sub-command:-Tables) works
 

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
6. [Recommended] Run the golang linter to analyzes & fix source code programming errors, bugs, stylistic errors, and suspicious constructs.
    ```
    golangci-lint run
    ```
   to install golangci-lint check [here](https://golangci-lint.run/usage/install/), config file `.golangci.yml` has been provided with this repo
7. To run test, use 
    ```
    # Edit the database environment variables on the "Makefile"
    make unit_tests
    make integration_tests
    make tests # Runs the above two test simultaneously 
    ```
8. To build the package use
    ```
    make build
    ```

**--- HAPPY HACKING ---**

## Contributors

<table>
<tr>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/faisaltheparttimecoder>
            <img src=https://avatars.githubusercontent.com/u/16798908?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Faisal Ali/>
            <br />
            <sub style="font-size:14px"><b>Faisal Ali</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/janpio>
            <img src=https://avatars.githubusercontent.com/u/183673?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Jan Piotrowski/>
            <br />
            <sub style="font-size:14px"><b>Jan Piotrowski</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/matt-song>
            <img src=https://avatars.githubusercontent.com/u/5155257?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Matt Song/>
            <br />
            <sub style="font-size:14px"><b>Matt Song</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/Zerpet>
            <img src=https://avatars.githubusercontent.com/u/1515757?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Aitor Pérez Cedres/>
            <br />
            <sub style="font-size:14px"><b>Aitor Pérez Cedres</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/andreasgan>
            <img src=https://avatars.githubusercontent.com/u/727125?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Andreas Gangsø/>
            <br />
            <sub style="font-size:14px"><b>Andreas Gangsø</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/art-frela>
            <img src=https://avatars.githubusercontent.com/u/32013768?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Artem/>
            <br />
            <sub style="font-size:14px"><b>Artem</b></sub>
        </a>
    </td>
</tr>
<tr>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/jujoramos>
            <img src=https://avatars.githubusercontent.com/u/8254001?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Juan José Ramos/>
            <br />
            <sub style="font-size:14px"><b>Juan José Ramos</b></sub>
        </a>
    </td>
    <td align="center" style="word-wrap: break-word; width: 150.0; height: 150.0">
        <a href=https://github.com/miguelff>
            <img src=https://avatars.githubusercontent.com/u/210307?v=4 width="100;"  style="border-radius:50%;align-items:center;justify-content:center;overflow:hidden;padding-top:10px" alt=Miguel Fernández/>
            <br />
            <sub style="font-size:14px"><b>Miguel Fernández</b></sub>
        </a>
    </td>
</tr>
</table>

## License

The Project is licensed under [MIT](https://github.com/pivotal-legacy/mock-data/blob/master/LICENSE)
