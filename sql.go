package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"strings"
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

type DBConstraintsByDataType struct {
	Colname string
	Dtype   string
}

type DBIndex struct {
	Tablename string
	Indexdef  string
}

type DBViolationRow struct {
	Row string
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
		Debugf("query: %s", query)
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
SELECT n.nspname AS SCHEMA, 
       c.relname AS table 
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
		Debugf("query: %s", query)
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
SELECT    n.nspname AS SCHEMA, 
          c.relname AS TABLE 
FROM      pg_catalog.pg_class c 
LEFT JOIN pg_catalog.pg_namespace n 
ON        n.oid = c.relnamespace 
WHERE     c.relkind IN ( 'r', '' ) 
AND       n.nspname <> 'pg_catalog' 
AND       n.nspname <> 'information_schema' 
AND       n.nspname !~ '^pg_toast' 
AND       n.nspname <> 'gp_toolkit' 
AND       c.relkind = 'r' 
AND       c.relstorage IN ('a','h') 
%s 
ORDER BY  1
`
	// add where clause
	query = fmt.Sprintf(query, whereClause)

	// execute the query
	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
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
SELECT   a.attname                                       AS COLUMN, 
         pg_catalog.Format_type(a.atttypid, a.atttypmod) AS datatype, 
         COALESCE( 
                    ( 
                    SELECT substring( pg_catalog.Pg_get_expr(d.adbin, d.adrelid) for 128 ) 
                    FROM   pg_catalog.pg_attrdef d 
                    WHERE  d.adrelid = a.attrelid 
                    AND    d.adnum = a.attnum 
                    AND    a.atthasdef ), '' ) AS sequence 
FROM     pg_catalog.pg_attribute a 
WHERE    a.attrelid = '%s' :: regclass 
AND      a.attnum > 0 
AND      NOT a.attisdropped 
ORDER BY a.attnum
`

	// add table information and execute the query
	query = fmt.Sprintf(query, tableName)
	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
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
SELECT          a.attname                                       AS COLUMN, 
                pg_catalog.Format_type(a.atttypid, a.atttypmod) AS datatype, 
                COALESCE( 
                           ( 
                           SELECT substring( pg_catalog.Pg_get_expr(d.adbin, d.adrelid) for 128 ) 
                           FROM   pg_catalog.pg_attrdef d 
                           WHERE  d.adrelid = a.attrelid 
                           AND    d.adnum = a.attnum 
                           AND    a.atthasdef ), '' ) AS sequence 
FROM            pg_catalog.pg_attribute a 
LEFT OUTER JOIN pg_catalog.pg_attribute_encoding e 
ON              e.attrelid = a.attrelid 
AND             e.attnum = a.attnum 
WHERE           a.attrelid = '%s' :: regclass 
AND             a.attnum > 0 
AND             NOT a.attisdropped 
ORDER BY        a.attnum
`

	// add table information and execute the query
	query = fmt.Sprintf(query, tableName)
	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
		Fatalf("Encountered error when getting all the columns from GPDB, err: %v", err)
	}

	return result
}

// Save all the DDL of the constraint ( like PK(p), FK(f), CK(c), UK(u) )
func GetPGConstraintDDL(conntype string) []DBConstraints {
	Debugf("Extracting the DDL of the %s constraints", conntype)
	var result []DBConstraints
	query := `
SELECT '"' 
       || n.nspname 
       || '"."' 
       || c.relname 
       || '"'                                         tablename, 
       con.conname                                    constraintname, 
       pg_catalog.Pg_get_constraintdef(con.oid, true) constraintKey 
FROM   pg_catalog.pg_class c, 
       pg_catalog.pg_constraint con, 
       pg_catalog.pg_namespace n 
WHERE  conrelid = c.oid 
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
		Debugf("query: %s", query)
		Fatalf("Encountered error when getting all the constriants from database, err: %v", err)
	}

	return result
}

