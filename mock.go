package main

import (
	"fmt"
	"os"
)

var (
	programName        = "mock"
	programVersion     = "v2.4"
	ExecutionTimestamp = TimeNow()
	Path               = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), programName, ExecutionTimestamp)
)

// The main function block
func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
