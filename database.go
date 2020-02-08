package main

// Execute demo database
func ExecuteDemoDatabase() {
	Infof("Create demo tables in the database: %s", cmdOptions.Database)

	// Execute the demo database dump
	_, err := ExecuteDB(demoDatabase())
	if err != nil {
		Fatalf("Failure in creating a demo tables in the database %s, err: %v", cmdOptions.Database, err)
	}

	Infof("Completed creating demo tables in the database: %s", cmdOptions.Database)
}