// Get all the Unique index from the database
func GetPGIndexDDL() []DBIndex {
	Debugf("Extracting the unique indexes")
	var result []DBIndex
	query := `
SELECT '"' 
       || schemaname 
       || '"."' 
       || tablename 
       || '"'   tablename, 
       indexdef indexdef 
FROM   pg_indexes 
WHERE  schemaname IN (SELECT nspname 
                      FROM   pg_namespace 
                      WHERE  nspname NOT IN ( 'pg_catalog', 'information_schema' 
                                              , 
                                              'pg_aoseg', 
                                              'gp_toolkit', 
                                              'pg_toast', 'pg_bitmapindex' )) 
       AND indexdef LIKE 'CREATE UNIQUE%' 
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
		Fatalf("Encountered error when getting all the constriants from database, err: %v", err)
	}

	return result
}

// Drop statement for the table
func GetConstraintsPertab(tabname string) []DBConstraintsByTable {
	Debugf("Extracting constraint info for table: %s", tabname)
	var result []DBConstraintsByTable
	query := `
SELECT * 
FROM   (SELECT n.nspname 
               || '.' 
               || c.relname                                   tablename, 
               con.conname                                    constraintname, 
               pg_catalog.Pg_get_constraintdef(con.oid, TRUE) constraintcol, 
               'constraint'                                   constrainttype 
        FROM   pg_catalog.pg_class c, 
               pg_catalog.pg_constraint con, 
               pg_namespace n 
        WHERE  c.oid = '%[1]s' :: regclass 
               AND conrelid = c.oid 
               AND n.oid = c.relnamespace 
               AND contype IN ( 'u', 'f', 'c', 'p' ) 
        UNION 
        SELECT schemaname 
               || '.' 
               || tablename tablename, 
               indexname    conname, 
               indexdef     concol, 
               'index'      contype 
        FROM   pg_indexes 
        WHERE  schemaname IN (SELECT nspname 
                              FROM   pg_namespace 
                              WHERE  nspname NOT IN ( 
                                     'pg_catalog', 'information_schema' 
                                     , 
                                     'pg_aoseg', 
                                     'gp_toolkit', 
                                     'pg_toast', 'pg_bitmapindex' )) 
               AND indexdef LIKE 'CREATE UNIQUE%[2]s' 
               AND '"' || schemaname || '"'
                   || '."' 
                   || tablename || '"' = '%[1]s') a 
ORDER  BY constrainttype 
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	query = fmt.Sprintf(query, tabname, "%")

	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
		Fatalf("Encountered error when getting constraints for table %s from database, err: %v", tabname, err)
	}

	return result
}

// Get the datatype of the column
func getDatatype(tab string, columns []string) []DBConstraintsByDataType {
	Debugf("Extracting constraint column data type info for table: %s", tab)
	var result []DBConstraintsByDataType
	whereClause := strings.Join(columns, "' or attname = '")
	whereClause = strings.Replace(whereClause, "attname = ' ", "attname = '", -1)
	query := `
SELECT attname                                     colname, 
       pg_catalog.Format_type(atttypid, atttypmod) dtype 
FROM   pg_attribute 
WHERE  attname = '%s' 
       AND attrelid = '%s' :: regclass 
`
	// db connection
	db := ConnectDB()
	defer db.Close()

	// add table information and execute the query
	query = fmt.Sprintf(query, whereClause, tab)
	_, err := db.Query(&result, query)
	if err != nil {
		Debugf("query: %s", query)
		Fatalf("Encountered error when getting data type for "+
			"building constrints for table %s from database, err: %v", tab, err)
	}

	return result
}

