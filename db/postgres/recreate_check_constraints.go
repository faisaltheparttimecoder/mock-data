package postgres

import (
	"fmt"
	"strings"
	"database/sql"
	"github.com/pivotal/mock-data/core"
)

//
// Currently there is issue with check constraints, so this check is just
// a dummy, will work on it once we have a proper idea
//


// fix Check constraints
func fixCheck(db *sql.DB, ck constraint) error {

	var TotalViolators int

	// Extract the key columns
	keys, err := core.ColExtractor(ck.column, `\(.*?\)`)
	if err != nil {
		return fmt.Errorf("Unable to determine CK violators key columns: %v", err)
	}
	cols := strings.Trim(keys, "()")

	// Extract the column name from it
	colname, err := core.ColExtractor(cols, `^(.*?)\s`)
	if err != nil {
		return fmt.Errorf("Unable to determine column name from keys: %v", err)
	}

	// Check if this table is violating any check constraints
	rows, err := db.Query(GetTotalCKViolator(ck.table, colname, cols))
	if err != nil {
		return fmt.Errorf("Unable to get the total primary key violators: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&TotalViolators)
		if err != nil {
			return fmt.Errorf("Unable to scan the total primary key violators: %v", err)
		}

	}

	log.Info(TotalViolators)

	return nil
}