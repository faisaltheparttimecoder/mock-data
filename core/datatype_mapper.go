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
	var floatkeywords = []string{"double precision", "numeric", "real", "money"}

	// Geometry datatypes
	var geoDataTypekeywords = []string{"path", "polygon", "line", "lseg", "box", "circle", "point"}

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

		// Generate Random date
		case strings.EqualFold(dt, "date"):
			value, err := RandomDate(fromyear, toyear)
			if err != nil {
				return "", fmt.Errorf("Build Date: %v", err)
			}
			return value, nil

		// Generate Random timestamp without timezone
		case strings.EqualFold(dt, "timestamp without time zone"):
			value, err := RandomTimestamp(fromyear, toyear)
			if err != nil {
				return "", fmt.Errorf("Build Timestamp  without timezone: %v", err)
			}
			return value, nil

		// Generate Random timestamp with timezone
		case strings.EqualFold(dt, "timestamp with time zone"):
			value, err := RandomTimestamptz(fromyear, toyear)
			if err != nil {
				return "", fmt.Errorf("Build Timestamp with timezone: %v", err)
			}
			return value, nil

		// Generate Random time without timezone
		case StringContains(dt, Intervalkeywords):
			value, err := RandomTime(fromyear, toyear)
			if err != nil {
				return "", fmt.Errorf("Build Time without timezone: %v", err)
			}
			return value, nil

		// Generate Random time with timezone
		case strings.EqualFold(dt, "time with time zone"):
			value, err := RandomTimetz(fromyear, toyear)
			if err != nil {
				return "", fmt.Errorf("Build Time with timezone: %v", err)
			}
			return value, nil

		// Generate Random ips
		case StringContains(dt, ipkeywords):
			return RandomIP(), nil

		// Generate Random boolean
		case strings.EqualFold(dt, "boolean"):
			return RandomBoolean(), nil

		// Generate Random text & bytea
		// not sure what the best way to generate a bytea data, so lets default to paragraph
		case strings.EqualFold(dt, "text"):
			return RandomParagraphs(), nil

		case strings.EqualFold(dt, "bytea"):
			return RandomBytea(1024 * 1024), nil

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
