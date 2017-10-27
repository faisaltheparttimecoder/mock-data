package core

import (
	"fmt"
	"strings"
)

// Data Generator
// It provided random data based on datatypes.
func BuildData(dt string) (interface{}, error) {
	// ranges of dates
	var fromyear = -10
	var toyear = 10

	// Time datatypes
	var Intervalkeywords = []string{"interval", "time without time zone"}

	// Networking datatypes
	var ipkeywords = []string{"inet", "cidr"}

	// Integer datatypes
	var intkeywords = []string{"smallint", "integer", "bigint"}
	var intranges = map[string]int{"smallint": 2767, "integer": 7483647, "bigint": 372036854775807}

	// Decimal datatypes
	var floatkeywords = []string{"double precision", "real", "money"}

	// Geometry datatypes
	var geoDataTypekeywords = []string{"path", "polygon", "line", "lseg", "box", "circle", "point"}

	switch {

		// Generate Random Integer
		case StringHasPrefix(dt, intkeywords):
			if strings.HasSuffix(dt, "[]") { // Its requesting for a array of data
				nonArraydt := strings.Replace(dt, "[]", "", 1)
				value, err := RandomIntArray(-intranges[nonArraydt], intranges[nonArraydt])
				if err != nil {
					return "", fmt.Errorf("Build Integer Array: %v", err)
				}
				return value, nil
			} else { // Not a array, but a single entry request
				value, err := RandomInt(-intranges[dt], intranges[dt])
				if err != nil {
					return "", fmt.Errorf("Build Integer: %v", err)
				}
				return value, nil
			}

		// Generate Random characters
		case strings.HasPrefix(dt, "character"):
			l, err := CharLen(dt)
			if err != nil {
				return "", fmt.Errorf("Getting Character Length: %v", err)
			}
			if strings.HasSuffix(dt, "[]") {
				value := RandomStringArray(l)
				return value, nil
			} else {
				value := RandomString(l)
				return value, nil
			}

		// Generate Random date
		case strings.HasPrefix(dt, "date"):
			if strings.HasSuffix(dt, "[]") {
				value, err := RandomDateArray(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Date Array: %v", err)
				}
				return value, nil
			} else {
				value, err := RandomDate(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Date: %v", err)
				}
				return value, nil
			}

		// Generate Random timestamp without timezone
		case strings.HasPrefix(dt, "timestamp without time zone"):
			if strings.HasSuffix(dt, "[]") {
				value, err := RandomTimestampArray(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Timestamp  without timezone Array: %v", err)
				}
				return value, nil
			} else {
				value, err := RandomTimestamp(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Timestamp  without timezone: %v", err)
				}
				return value, nil
			}

		// Generate Random timestamp with timezone
		case strings.HasPrefix(dt, "timestamp with time zone"):
			if strings.HasSuffix(dt, "[]") {
				value, err := RandomTimestamptzArray(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Timestamp with timezone Array: %v", err)
				}
				return value, nil
			} else {
				value, err := RandomTimestamptz(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Timestamp with timezone: %v", err)
				}
				return value, nil
			}

		// Generate Random time without timezone
		case StringHasPrefix(dt, Intervalkeywords):
			if strings.HasSuffix(dt, "[]") {
				value, err := RandomTimeArray(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Array Time without timezone: %v", err)
				}
				return value, nil
			} else {
				value, err := RandomTime(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Time without timezone: %v", err)
				}
				return value, nil
			}

		// Generate Random time with timezone
		case strings.HasPrefix(dt, "time with time zone"):
			if strings.HasSuffix(dt, "[]") {
				value, err := RandomTimetzArray(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Time with timezone array: %v", err)
				}
				return value, nil
			} else {
				value, err := RandomTimetz(fromyear, toyear)
				if err != nil {
					return "", fmt.Errorf("Build Time with timezone: %v", err)
				}
				return value, nil
			}

		// Generate Random ips
		case StringContains(dt, ipkeywords):
			return RandomIP(), nil

		// Generate Random boolean
		case strings.EqualFold(dt, "boolean"):
			return RandomBoolean(), nil

		// Generate Random text
		case strings.HasPrefix(dt, "text"):
			if strings.HasSuffix(dt, "[]") {
				return RandomParagraphsArray(), nil
			} else {
				return RandomParagraphs(), nil
			}

		// Generate Random text & bytea
		case strings.EqualFold(dt, "bytea"):
			return RandomBytea(1024 * 1024), nil

		// Generate Random float values
		case StringHasPrefix(dt, floatkeywords):
			if strings.HasSuffix(dt, "[]") { // Float array
				value, err := RandomfloatArray(1, intranges["smallint"], 3)
				if err != nil {
					return "", fmt.Errorf("Build Float Array: %v", err)
				}
				return value, nil
			} else { // non float array
				value, err := RandomFloat(1, intranges["smallint"], 3)
				if err != nil {
					return "", fmt.Errorf("Build Float: %v", err)
				}
				return value, nil
			}

		// Generate Random numeric values with precision
		case strings.HasPrefix(dt, "numeric"):
			max, precision, err := FloatPrecision(dt)
			if err != nil {
				return "", fmt.Errorf("Build Numeric: %v", err)
			}
			if strings.HasSuffix(dt, "[]") { // Numeric Array
				value, err := RandomfloatArray(0, max, precision)
				if err != nil {
					return "", fmt.Errorf("Build Numeric Float Array: %v", err)
				}
				return value, nil
			} else { // Non numeric array
				value, err := RandomFloat(0, max, precision)
				if err != nil {
					return "", fmt.Errorf("Build Numeric Float Array: %v", err)
				}
				value = TruncateFloat(value, max, precision)
				return value, nil
			}

		// Random bit generator
		case strings.HasPrefix(dt, "bit"):
			l, err := CharLen(dt)
			if err != nil {
				return "", fmt.Errorf("Build bit: %v", err)
			}
			if strings.HasSuffix(dt, "[]") {
				value := RandomBitArray(l)
				return value, nil
			} else {
				value := RandomBit(l)
				return value, nil
			}

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

		// Random Json
		case strings.HasPrefix(dt, "json"):
			return RandomJson(), nil

		// Random XML
		case strings.EqualFold(dt, "xml"):
			return RandomXML(), nil

		// Random Text Search Query
		case strings.EqualFold(dt, "tsquery"):
			return RandomTSQuery(), nil

		// Random Text Search Vector
		case strings.EqualFold(dt, "tsvector"):
			return RandomTSVector(), nil

		// Random Log Sequence number
		case strings.EqualFold(dt, "pg_lsn"):
			return RandomLSN(), nil

		// Random Log Sequence number
		case strings.EqualFold(dt, "txid_snapshot"):
			return RandomTXID(), nil

		// Random GeoMetric data
		case StringContains(dt, geoDataTypekeywords):
			var randomInt int
			if dt == "path" || dt == "polygon" {
				randomInt, _ = RandomInt(1, 5)
			} else {
				randomInt, _ = RandomInt(1, 2)
			}
			return RandomGeometricData(randomInt, dt), nil
		default:
			return "", fmt.Errorf("Unsupported datatypes found: %v", dt)
	}

	return "", nil
}
