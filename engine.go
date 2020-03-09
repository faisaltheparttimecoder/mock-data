package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ranges of dates
	fromYear = -10
	toYear   = 10

	// Time data types
	intervalKeywords = []string{"interval", "time without time zone"}

	// Networking data types
	ipKeywords = []string{"inet", "cidr"}

	// Integer data types
	intKeywords = []string{"smallint", "integer", "bigint", "oid"}
	intRanges   = map[string]int{"smallint": 2767, "integer": 7483647, "bigint": 372036854775807, "oid": 7483647}

	// Decimal data types
	floatKeywords = []string{"double precision", "real", "money"}

	// Geometry data types
	geoDataTypekeywords = []string{"path", "polygon", "line", "lseg", "box", "circle", "point"}
)

// Data Generator
// It provided random data based on data types.
func BuildData(dt string) (interface{}, error) {
	var value interface{}
	if StringHasPrefix(dt, intKeywords) { // Integer builder
		return buildInteger(dt)
	} else if strings.HasPrefix(dt, "character") { // String builder
		return buildCharacter(dt)
	} else if strings.HasPrefix(dt, "date") { // Date builder
		return buildDate(dt)
	} else if strings.HasPrefix(dt, "timestamp") { // Timestamp builder
		return buildTimeStamp(dt)
	} else if StringHasPrefix(dt, intervalKeywords) { // Generate Random time without timezone
		return buildInterval(dt)
	} else if strings.HasPrefix(dt, "time with time zone") { // Generate Random time with timezone
		return buildTimeWithTz(dt)
	} else if StringHasPrefix(dt, ipKeywords) { // Generate Random ips
		return buildIps(dt)
	} else if strings.HasPrefix(dt, "boolean") { // Generate Random boolean
		return buildBoolean(dt)
	} else if strings.HasPrefix(dt, "text") { // Generate Random text
		return buildText(dt)
	} else if strings.EqualFold(dt, "bytea") { // Generate Random bytea
		return buildBytea(dt)
	} else if StringHasPrefix(dt, floatKeywords) { // Generate Random float values
		return buildFloat(dt)
	} else if strings.HasPrefix(dt, "numeric") { // Generate Random numeric values with precision
		return buildNumeric(dt)
	} else if strings.HasPrefix(dt, "bit") { // Random bit generator
		return buildBit(dt)
	} else if strings.HasPrefix(dt, "uuid") { // Random UUID generator
		return buildUuid(dt)
	} else if strings.HasPrefix(dt, "macaddr") { // Random MacAddr Generator
		return buildMacAddr(dt)
	} else if strings.HasPrefix(dt, "json") { // Random Json
		return buildJson(dt)
	} else if strings.HasPrefix(dt, "xml") { // Random Xml
		return buildXml(dt)
	} else if strings.HasPrefix(dt, "tsquery") { // Random Text Search Query
		return buildTsQuery(dt)
	} else if strings.HasPrefix(dt, "tsvector") { // Random Text Search Vector
		return buildTsVector(dt)
	} else if strings.HasPrefix(dt, "pg_lsn") { // Random Log Sequence number
		return buildLseg(dt)
	} else if strings.HasPrefix(dt, "txid_snapshot") { // Random transaction XID snapshot
		return buildTxidSnapShot(dt)
	} else if StringHasPrefix(dt, geoDataTypekeywords) { // Random GeoMetric data
		return buildGeometry(dt)
	} else { // if these are not the defaults, the ony custom we allow is enum data type, check if its them
		return buildEnumDatatypes(dt)
	}
	return value, nil
}

// Build Integer
func buildInteger(dt string) (interface{}, error) {
	isItArray, t := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("int", dt, -intRanges[t], intRanges[t])
	}
	value := RandomInt(-intRanges[dt], intRanges[dt])
	return value, nil
}

// Character Builder
func buildCharacter(dt string) (interface{}, error) {
	l, err := CharLen(dt)
	if err != nil {
		return "", fmt.Errorf("getting character length: %v", err)
	}
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("string", dt, 0, l)
	}
	return RandomString(l), nil
}

// Date Builder
func buildDate(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("date", dt, fromYear, toYear)
	}
	return RandomDate(fromYear, toYear)
}

// Time Builder
func buildTimeWithTz(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("timetz", dt, fromYear, toYear)
	}
	return RandomTimeTz(fromYear, toYear)
}

// Ip's builder
func buildIps(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("IP", dt, 0, 0)
	}
	return RandomIP(), nil
}

// Boolean builder
func buildBoolean(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("bool", dt, 0, 0)
	}
	return RandomBoolean(), nil
}

