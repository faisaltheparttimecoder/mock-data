package main

var (
	programName    = "mock"
	programVersion = "v2.0"
)

// The main function block
func main() {
	// Execute the cobra CLI & run the program
	rootCmd.Execute()
}
