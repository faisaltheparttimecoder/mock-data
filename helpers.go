package main

import (
	"bufio"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Extract the current time now.
func TimeNow() string {
	return time.Now().Format("20060102150405")
}

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

// Execute statement in the database
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
	// Turn off the progress bar when the Debug is one
	if cmdOptions.Debug {
		return &progressbar.ProgressBar{}
	}

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
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		Fatalf("error in compiling the string to remove special characters: %v", err)
	}
	return reg.ReplaceAllString(s, "")
}

// Inserting a array needs all the single quotes escaped
// the below function does just that
// i.e. If its array then replace " with escape to load to database
func FormatForArray(s string, isItArray bool) string {
	if isItArray {
		return fmt.Sprintf("\"%s\"", s)
	} else {
		return s
	}
}

// Prompt for confirmation
func YesOrNoConfirmation() string {
	Debugf("Promoting for yes or no confirmation")
	var YesOrNo = map[string]string{"y": "y", "ye": "y", "yes": "y", "n": "n", "no": "n"}
	question := "Are you sure the program %s can continue loading the fake data? " +
		"FYI, For faking data to the database %s the constraints on the database will be dropped. \n" +
		"NOTE: \n" +
		" 1. These constraints will be backed up & saved onto to directory\n" +
		" 2. At the end of the program there will be an attempt " +
		"to restore it, unless ignore (-i) flag is set when the restore of constraints will be ignored.\n" +
		"Choose (Yy/Nn): "

	// Start the new scanner to get the user input
	fmt.Printf(question, programName, cmdOptions.Database)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		// The choice entered
		choiceEntered := input.Text()

		// If its a valid value move on
		if YesOrNo[strings.ToLower(choiceEntered)] == "y" { // Is it Yes
			return choiceEntered
		} else if YesOrNo[strings.ToLower(choiceEntered)] == "n" { // Is it No
			Info("Canceling as per user request...")
			os.Exit(0)
		} else { // Invalid choice, ask to re-enter
			fmt.Println("Invalid Choice: Please enter Yy/Nn, try again.")
			return YesOrNoConfirmation()
		}
	}

	return ""
}

// Ignore these errors, else error out
func IgnoreError(e string, ignoreMsg string, failureMsg string) {
	if !strings.HasSuffix(e, ignoreMsg) {
		Fatalf(failureMsg)
	}
}

// If the random value of numeric datatype is greater than specifed, it ends up with
// i.e error "numeric field overflow"
// The below helper helps to reduce the size of the value
func TruncateFloat(f float64, max, precision int) float64 {
	stringFloat := strconv.FormatFloat(f, 'f', precision, 64)
	if len(stringFloat) > max {
		f = math.Log10(f)
	}
	return f
}

// Extract Float precision from the float datatypes
func FloatPrecision(dt string) (int, int, error) {

	// check if brackets exists, if it doesn't then add some virtual values
	if !BracketsExists(dt) && strings.HasSuffix(dt, "[]") {
		dt = strings.Replace(dt, "[]", "", 1) + "(5,3)[]"
	} else if !BracketsExists(dt) && !strings.HasSuffix(dt, "[]") {
		dt = dt + "(5,3)"
	}
	// Get the ranges in the brackets
	var rgx = regexp.MustCompile(`\((.*?)\)`)
	rs := rgx.FindStringSubmatch(dt)
	split := strings.Split(rs[1], ",")
	m, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, fmt.Errorf("float Precision (min): %v", err)
	}
	p, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, fmt.Errorf("float Precision (precision): %v", err)
	}
	return m, p, nil
}

// Column Extractor from the provided constraint key
func ColExtractor(conkey, regExp string) (string, error) {
	var rgx = regexp.MustCompile(regExp)
	rs := rgx.FindStringSubmatch(conkey)
	if len(rs) > 0 {
		return rs[0], nil
	} else {
		return "", fmt.Errorf("unable to extract the columns from the constraint key")
	}
	return "", nil
}

// If given a datatype see if it has a bracket or not.
func BracketsExists(dt string) bool {
	var rgx = regexp.MustCompile(`\(.*\)`)
	rs := rgx.FindStringSubmatch(dt)
	if len(rs) > 0 {
		return true
	} else {
		return false
	}
}

// Does the string contain the substring
func isSubStringAvailableOnString(s string, criteria string) bool {
	var re = regexp.MustCompile(fmt.Sprintf("%s", criteria))
	if re.MatchString(s) {
		return true
	}
	return false
}

// Built a method to find if the values exits with a slice
func StringContains(item string, slice []string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

// Build a method to find if the value starts with specific word within a slice
func StringHasPrefix(item string, slice []string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		if strings.HasPrefix(item, s) {
			set[item] = struct{}{}
		}
	}
	_, ok := set[item]
	return ok
}

// Extract total characters that the datatype char can store.
func CharLen(dt string) (int, error) {
	var rgx = regexp.MustCompile(`\((.*?)\)`)
	var returnValue int
	var err error
	rs := rgx.FindStringSubmatch(dt)
	if len(rs) > 0 { // If the datatypes has number of value defined
		returnValue, err = strconv.Atoi(rs[1])
		if err != nil {
			return 0, err
		}
	} else {
		returnValue = 1
	}
	return returnValue, nil
}