// Text builder
func buildText(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("text", dt, 0, 0)
	}
	return RandomParagraphs(), nil
}

// Bytea builder
func buildBytea(dt string) (interface{}, error) {
	return RandomBytea(1024 * 1024), nil
}

// Float builder
func buildFloat(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("float", dt, 1, intRanges["smallint"])
	}
	return RandomFloat(1, intRanges["smallint"], 3), nil
}

// Numeric values with precision builder
func buildNumeric(dt string) (interface{}, error) {
	max, precision, err := findNumberPrecision(dt)
	value := RandomFloat(1, max, precision)
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("numericFloat", dt, 0, max)
	}
	return TruncateFloat(value, max, precision), err
}

// Bit builder
func buildBit(dt string) (interface{}, error) {
	l, err := CharLen(dt)
	if err != nil {
		return "", fmt.Errorf("build bit: %v", err)
	}
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("bit", dt, 0, l)
	}
	return RandomBit(l), nil
}

// Uuid Builder
func buildUuid(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("uuid", dt, 0, 0)
	}
	return RandomUUID(), nil
}

// Mac Addr builder
func buildMacAddr(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("macaddr", dt, 0, 0)
	}
	return RandomMacAddress(), nil
}

// Json builder
func buildJson(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return JsonXmlArrayGenerator("json"), nil
	}
	return RandomJson(false), nil
}

// Xml Builder
func buildXml(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return JsonXmlArrayGenerator("xml"), nil
	}
	return RandomXML(false), nil
}

// Ts Query Builder
func buildTsQuery(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("tsquery", dt, 0, 0)
	}
	return RandomTSQuery(), nil
}

// Ts Vector Builder
func buildTsVector(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("tsvector", dt, 0, 0)
	}
	return RandomTSVector(), nil
}

// Log Sequence Builder
func buildLseg(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("pg_lsn", dt, fromYear, toYear)
	}
	return RandomLSN(), nil
}

// Transaction Xid Builder
func buildTxidSnapShot(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("txid_snapshot", dt, fromYear, toYear)
	}
	return RandomTXID(), nil
}

// GeoMetric data builder
func buildGeometry(dt string) (interface{}, error) {
	var randomInt int
	if dt == "path" || dt == "polygon" {
		randomInt = RandomInt(1, 5)
	} else {
		randomInt = RandomInt(1, 2)
	}
	isItArray, t := isDataTypeAnArray(dt)
	if isItArray {
		return GeometricArrayGenerator(randomInt, t), nil
	}
	return RandomGeometricData(randomInt, dt, false), nil
}

// TimeStamp Builder
func buildTimeStamp(dt string) (interface{}, error) {
	isItArray, t := isDataTypeAnArray(dt)
	if t == "timestamp without time zone" { // Without time zone
		if isItArray {
			return ArrayGenerator("timestamp", dt, fromYear, toYear)
		}
		return RandomTimestamp(fromYear, toYear)
	} else if t == "timestamp with time zone" { // With time zone
		if isItArray {
			return ArrayGenerator("timestamptz", dt, fromYear, toYear)
		}
		return RandomTimeStampTz(fromYear, toYear)
	} else if regexp.MustCompile(`timestamp\([0-6]\) without time zone`).MatchString(dt) ||
		regexp.MustCompile(`timestamp\([0-6]\) with time zone`).MatchString(dt) { // time zone with precision
		if isItArray {
			return ArrayGenerator("timestamptzWithDecimals", dt, fromYear, toYear)
		}
		return RandomTimeStampTzWithDecimals(fromYear, toYear, findTimeStampDecimal(dt))
	}
	return "", nil
}

// Interval builder
func buildInterval(dt string) (interface{}, error) {
	isItArray, _ := isDataTypeAnArray(dt)
	if isItArray {
		return ArrayGenerator("time", dt, fromYear, toYear)
	}
	return RandomTime(fromYear, toYear)
}

// Find the number of decimal points needed
func findTimeStampDecimal(dt string) int {
	tsReg := regexp.MustCompile(`\([0-6]\)`)
	decimal, _ := strconv.Atoi(strings.Split(tsReg.FindString(dt), "")[1])
	return decimal
}

// Find number precision
func findNumberPrecision(dt string) (int, int, error) {
	max, precision, err := FloatPrecision(dt)
	if err != nil {
		return 0, 0, fmt.Errorf("build numeric: %v", err)
	}
	return max, precision, nil
}

// Does the data type have array
func isDataTypeAnArray(dt string) (bool, string) {
	if strings.HasSuffix(dt, "[]") {
		t := strings.Replace(dt, "[]", "", 1)
		return true, t
	}
	return false, dt
}

