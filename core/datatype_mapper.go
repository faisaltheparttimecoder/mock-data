package core

import (
	"fmt"
	"strings"
)

// Data Generator
func BuildData(dt string) (interface{}, error) {

	var datekeywords = []string{"date", "time without time zone", "time with time zone", "timestamp without time zone", "timestamp with time zone"}
	var ipkeywords = []string{"inet", "cidr"}
	var intkeywords = []string{"smallint", "integer", "bigint"}
	var intranges = map[string]int{"smallint": 767, "integer": 7483647, "bigint": 372036854775807}
	var floatkeywords = []string{"double precision", "numeric", "real"}

	switch {

	// Generate Random Integer
	case StringContains(dt, intkeywords):
		value, err := RandomInt(-intranges[dt], intranges[dt])
		if err != nil {
			return "", fmt.Errorf("Build Integer: %v", err)
		}
		return value, nil

		// Generate Random characters
	case strings.HasPrefix(dt, "character"):
		l, err := CharLen(dt)
		if err != nil {
			return "", fmt.Errorf("Build character: %v", err)
		}
		value := RandomString(l)
		return value, nil

		// Generate Random date, timestamp etc
	case StringContains(dt, datekeywords):
		value, err := RandomDate(-10, 10)
		if err != nil {
			return "", fmt.Errorf("Build Date: %v", err)
		}
		return value, nil

		// Generate Random ips
	case StringContains(dt, ipkeywords):
		return RandomIP(), nil

		// Generate Random boolean
	case strings.EqualFold(dt, "boolean"):
		return RandomBoolean(), nil

		// Generate Random text
	case strings.EqualFold(dt, "text"):
		return RandomParagraphs(), nil

		// Generate Random float values
	case StringContains(dt, floatkeywords):
		value, err := RandomFloat(1, intranges["smallint"], 3)
		if err != nil {
			return "", fmt.Errorf("Build Float: %v", err)
		}
		return value, nil

		// Generate Random numeric values
	case strings.HasPrefix(dt, "numeric"):
		max, precision, err := FloatPrecision(dt)
		if err != nil {
			return "", fmt.Errorf("Build Numeric: %v", err)
		}
		value, err := RandomFloat(0, max, precision)
		value = TruncateFloat(value, max, precision)
		if err != nil {
			return "", fmt.Errorf("Build Numeric: %v", err)
		}
		return value, nil

		// Random bit generator
	case strings.HasPrefix(dt, "bit"):
		l, err := CharLen(dt)
		if err != nil {
			return "", fmt.Errorf("Build bit: %v", err)
		}
		value := RandomBit(l)
		return value, nil

		// Random UUID generator
	case strings.HasPrefix(dt, "uuid"):
		uuid, err := RandomUUID()
		if err != nil {
			return "", fmt.Errorf("Build UUID: %v", err)
		}
		return uuid, nil

		// Random MacAddr Generator
	case strings.HasPrefix(dt, "macaddr"):
		return RandomMacAddress(), nil

	default:
		return "", fmt.Errorf("Unsupported datatypes found: %v", dt)
	}

	return "", nil
}
