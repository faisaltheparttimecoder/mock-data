# Mock-data [![Go Version](https://img.shields.io/badge/go-v1.7.4-green.svg?style=flat-square)](https://golang.org/dl/)

    Here are my tables
    Load them [with data] for me
    I don't care how
    
Mock-data is the result of a Pivotal internal hackathon in July 2017. The original idea behind it is to allow users to test database queries with sets of fake data in any pre-defined table.

With Mock-data users can have their own tables defined with any particular datatypes. It's only needed to provide the target table(s) and the number of rows of randomly generated data to insert.

An ideal environment to make Mock-data work without any errors would be 
+ Tables with no constraints
+ No custom datatypes

On a second iteration work has been done to ensure proper functioning with database constraints such as primary keys, unique keys or foreign keys. However, please **DO MAKE SURE TO TAKE A BAKCUP** of your database before you mock data in it as it has not been tested extensively.

Check on the "Known Issues" section below for more information about current identified bugs.

# Important information and disclaimer

Mock-data idea is to generate fake data in new test cluster and it is **NOT TO BE USED IN PRODUCTION ENVIRONMENTS**. Please ensure you have a backup of your database before running Mock-data in an environment you can't afford losing.

# Supported database engines

+ PostgresSQL
+ Greenplum Database
+ HAWQ/HDB
+ MySQL (coming soon?) 
+ Oracle (coming soon?) 

# Supported datatypes

+ All datatypes that are listed on the [postgres datatype](https://www.postgresql.org/docs/9.6/static/datatype.html) website are supported
+ As Greenplum / HAWQ are both base from postgres, the supported postgres datatype also apply in their case

# Dependencies

+ [Data Faker](https://github.com/icrowley/fake) by icrowley
+ [Progress bar](https://github.com/vbauerster/mpb) by vbauerster
+ [Postgres Driver](https://github.com/lib/pq) by lib/pq
+ [Go Logger](https://github.com/op/go-logging) by op/go-logging

# How it works.

+ PARSES the CLI arguments
+ CHECKS if the database connection can be established
+ IF all database flag is set, then extract all the tables in the database
+ ELSE IF tables are specified then uses only target tables
+ CREATES a backup of all constraints (PK, UK, CK, FK ) and unique indexes (due to cascade nature of the drop constraints)
+ STORES this constraint/unique index information in memory
+ REMOVES all the constraints on the table
+ STARTS loading random data based on the columns datatype
+ READS all the constraints information from memory
+ FIXES PK and UK initially
+ FIXES FK
+ CHECK constraints are ignored (coming soon?)
+ LOADS constraints that it had backed up (Mock-data can fail at this stage if its not able to fix the constraint violations)

# Usage

```
USAGE: mockd <DATABASE ENGINE> <OPTIONS>
DATABASE ENGINE:
	postgres        Postgres database
	greenplum       Greenplum database
	hdb             Hawq Database
	help            Show help
OPTIONS:
	Execute "mockd <database engine> -h" for all the database specific options
```

# How to use it

### Users

[Download](https://github.com/pivotal/mock-data/releases/tag/v1.0) the latest release and you're ready to go!

**NOTE:** if you have the datatype UUID defined on the table, make sure you have the execute "uuidgen" installed on the OS.  

### Developers

+ Clone the github repo

```
git clone https://github.com/pivotal/mock-data.git
```

or use "go get" to download the source after setting the GOPATH

```
go get github.com/pivotal/mock-data
```

+ Download all the dependencies

```
go get github.com/icrowley/fake
go get github.com/vbauerster/mpb
go get github.com/lib/pq
go get github.com/op/go-logging
```

+ You can modify the code and execute using command before creating a build

```
go run *.go <database engine>
```

+ To build binaries for different OS, you can use for eg.s

```
env GOOS=linux GOARCH=amd64 go build
```

# Command Reference

+ For PostgresSQL / Greenplum Database / HAWQ

```
XXXXX:bin XXXXX ./mockd-mac postgres -help
2017-07-16 10:58:43.609:INFO > Parsing all the command line arguments
Usage of postgres:
  -d string
    	The database name where the table resides (default "postgres")
  -h string
    	The hostname that can be used to connect to the database (default "localhost")
  -i	Ignore checking and fixing constraint issues
  -n int
    	The total number of mocked rows that is needed (default 1)
  -p int
    	The port that is used by the database engine (default 5432)
  -t string
    	The table name to be filled in with mock data
  -u string
    	The username that can be used to connect to the database (default "postgres")
  -w string
    	The password for the user that can be used to connect to the database
  -x	Mock all the tables in the database
```

# Examples

+ Mock one table will random data

```
bin/mockd-mac <dbengine> -n <total rows> -u <user> -d <database> -t <table>
```

![single table](https://github.com/pivotal/mock-data/blob/master/img/singletable.gif)

+ Mock multiple table with random data

```
bin/mockd-mac <dbengine> -n <total rows> -u <user> -d <database> -t <table1>,<table2>,....
```

![multiple table](https://github.com/pivotal/mock-data/blob/master/img/multipletable.gif)

+ Mock entire database

```
bin/mockd-mac <dbengine> -n <total rows> -u <user> -d <database> -x
```

![All Database](https://github.com/pivotal/mock-data/blob/master/img/alldb.gif)

# Known Issues

1. If you have a unique index on a foreign key column then there are chance the constraint creation would fail, since mockd doesn't pick up unique value for foriegn key value it picks up random values from the reference table.
2. Still having issues with Check constraint, only check that works is "COLUMN > 0"
3. On Greenplum Datbase/HAWQ partition tables are not supported (due to check constraint issues defined above)
4. Custom datatypes are not supported

# Collaborate

You can sumbit issues or pull request via [github](https://github.com/pivotal/mock-data) and we will try our best to fix them.

# Authors

[![Ignacio](https://img.shields.io/badge/github-Ignacio_Elizaga-green.svg?style=social)](https://github.com/ielizaga) [![Aitor](https://img.shields.io/badge/github-Aitor_Cedres-green.svg?style=social)](https://github.com/Zerpet) [![Juan](https://img.shields.io/badge/github-Juan_Ramos-green.svg?style=social)](https://github.com/jujoramos) [![Faisal](https://img.shields.io/badge/github-Faisal_Ali-green.svg?style=social)](https://github.com/faisaltheparttimecoder) [![Adam](https://img.shields.io/badge/github-Adam_Clevy-green.svg?style=social)](https://github.com/adamclevy)