// Random array generator for array datatypes
func ArrayGenerator(dt, originalDt string, min, max int) (string, error) {
	maxValues := RandomInt(1, 6) // Getting the value of iterators
	maxIteration := RandomInt(1, 3)
	var resultArrayCollector []string

	// Create a array of results
	for i := 0; i < maxIteration; i++ { // Max number of arrays
		var resultArray []string
		for j := 0; j < maxValues; j++ { // max number of values in a array.
			value, err := randomDataByDataTypeForArray(dt, originalDt, min, max)
			if err != nil {
				return "", fmt.Errorf("error when generating array for datatype %s, err: %v", dt, err)
			}
			resultArray = append(resultArray, value)
		}
		resultArrayCollector = append(resultArrayCollector, strings.Join(resultArray, ","))
	}
	return fmt.Sprintf("{%s}", strings.Join(resultArrayCollector, ",")), nil
}

// Send in appropriate random value based on data type
func randomDataByDataTypeForArray(dt, originalDt string, min, max int) (string, error) {
	if dt == "int" {
		return strconv.Itoa(RandomInt(min, max)), nil
	} else if dt == "string" {
		return RandomString(max), nil
	} else if dt == "date" {
		return RandomDate(min, max)
	} else if dt == "timestamp" {
		return RandomTimestamp(min, max)
	} else if dt == "timestamptz" {
		return RandomTimeStampTz(min, max)
	} else if dt == "timestamptzWithDecimals" {
		return RandomTimeStampTzWithDecimals(min, max, findTimeStampDecimal(originalDt))
	} else if dt == "time" {
		return RandomTime(min, max)
	} else if dt == "float" {
		value := RandomFloat(0, max, 3)
		return fmt.Sprintf("%v", TruncateFloat(value, max, 3)), nil
	} else if dt == "numericFloat" {
		m, precision, err := FloatPrecision(dt)
		value := RandomFloat(1, m, precision)
		return fmt.Sprintf("%v", value), err
	} else if dt == "bit" {
		return RandomBit(max), nil
	} else if dt == "text" {
		return fake.WordsN(1), nil
	} else if dt == "timetz" {
		return RandomTimeTz(min, max)
	} else if dt == "bool" {
		if RandomBoolean() {
			return "true", nil
		} else {
			return "false", nil
		}
	} else if dt == "IP" {
		return RandomIP(), nil
	} else if dt == "macaddr" {
		return RandomMacAddress(), nil
	} else if dt == "uuid" {
		return RandomUUID(), nil
	} else if dt == "txid_snapshot" {
		return RandomTXID(), nil
	} else if dt == "pg_lsn" {
		return RandomLSN(), nil
	} else if dt == "tsquery" {
		return RandomTSQuery(), nil
	} else if dt == "tsvector" {
		return RandomTSVector(), nil
	} else {
		return "", fmt.Errorf("unsupported datatypes found in array %v", dt)
	}
}

// Random geometric array generators
func GeometricArrayGenerator(maxInt int, geometryType string) string {
	// Getting the value of iterators
	maxIterations := RandomInt(1, 6)
	var resultArray []string

	if geometryType == "box" {
		value := RandomGeometricData(maxInt, geometryType, false)
		resultArray = append(resultArray, value)
	} else {
		for i := 0; i < maxIterations; i++ { // Max number of arrays
			value := RandomGeometricData(maxInt, geometryType, true)
			resultArray = append(resultArray, value)
		}
	}

	return fmt.Sprintf("{%s}", strings.Join(resultArray, ","))
}

// Random XML & Json array generators.
func JsonXmlArrayGenerator(dt string) string {
	// Getting the value of iterators
	maxIterations := RandomInt(1, 6)
	var resultArray []string
	var value string
	for i := 0; i < maxIterations; i++ { // Max number of arrays

		switch dt { // Choose the appropriate random data generators
		case "json":
			value = fmt.Sprintf("\"%s\"", RandomJson(true))
		case "xml":
			value = fmt.Sprintf("\"%s\"", RandomXML(true))
		}

		resultArray = append(resultArray, value)
	}
	return fmt.Sprintf("{%s}", strings.Join(resultArray, ","))
}

// Enum datatypes
func buildEnumDatatypes(dt string) (string, error) {
	// Check if the data type is ENUM
	enumOutput := checkEnumDatatype(dt)

	// if there are none then pass in the error back to user
	if len(enumOutput) <= 0 {
		return "", fmt.Errorf("unsupported datatypes found: %v", dt)
	}

	// found some output, lets pick some random value
	n := RandomValueFromLength(len(enumOutput))
	return enumOutput[n].EnumValue, nil
}
