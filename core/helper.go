package core

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gosuri/uiprogress"
)

// Built a method to find if the values exits with a slice
func StringContains(item string, slice []string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Extract total characters from the datatypes
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

// Extract Float precision extractor
func FloatPrecision(dt string) (int, int, error) {
	var rgx = regexp.MustCompile(`\((.*?)\)`)
	rs := rgx.FindStringSubmatch(dt)
	split := strings.Split(rs[1], ",")
	m, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Float Precision (min): %v", err)
	}
	p, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Float Precision (precision): %v", err)
	}
	return m, p, nil
}

// Remove values from the float to avoid numeric overflow issue
func TruncateFloat(f float64, max, precision int) float64 {
	stringFloat := strconv.FormatFloat(f, 'f', precision, 64)
	if len(stringFloat) > max {
		f = math.Log10(f)
	}
	return f
}

// Progress Bar
var bar *uiprogress.Bar

func ProgressBar(steps int, tabname string) {

	// start rendering
	uiprogress.Start()
	bar = uiprogress.AddBar(steps).AppendCompleted().PrependElapsed() // Add a new bar

	// prepend the current step to the bar
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		return "(Mocking Table: " + tabname + ")"
	})

}

// Increment Progress bar
func IncrementBar() {
	bar.Incr()
}
