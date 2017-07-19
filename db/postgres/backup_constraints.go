package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ielizaga/mockd/core"
	"strings"
)

type constraint struct {
	table, column string
}

var (
	savedConstraints = map[string][]constraint{"PRIMARY":{}, "CHECK":{}, "UNIQUE":{}, "FOREIGN":{}}
	constraints = []string{"p", "f", "u", "c"}
	ignoreErr = []string{
		"pq: multiple primary keys for table",
		"already exists"}
)

// Backup DDL of objects which are going to drop to
// allow faster and smooth transition of inputting data.
func BackupDDL(db *sql.DB, timestamp string) error {

	// Constraints
	err := backupConstraints(db, timestamp)
	if err != nil {
		return err
	}

	// Unique Index
	err = backupIndexes(db, timestamp)
	if err != nil {
		return err
	}

	return nil
}


// Backup all the constraints
func backupConstraints(db *sql.DB, timestamp string) error {

	for _, constr := range constraints {
		filename := "mockd_constriant_backup_" + constr + "_" + timestamp + ".sql"
		rows, err := db.Query(GetPGConstraintDDL(constr))
		for rows.Next() {
			var table, conname, conkey string

			// Scan and store the rows
			err = rows.Scan(&table, &conname, &conkey)
			if err != nil {
				return fmt.Errorf("Error extracting the rows of the list of constraint DDL: %v", err)
			}

			// Generate the DDL command
			message := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;\n", table, conname, conkey)

			// write this to the file
			err = core.WriteToFile(filename, message)
			if err != nil {
				return fmt.Errorf("Error in saving the constraints DDL to the file: %v", err)
			}

			// Before dropping the constraints ensure we have the information
			// saved about the state of what constraints this table had
			ctype, err := constraintFinder(constr)
			if err != nil {
				return err
			}
			savedConstraints[ctype] = append(savedConstraints[ctype], constraint{table: table, column:conkey})
		}
	}

	return nil
}

// Backup all the unique index
func backupIndexes(db *sql.DB, timestamp string) error {

	filename :="mockd_index_backup_u_" + timestamp + ".sql"
	rows, err := db.Query(GetPGIndexDDL())
	for rows.Next() {
		var table, index string
		// Scan and store the rows
		err = rows.Scan(&table, &index)
		if err != nil {
			return fmt.Errorf("Error extracting the rows of the list of Index DDL: %v", err)
		}

		// Generate the DDL command
		message := fmt.Sprintf("%s;\n", index)

		// write this to the file
		err = core.WriteToFile(filename, message)
		if err != nil {
			return fmt.Errorf("Error in saving the index DDL to the file: %v", err)
		}

		// Save all the index information
		savedConstraints["UNIQUE"] = append(savedConstraints["UNIQUE"], constraint{table: table, column:index})
	}

	return nil
}

func RemoveConstraints(db *sql.DB, table string) error {

	var statment string

	// Obtain all the constraints on the table that we going to load data to
	rows, err := db.Query(GetConstraintsPertab(table))
	if err != nil {
		return err
	}

	// scan through the rows and generate the drop command
	for rows.Next() {
		var tab, conname, concol, contype string

		// Scan and store the rows
		err = rows.Scan(&tab, &conname, &concol,  &contype)
		if err != nil {
			return fmt.Errorf("Error extracting all the constriant list on the table: %v", err)
		}

		// Generate the DROP DDL command
		if contype == "index" { // if the constriant is a index
			statment = fmt.Sprintf("DROP INDEX %s CASCADE;", conname)
		} else { // if the constraint is a constraint
			statment = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s CASCADE;", table, conname)
		}

		// Execute the statement
		_, err = db.Exec(statment)
		if err != nil {
			// Ignore does not exist error eg.s the primary key is dropped
			// then the index also goes along with it , so no need to panic here
			errmsg := fmt.Sprintf("%s", err)
			if !strings.HasSuffix(errmsg, "does not exist") {
				return err
			}
		}
	}

	return nil
}

// What is the type of constraints is it
func constraintFinder(contype string) (string, error) {
	switch {
	// Check constraint
	case strings.Contains(contype, "c"):
		return "CHECK", nil
		// Primary constraint
	case strings.Contains(contype, "p"):
		return "PRIMARY", nil
		// Foreign constraint
	case strings.Contains(contype, "f"):
		return "FOREIGN", nil
		// Unique constraint
	case strings.Contains(contype, "u"):
		return "UNIQUE", nil
	default:
		return "", fmt.Errorf("Cannot understand the type of constraints")
	}
	return "", nil
}