// Primary key violation check
func getTotalPKViolator(tab, cols string) int {
	var total int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM ( %s ) a`, getPKViolator(tab, cols))

	// db connection
	db := ConnectDB()
	defer db.Close()

	_, err := db.Query(pg.Scan(&total), query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when execute the query to extract pk violators: %v", err)
	}

	return total
}

// Total Primary Key violator
func getPKViolator(tab, cols string) string {
	return fmt.Sprintf(`SELECT %[1]s FROM %[2]s GROUP BY %[1]s HAVING COUNT(*) > 1`, cols, tab)
}

// Get the list of the PK violators
func GetPKViolators(tab, cols string) []DBViolationRow {
	Debugf("Extracting the unique violations for table %s and column %s", tab, cols)
	var result []DBViolationRow

	// db connection
	db := ConnectDB()
	defer db.Close()

	query := strings.Replace(getPKViolator(tab, cols), "SELECT "+cols, "SELECT "+cols+" AS row", -1)
	_, err := db.Query(&result, query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when execute the query to extract pk violators for table %s: %v", tab, err)
	}

	return result
}

// Fix PK Violators
func UpdatePKKey(tab, col, whichrow, newdata string) string {
	query := `
UPDATE %[1]s 
SET    %[2]s = '%[3]s' 
WHERE  ctid = 
       ( 
              SELECT ctid 
              FROM   %[1]s 
              WHERE  %[2]s = '%[4]s' limit 1 )
`
	query = fmt.Sprintf(query, tab, col, newdata, whichrow)
	_, err := ExecuteDB(query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when updating the primary key for table %s, err: %v", tab, err)
	}
	return ""
}

// Get the foreign violations keys
func getFKViolators(key ForeignKey) string {
	query := `
SELECT %[1]s 
FROM   %[2]s 
WHERE %[1]s  NOT IN 
       ( 
              SELECT %[3]s 
              FROM   %[4]s )
`
	return fmt.Sprintf(query, key.Column, key.Table, key.Refcolumn, key.Reftable)
}

// Get total FK violators
func GetTotalFKViolators(key ForeignKey) int {
	var total int
	query := `SELECT COUNT(*) FROM (%s) a`
	query = fmt.Sprintf(query, getFKViolators(key))

	// db connection
	db := ConnectDB()
	defer db.Close()

	_, err := db.Query(pg.Scan(&total), query)
	if err != nil {
		fmt.Println()
		Debugf("Query: %s", query)
		Fatalf("Error when execute the query to total rows of foreign keys for table %s: %v", key.Table, err)
	}

	return total
}

// Total rows of the table
func TotalRows(tab string) int {
	var total int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tab)

	// db connection
	db := ConnectDB()
	defer db.Close()

	_, err := db.Query(pg.Scan(&total), query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when execute the query to total rows: %v", err)
	}

	return total
}

// Get the list of the FK violators
func GetFKViolators(key ForeignKey) []DBViolationRow {
	Debugf("Extracting the foreign violations for table %s and column %s", key.Table, key.Reftable)
	var result []DBViolationRow

	// db connection
	db := ConnectDB()
	defer db.Close()

	query := strings.Replace(getFKViolators(key), "SELECT "+key.Column, "SELECT "+key.Column+" AS row", -1)
	_, err := db.Query(&result, query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when execute the query to extract fk violators for table %s: %v", key.Table, err)
	}

	return result
}

// Update FK violators with rows from the referenced table
func UpdateFKeys(key ForeignKey, totalRows int, whichRow string) {
	query := `
UPDATE %[1]s 
SET    %[2]s = 
       ( 
              SELECT %[3]s 
              FROM   %[4]s offset floor(random()*%[5]d) limit 1) 
WHERE  %[2]s = '%[6]s'
`
	query = fmt.Sprintf(query, key.Table, key.Column, key.Refcolumn, key.Reftable, totalRows, whichRow)
	_, err := ExecuteDB(query)
	if err != nil {
		fmt.Println()
		Debugf("query: %s", query)
		Fatalf("Error when updating the foreign key for table %s, err: %v", key.Table, err)
	}
}

// Delete the violating key
func deleteViolatingConstraintKeys(tab string, column string) error {
	Debugf("Deleting the rows of the table that violate the constraints: %s:(%s)", tab, column )
	query := `
DELETE 
FROM   %[1]s 
WHERE  ( 
              %[2]s) IN 
       ( 
                SELECT   %[2]s 
                FROM     %[1]s 
                GROUP BY %[2]s 
                HAVING   count(*) > 1);
`
	query = fmt.Sprintf(query, tab, column)
	_, err := ExecuteDB(query)
	if err != nil {
		return err
	}
	return nil
}
