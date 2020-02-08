package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"regexp"
	"strings"
)

// Create a database connection
func ConnectDB() *pg.DB {
	setDBDefaults()
	addr := fmt.Sprintf("%s:%d", cmdOptions.Hostname, cmdOptions.Port)
	return pg.Connect(&pg.Options{
		User:     cmdOptions.Username,
		Password: cmdOptions.Password,
		Database: cmdOptions.Database,
		Addr:     addr,
	})
}

// Execute queries in the database
func ExecuteDB(stmt string) (pg.Result, error) {
	// Connect to database
	db := ConnectDB()
	defer db.Close()

	// Execute the statement
	return db.Exec(stmt)
}

// Set database defaults if no options available
func setDBDefaults() {
	if IsStringEmpty(cmdOptions.Database) {
		cmdOptions.Database = "postgres"
	}
	if IsStringEmpty(cmdOptions.Username) {
		cmdOptions.Username = "postgres"
	}
	if IsStringEmpty(cmdOptions.Password) {
		cmdOptions.Password = "postgres"
	}
	if cmdOptions.Port == 0 {
		cmdOptions.Port = 5432
	}
	if IsStringEmpty(cmdOptions.Hostname) {
		cmdOptions.Hostname = "localhost"
	}
}

// is string empty
func IsStringEmpty(s string) bool {
	if strings.TrimSpace(s) != "" {
		return false
	}
	return true
}

// Progress Bar
func StartProgressBar(text string, max int) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetBytes(10000),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s[reset]", text)),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",

		}))
}

// Remove all special characters
// Though we allow users to have their own table and column prefix, postgres have limitation on the characters
// used, so we ensure that we only use valid characters from the string
func RemoveSpecialCharacters(s string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		Fatalf("error in compiling the string to remove special characters: %v", err)
	}
	return reg.ReplaceAllString(s, "")
}

// Inserting a array needs all the single quotes escaped
// the below function does just that
// i.e. If its array then replace " with escape to load to database
func FormatForArray(s string , isItArray, addQuotes bool) string {
	if isItArray {
		if addQuotes {
			s = fmt.Sprintf("\"%s\"", s)
		}
		return strings.Replace(s, "\"", "\\\"", -1)
	} else {
		return s
	}
}