package core

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"math/rand"
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
				ArrayArgs = map[string]interface{}{"intmin": -intranges[nonArraydt], "intmax": intranges[nonArraydt]}
				value, err := ArrayGenerator("int")
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
				ArrayArgs["strlen"] = l
				value, _ := ArrayGenerator("string")
				return value, nil
			} else {
				value := RandomString(l)
				return value, nil
			}

		// Generate Random date
		case strings.HasPrefix(dt, "date"):
			if strings.HasSuffix(dt, "[]") {
				ArrayArgs = map[string]interface{}{"fromyear": fromyear, "toyear": toyear}
				value, err := ArrayGenerator("date")
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
				ArrayArgs = map[string]interface{}{"fromyear": fromyear, "toyear": toyear}
				value, err := ArrayGenerator("timestamp")
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

		/* 
		=== Updated at 2019-06-26 ===
		Added below code to handle data type from timestamp(0) to timestampe(6) with/without time zone
		The data type can be find in this doc: https://gpdb.docs.pivotal.io/5200/ref_guide/data_types.html 
		*/
		case regexp.MustCompile(`timestamp\([0-6]\) without time zone`).MatchString(dt), regexp.MustCompile(`timestamp\([0-6]\) with time zone`).MatchString(dt):
			value, err := RandomTimestamp(fromyear, toyear)     // get a random timestamp with format like: 2018-04-10 01:19:22
			if err != nil {
				return "", fmt.Errorf("Build Timestamp[p] without timezone: %v", err)
			}
	    	ts_reg := regexp.MustCompile(`\([0-6]\)`)
			decimal, _ := strconv.Atoi( strings.Split(ts_reg.FindString(dt),"")[1] )  // capture the decimal in timestamp[x]
	    	var timestamp_decimal string
			for i := 0; i < decimal; i++ {
				timestamp_decimal = timestamp_decimal + strconv.Itoa(rand.Intn(9)) // use rand() to generate random decimal in timestamp
			}
			if len(timestamp_decimal) > 0 {
				value = value + "." + timestamp_decimal
			}
			return value,nil
		/* End of Updated */
		
		// Generate Random timestamp with timezone
		case strings.HasPrefix(dt, "timestamp with time zone"):
			if strings.HasSuffix(dt, "[]") {
				ArrayArgs = map[string]interface{}{"fromyear": fromyear, "toyear": toyear}
				value, err := ArrayGenerator("timestamptz")
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
				ArrayArgs = map[string]interface{}{"fromyear": fromyear, "toyear": toyear}
				value, err := ArrayGenerator("time")
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
				ArrayArgs = map[string]interface{}{"fromyear": fromyear, "toyear": toyear}
				value, err := ArrayGenerator("timetz")
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
		case StringHasPrefix(dt, ipkeywords):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("IP")
				return value, nil
			} else {
				return RandomIP(), nil
			}

		// Generate Random boolean
		case strings.HasPrefix(dt, "boolean"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("bool")
				return value, nil
			} else {
				return RandomBoolean(), nil
			}

		// Generate Random text
		case strings.HasPrefix(dt, "text"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("text")
				return value, nil
			} else {
				return RandomParagraphs(), nil
			}

		// Generate Random text & bytea
		case strings.EqualFold(dt, "bytea"):
			return RandomBytea(1024 * 1024), nil

		// Generate Random float values
		case StringHasPrefix(dt, floatkeywords):
			if strings.HasSuffix(dt, "[]") { // Float array
				ArrayArgs = map[string]interface{}{"floatmin": 1, "floatmax": intranges["smallint"], "floatprecision": 3}
				value, err := ArrayGenerator("float")
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
				ArrayArgs = map[string]interface{}{"floatmin": 0, "floatmax": max, "floatprecision": precision}
				value, err := ArrayGenerator("float")
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
				ArrayArgs["bitlen"] = l
				value, err := ArrayGenerator("bit")
				if err != nil {
					return "", fmt.Errorf("Build bit array: %v", err)
				}
				return value, nil
			} else {
				value := RandomBit(l)
				return value, nil
			}

		// Random UUID generator
		case strings.HasPrefix(dt, "uuid"):
			if strings.HasSuffix(dt, "[]") {
				uuid, err := ArrayGenerator("uuid")
				if err != nil {
					return "", fmt.Errorf("Build UUID Array: %v", err)
				}
				return uuid, nil
			} else {
				uuid, err := RandomUUID()
				if err != nil {
					return "", fmt.Errorf("Build UUID: %v", err)
				}
				return uuid, nil
			}

		// Random MacAddr Generator
		case strings.HasPrefix(dt, "macaddr"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("macaddr")
				return value, nil
			} else {
				return RandomMacAddress(), nil
			}

		// Random Json
		case strings.HasPrefix(dt, "json"):
			if strings.HasSuffix(dt, "[]") {
				return JsonXmlArrayGenerator("json"), nil
			} else {
				return RandomJson(false), nil
			}

		// Random XML
		case strings.HasPrefix(dt, "xml"):
			if strings.HasSuffix(dt, "[]") {
				return JsonXmlArrayGenerator("xml"), nil
			} else {
				return RandomXML(false), nil
			}

		// Random Text Search Query
		case strings.HasPrefix(dt, "tsquery"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("tsquery")
				return value, nil
			} else {
				return RandomTSQuery(), nil
			}

		// Random Text Search Vector
		case strings.HasPrefix(dt, "tsvector"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("tsquery")
				return value, nil
			} else {
				return RandomTSVector(), nil
			}

		// Random Log Sequence number
		case strings.HasPrefix(dt, "pg_lsn"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("pg_lsn")
				return value, nil
			} else {
				return RandomLSN(), nil
			}

		// Random Log Sequence number
		case strings.HasPrefix(dt, "txid_snapshot"):
			if strings.HasSuffix(dt, "[]") {
				value, _ := ArrayGenerator("txid_snapshot")
				return value, nil
			} else {
				return RandomTXID(), nil
			}

		// Random GeoMetric data
		case StringHasPrefix(dt, geoDataTypekeywords):
			var randomInt int
			if dt == "path" || dt == "polygon" {
				randomInt, _ = RandomInt(1, 5)
			} else {
				randomInt, _ = RandomInt(1, 2)
			}
			if strings.HasSuffix(dt, "[]") {
				dtype := strings.Replace(dt, "[]", "", 1)
				value := GeometricArrayGenerator(randomInt, dtype)
				return value, nil
			} else {
				return RandomGeometricData(randomInt, dt, false), nil
			}


		// If there is no datatype found then send the below message
		default:
			return "", fmt.Errorf("Unsupported datatypes found: %v", dt)
	}

	return "", nil
}
