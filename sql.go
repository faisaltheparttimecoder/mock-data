package main

import (
	"fmt"
	"github.com/go-pg/pg"
)

var GreenplumOrPostgres = "greenplum"

type DBTables struct {
	Schema string
	Table  string
}

type DBColumns struct {
	Column   string
	Datatype string
	Sequence string
}

type DBConstraints struct {
	Tablename      string
	Constraintname string
	Constraintkey  string
}

type DBConstraintsByTable struct {
	Tablename      string
	Constraintname string
	Constraintcol  string
	Constrainttype string
}

type DBIndex struct {
	Tablename string
	Indexdef  string
}

// Connection check and database version
func dbVersion() {
	Debug("Checking the version of the database")

	// db connection
	db := ConnectDB()
	defer db.Close()

	// Fire the database version query and check if there is no
	// database error
	var version string
	query := "SELECT version()"
	_, err := db.QueryOne(pg.Scan(&version), query)
	if err != nil {
		Fatalf("Encountered error when connecting to the database, err: %v", err)
	}
	Infof("Version of the database: %s", version)
	postgresOrGreenplum()
}

// Is this postgres or greenplum database
func postgresOrGreenplum() {
	Debug("Checking if this a greenplum or postgres DB")

	// Only greenplum has this table
	query := "select * from gp_segment_configuration"
	_, err := ExecuteDB(query)
	if err != nil {
		GreenplumOrPostgres = "postgres"
	}

	Infof("The flavour of postgres is: %s", GreenplumOrPostgres)
}

// Postgres get all the tables
func allTablesPostgres(whereClause string) []DBTables {
	Debug("Extracting the tables info from the postgres database")
	var result []DBTables

	// db connection
	db := ConnectDB()
	defer db.Close()

	// The query
	query := `
SELECT n.nspname as schema, c.relname as table
FROM   pg_catalog.pg_class c 
       LEFT JOIN pg_catalog.pg_namespace n 
              ON n.oid = c.relnamespace 
WHERE  c.relkind IN ( 'r', '' ) 
       AND n.nspname <> 'pg_catalog' 
       AND n.nspname <> 'information_schema' 
       AND n.nspname !~ '^pg_toast' 
       AND n.nspname !~ '^gp_toolkit' 
       AND c.relkind = 'r' 
       %s
ORDER  BY 1
`
	// add where clause
	query = fmt.Sprintf(query, whereClause)

	// execute the query
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the tables from postgres database, err: %v", err)
	}

	return result
}

// Greenplum / HDB get all the tables
func allTablesGPDB(whereClause string) []DBTables {
	Debug("Extracting the tables info from greenplum database")
	var result []DBTables

	// db connection
	db := ConnectDB()
	defer db.Close()

	// The query
	query := `
SELECT n.nspname as schema, c.relname as table
FROM   pg_catalog.pg_class c 
       LEFT JOIN pg_catalog.pg_namespace n 
              ON n.oid = c.relnamespace 
WHERE  c.relkind IN ( 'r', '' ) 
       AND n.nspname <> 'pg_catalog' 
       AND n.nspname <> 'information_schema' 
       AND n.nspname !~ '^pg_toast' 
       AND n.nspname <> 'gp_toolkit' 
       AND c.relkind = 'r' 
       AND c.relstorage IN ('a', 'h') 
       %s
ORDER  BY 1
`
	// add where clause
	query = fmt.Sprintf(query, whereClause)

	// execute the query
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the tables from GPDB, err: %v", err)
	}

	return result
}

// Extract Column & DataType Postgres
func columnExtractorPostgres(schema, table string) []DBColumns {
	tableName := fmt.Sprintf("%s.\"%s\"", schema, table)
	Debugf("Extracting the column information from postgres database for table: %s", tableName)
	var result []DBColumns

	// db connection
	db := ConnectDB()
	defer db.Close()

	// The query
	query := `
SELECT 
  a.attname as column, 
  pg_catalog.Format_type(a.atttypid, a.atttypmod) as datatype, 
  COALESCE(
    (
      SELECT 
        substring(
          pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128
        ) 
      FROM 
        pg_catalog.pg_attrdef d 
      WHERE 
        d.adrelid = a.attrelid 
        AND d.adnum = a.attnum 
        AND a.atthasdef
    ), 
    ''
  ) as sequence
FROM 
  pg_catalog.pg_attribute a 
WHERE 
  a.attrelid = '%s' :: regclass 
  AND a.attnum > 0 
  AND NOT a.attisdropped 
ORDER BY 
  a.attnum
`

	// add table information and execute the query
	query = fmt.Sprintf(query, tableName)
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the columns from Postgres, err: %v", err)
	}

	return result
}

