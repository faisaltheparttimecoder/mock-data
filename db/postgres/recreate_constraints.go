package postgres

import (
	"database/sql"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pivotal/mock-data/core"
)

var (
	log = logging.MustGetLogger("mockd")
)

// fix the data loaded so that we can reenable the constraints
func FixConstraints(db *sql.DB, timestamp string, debug bool) error {

	// Fix the constraints in this order
	//var constr = []string{"PRIMARY", "UNIQUE", "CHECK", "FOREIGN"}
	var constr = []string{"PRIMARY", "UNIQUE", "FOREIGN"}
	for _, v := range constr {
		log.Infof("Checking for any %s KEYS, fixing them if there is any violations", v)
		for _, con := range savedConstraints[v] {
			switch {
				case v == "PRIMARY":
					err := fixPKey(db, con, v, debug)
					if err != nil {
						return err
					}
				case v == "UNIQUE": // Run the same logic as primary key
					err := fixPKey(db, con, v, debug)
					if err != nil {
						return err
					}
				case v == "CHECK":
					err := fixCheck(db, con)
					if err != nil {
						return err
					}
				case v == "FOREIGN":
					err := fixFKey(db, con, debug)
					if err != nil {
						return err
					}
			}
		}
	}

	// Recreate constraints
	failureDetected, err := recreateAllConstraints(db, timestamp)
	if failureDetected || err != nil {
		return err
	}

	return nil
}

// Recreate all the constraints of the database ( in case we have dropped any )
func recreateAllConstraints(db *sql.DB, timestamp string) (bool, error) {

	var AnyErrors bool
	log.Info("Starting to recreating all the constraints of the table ...")

	// list the backup files collected.
	for _, con := range constraints {
		backupFile, err := core.ListFile(".", "*_"+con+"_"+timestamp+".sql")
		if err != nil {
			return AnyErrors, fmt.Errorf("Error in getting the list of backup files: %v", err)
		}

		// run it only if we do find the backup file
		if len(backupFile) > 0 {
			contents, err := core.ReadFile(backupFile[0])
			if err != nil {
				return AnyErrors, fmt.Errorf("Error in reading the backup files: %v", err)
			}

			// Recreate all constraints one by one, if we can't create it then display the message
			// on the screen and continue with the rest, since we don't want it to fail if we cannot
			// recreate constraint of a single table.
			for _, content := range contents {
				_, err = db.Exec(content)
				if err != nil && !core.IgnoreErrorString(fmt.Sprintf("%s", err), ignoreErr) {
					AnyErrors = true
					log.Errorf("Failed to create constraints: \"%v\"", content)
				}
			}
		}
	}

	// If any error detected, tell the user about it
	if AnyErrors {
		return AnyErrors, fmt.Errorf("Detected failure in creating constraints... ")
	} else { // else we are all good.
		return AnyErrors, nil
	}

}
