# MockD [![Go Version](https://img.shields.io/badge/go-v1.7.4-green.svg?style=flat-square)](https://golang.org/dl/) [![MIT License](https://img.shields.io/badge/License-MIT_License-green.svg?style=flat-square)](https://github.com/ielizaga/mockd/blob/master/LICENSE)
                                                                                                                                                                                                       
MockD is built of simple catchphrase

    Here are my tables
    Load data for me
    I don't care how

with MockD you can have your own tables defined with datatypes, just supply the connection, total data needed and the tables and leave MockD to build random generated dataset on those tables.

MockD does try to ensure it acknowledge the constraints in the database and try to fix PK / UK / FK constraints, but it not foolproof. hence its **HIGHLY RECOMMEDED** you do take a backup of your database before you ask MockD to load data.

An ideal environment to make MockD work without any error's would be 

+ No constraints
+ No custom datatypes

Check on "Known Issues" below for things are currently not working

# Important Information

MockD program is created for generating test data and **NOT TO BE USED IN PRODUCTION DATABASE**, ensure you have the backup of the database before asking MockD to load data for you.

# Supported database

+ Postgres
+ Greenplum (GPDB)
+ Hawq (HDB) 
+ MySQL / Oracle ( coming soon <img src="https://media.giphy.com/media/MhS3BBBdYAxEc/giphy.gif" width="18" height="10"> ) 

# Supported datatypes

+ All datatypes that are listed on the [postgres datatype](https://www.postgresql.org/docs/9.6/static/datatype.html) website are supported
+ As Greenplum / Hawq is a fork from postgres, the supported postgres datatype also applies here.

# Go Library Package Dependencies

+ [Data Faker](https://github.com/icrowley/fake) by icrowley
+ [Progress bar](https://github.com/vbauerster/mpb) by vbauerster
+ [Postgres Driver](https://github.com/lib/pq) by lib/pq
+ [Go Logger](https://github.com/op/go-logging) by op/go-logging

# How it works.

+ Parse the argument parameters
+ Check if the database connection can be established and queries work.
+ If all database flag is set, then extract all the tables in the database
+ if tables are specified then uses only those tables.
+ First creates a backup of all constraints ( like PK (primary key), UK, CK, FK ) and also unique index ( due to cascade nature of the drop constraints )
+ Also store this constraint / Unique index information on the memory
+ Before loading it will remove all the constraints on the table
+ Once done, it will start its loaded by picking random data based on the columns datatype.
+ Once its done, it read all the constraints information from the memory
+ fixes PK and UK initially
+ then fixes FK
+ Check constraints are ignored ( since its not easy to translate it )
+ Once done it loads all the constraints that it had backed up ( MockD fails mostly at this stage if its not able to autofix the constraints violation )
+ Program ends.

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

### MockD User

+ [Download](https://github.com/ielizaga/mockd/archive/master.zip) the github repo or clone the github repo

```
git clone https://github.com/ielizaga/mockd.git
```

+ Navigate to the "bin" directory
+ Use the appropriate binary that matches your OS

```
XXXXX:bin XXXXX ls -ltr
total 26832
-rwxr-xr-x  1 XXXXX  XXXXX  6878692 Jul 16 10:44 mockd-mac
-rwxr-xr-x  1 XXXXX  XXXXX  6853562 Jul 16 10:44 mockd-linux
```

**NOTE:** if you have the datatype UUID defined on the table, make sure you have the execute "uuidgen" installed on the OS.  

### Developers

+ [Download](https://github.com/ielizaga/mockd/archive/master.zip) the github repo or clone the github repo

```
git clone https://github.com/ielizaga/mockd.git
```

or use "go get" to download the source after setting the GOPATH

```
go get github.com/ielizaga/mockd
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

+ For postgres / greenplum / hawq 

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

![single table](https://github.com/ielizaga/mockd/blob/master/img/singletable.gif)

+ Mock multiple table with random data

```
bin/mockd-mac <dbengine> -n <total rows> -u <user> -d <database> -t <table1>,<table2>,....
```

![multiple table](https://github.com/ielizaga/mockd/blob/master/img/multipletable.gif)

+ Mock entire database

```
bin/mockd-mac <dbengine> -n <total rows> -u <user> -d <database> -x
```

![All Database](https://github.com/ielizaga/mockd/blob/master/img/alldb.gif)

# Known Issues

1. If you uses a small amount of rows to mock and you have defined a composite PK index, there are chances that the primary key would fail ( this is due to foreign key fix only runs after primary key fix )
2. Still having issues with Check constraint, only check that works is "COLUMN > 0" & nothing else
3. On Greenplum/HDB partition tables are not supported (due to check constraint issues defined above). 
4. Custom datatypes are not supported


# Community

You can sumbit issues or pull request via [github](https://github.com/ielizaga/mockd) and we will try our level best to fix them..

# Authors

[![Ignacio](https://img.shields.io/badge/github-Ignacio_Elizaga-green.svg?style=social)](https://github.com/ielizaga) [![Aitor](https://img.shields.io/badge/github-Aitor_Cedres-green.svg?style=social)](https://github.com/Zerpet) [![Juan](https://img.shields.io/badge/github-Juan_Ramos-green.svg?style=social)](https://github.com/jujoramos) [![Faisal](https://img.shields.io/badge/github-Faisal_Ali-green.svg?style=social)](https://github.com/faisaltheparttimecoder) [![Adam](https://img.shields.io/badge/github-Adam_Clevy-green.svg?style=social)](https://github.com/adamclevy)