// Extract Column & DataType GPDB
func columnExtractorGPDB(schema, table string) []DBColumns {
	tableName := fmt.Sprintf("%s.\"%s\"", schema, table)
	Debugf("Extracting the column information from postgres database for table: %s", tableName)
	var result []DBColumns

	// db connection
	db := ConnectDB()
	defer db.Close()

	// The query
	query := `
SELECT 
  a.attname as column, 
  pg_catalog.Format_type(a.atttypid, a.atttypmod) as datatype, 
  COALESCE(
    (
      SELECT 
        substring(
          pg_catalog.pg_get_expr(d.adbin, d.adrelid) for 128
        ) 
      FROM 
        pg_catalog.pg_attrdef d 
      WHERE 
        d.adrelid = a.attrelid 
        AND d.adnum = a.attnum 
        AND a.atthasdef
    ), 
    ''
  ) as sequence
FROM 
  pg_catalog.pg_attribute a 
  LEFT OUTER JOIN pg_catalog.pg_attribute_encoding e ON e.attrelid = a.attrelid 
  AND e.attnum = a.attnum 
WHERE 
  a.attrelid = '%s' :: regclass 
  AND a.attnum > 0 
  AND NOT a.attisdropped 
ORDER BY 
  a.attnum
`

	// add table information and execute the query
	query = fmt.Sprintf(query, tableName)
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the columns from GPDB, err: %v", err)
	}

	return result
}

// Save all the DDL of the constraint ( like PK(p), FK(f), CK(c), UK(u) )
func GetPGConstraintDDL(conntype string) []DBConstraints {
	Debugf("Extracting the DDL of the %s constraints", conntype)
	var result []DBConstraints
	query := `
SELECT n.nspname || '.' || c.relname tablename, 
	con.conname constraintname,
    pg_catalog.pg_get_constraintdef(con.oid, true) constraintKey
FROM  pg_catalog.pg_class c,
      pg_catalog.pg_constraint con,
      pg_catalog.pg_namespace n
WHERE conrelid = c.oid
AND n.oid = c.relnamespace
AND contype = '%s'
ORDER  BY tablename
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	query = fmt.Sprintf(query, conntype)
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the constriants from database, err: %v", err)
	}

	return result
}

// Get all the Unique index from the database
func GetPGIndexDDL() []DBIndex {
	Debugf("Extracting the unique indexes")
	var result []DBIndex
	query := `
SELECT schemaname ||'.'|| tablename tablename,
       indexdef indexdef
FROM   pg_indexes
WHERE  schemaname IN (
	SELECT nspname
	FROM   pg_namespace
	WHERE  nspname NOT IN (
       'pg_catalog',
       'information_schema',
       'pg_aoseg',
       'gp_toolkit',
       'pg_toast', 'pg_bitmapindex' )
	)
AND indexdef LIKE 'CREATE UNIQUE%'
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting all the constriants from database, err: %v", err)
	}

	return result
}

// Drop statement for the table
func GetConstraintsPertab(tabname string) []DBConstraintsByTable {
	Debugf("Extracting constraint info for table: %s", tabname)
	var result []DBConstraintsByTable
	query := `
SELECT * FROM (
	SELECT n.nspname || '.' || c.relname tablename,
		con.conname constraintname,
    	pg_catalog.pg_get_constraintdef(con.oid, true) constraintcol,
    	'constraint' constrainttype 
	FROM  pg_catalog.pg_class c,
          pg_catalog.pg_constraint con,
          pg_namespace n
	WHERE  c.oid = '%[1]s'::regclass
	AND conrelid = c.oid
	AND n.oid = c.relnamespace
	AND contype IN ('u','f','c','p')
	UNION
	SELECT schemaname || '.' || tablename tablename,
	    indexname conname,
		indexdef concol,
	   'index' contype
	FROM   pg_indexes 
	WHERE  schemaname IN (SELECT nspname
    					  FROM   pg_namespace
                          WHERE  nspname NOT IN (
                          	'pg_catalog',
		                  	'information_schema',
		                  	'pg_aoseg',
		                  	'gp_toolkit',
                           	'pg_toast', 
							'pg_bitmapindex')
						)
	AND indexdef LIKE 'CREATE UNIQUE%s'
	AND schemaname || '.' || tablename = '%[1]s'
	) a 
ORDER BY constrainttype
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	query = fmt.Sprintf(query, tabname, "%")
	_, err := db.Query(&result, query)
	if err != nil {
		Fatalf("Encountered error when getting constraints for table %s from database, err: %v", tabname, err)
	}

	return result
}
