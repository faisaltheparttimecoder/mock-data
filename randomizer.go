package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

// Set the seed value of the random generator
var r *rand.Rand

func init() {
	// nolint:gosec
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString: Random String generator
func RandomString(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// RandomInt: Random Number generator based on the min and max specified
func RandomInt(min, max int) int {
	if min >= max {
		return r.Intn(min-max) + min
	}
	return r.Intn(max-min) + min
}

// RandomBytea: Random Bytea data
func RandomBytea(maxlen int) []byte {
	result := make([]byte, r.Intn(maxlen)+1)
	for i := range result {
		result[i] = byte(r.Intn(255))
	}
	return result
}

// RandomFloat: Random Float generator based on precision specified
func RandomFloat(min, max, precision int) float64 {
	output := math.Pow(10, float64(precision))
	randNumber := float64(min) + r.Float64()*float64(max-min)*100
	return math.Round(randNumber) / output
}

// RandomCalenderDateTime: Random calender date time generator
func RandomCalenderDateTime(fromyear, toyear int) (time.Time, error) {
	if fromyear > toyear {
		return time.Now(), errors.New("number of years behind is greater than number of years in future")
	}
	min := time.Now().AddDate(fromyear, 0, 0).Unix()
	max := time.Now().AddDate(toyear, 0, 0).Unix()
	delta := max - min
	sec := r.Int63n(delta) + min
	return time.Unix(sec, 0), nil
}

// RandomDate: Random date
func RandomDate(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02"), nil
}

// RandomTimestamp: Random Timestamp without time zone
func RandomTimestamp(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02 15:04:05"), nil
}

// RandomTimeStampTz: Random Timestamp with time zone
func RandomTimeStampTz(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02 15:04:05.000000"), nil
}

// RandomTimeStampTzWithDecimals: Random Timestamp with decimals
func RandomTimeStampTzWithDecimals(fromyear, toyear, decimal int) (string, error) {
	var timestampDecimal string
	d, err := RandomTimestamp(fromyear, toyear)
	if err != nil {
		return "", fmt.Errorf("randomizer with timestamp[p] without timezone failed: %w", err)
	}
	// use rand() to generate random decimal in timestamp
	for i := 0; i < decimal; i++ {
		timestampDecimal = timestampDecimal + strconv.Itoa(r.Intn(9))
	}
	if len(timestampDecimal) > 0 {
		d = d + "." + timestampDecimal
	}
	return d, nil
}

// RandomTime: Random Time without time zone
func RandomTime(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("15:04:05"), nil
}

// RandomTimeTz: Random Timestamp without time zone
func RandomTimeTz(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("15:04:05.000000"), nil
}

// RandomBoolean: Random bool generator based on if number is even or not
func RandomBoolean() bool {
	number := RandomInt(1, 9999)
	return number%2 == 0
}

// RandomParagraphs: Random Paragraphs
func RandomParagraphs() string {
	n := RandomInt(1, 5)
	return fake.ParagraphsN(n)
}

// RandomCiText: Random CiText generator
func RandomCiText() string {
	return strings.Title(fake.Words())
}

// RandomIP: Random IPv6 & IPv4 Address
func RandomIP() string {
	number := RandomInt(1, 9999)
	var ip string
	if ip = fake.IPv6(); number%2 == 0 {
		ip = fake.IPv4()
	}
	return ip
}

// RandomBit: Random bit
func RandomBit(max int) string {
	var bitValue string
	for i := 0; i < max; i++ {
		if RandomBoolean() {
			bitValue = bitValue + "1"
		} else {
			bitValue = bitValue + "0"
		}
	}
	return bitValue
}

// RandomUUID: Random UUID
func RandomUUID() string {
	return uuid.New().String()
}

// RandomMacAddress: Random Mac Address
func RandomMacAddress() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1))
}

// RandomTSQuery: Random Text Search Query
func RandomTSQuery() string {
	number := RandomInt(1, 9999)
	switch number % 5 { // TODO: replace magic number 5 to symbol constant. What is mean 5? Why exactly 5?
	case 0:
		return fake.WordsN(1) + " & " + fake.WordsN(1)
	case 1: // TODO: replace magic number to symbol constant. What is mean 1 or 2 or 3?
		return fake.WordsN(1) + " | " + fake.WordsN(1)
	case 2: // TODO: replace magic number to symbol constant
		return fake.WordsN(1) + " | " + fake.WordsN(1)
	case 3: // TODO: replace magic number to symbol constant
		return fake.WordsN(1) + " & " + fake.WordsN(1) + "  & ! " + fake.WordsN(1)
	default:
		return fake.WordsN(1) + " & ( " + fake.WordsN(1) + " | " + fake.WordsN(1) + " )"
	}
}

// RandomTSVector: Random Text Search Query
func RandomTSVector() string {
	return fake.SentencesN(fake.Day())
}

// RandomGeometricData: Random Geometric data
func RandomGeometricData(randomInt int, GeoMetry string, IsItArray bool) string {
	var data string
	if GeoMetry == "point" { // Syntax for point data type
		data = fmt.Sprintf("%d,%d",
			RandomInt(1, 999), RandomInt(1, 999))
		return FormatForArray(data, IsItArray)
	} else if GeoMetry == "circle" { // Syntax for circle data type
		data = fmt.Sprintf("<(%d,%d),%d>",
			RandomInt(1, 999), RandomInt(1, 999), RandomInt(1, 999))
		return FormatForArray(data, IsItArray)
	} else { // Syntax for rest
		data = fmt.Sprintf("%d,%d,%d,%d",
			RandomInt(1, 999), RandomInt(1, 999),
			RandomInt(1, 999), RandomInt(1, 999))
		return FormatForArray(data, IsItArray)
	}
}

// RandomLSN: Random Log Sequence Number
func RandomLSN() string {
	return fmt.Sprintf("%02x/%02x",
		RandomString(1), RandomString(4))
}

// RandomTXID: Random transaction XID
func RandomTXID() string {
	x, _ := strconv.Atoi(fake.DigitsN(8))
	y, _ := strconv.Atoi(fake.DigitsN(8))
	var z string
	if z = fmt.Sprintf("%v:%v:", x, y); x > y { // left side of ":" should be always less than right side
		z = fmt.Sprintf("%v:%v:", y, x)
	}
	return z
}

// RandomJSON: Random JSON generator
func RandomJSON(IsItArray bool) string {
	jsonData := fmt.Sprintf(JsonSkeleton(), RandomString(24),
		fake.DigitsN(10), RandomUUID(), strconv.FormatBool(RandomBoolean()), fake.Digits(), fake.DigitsN(2),
		fake.DomainName(), fake.WordsN(1), fake.DigitsN(2), fake.UserName(), fake.Color(), fake.FullName(),
		fake.Gender(), fake.Company(), fake.EmailAddress(), fake.Phone(), fake.StreetAddress(), fake.Zip(),
		fake.State(), fake.Country(), fake.WordsN(12), RandomIP(), fake.JobTitle(),
		strconv.Itoa(fake.Year(2000, 2050)), strconv.Itoa(fake.MonthNum()), strconv.Itoa(fake.Day()),
		fake.DigitsN(2), fake.DigitsN(2), fake.DigitsN(2), fake.DigitsN(1), fake.DigitsN(2),
		fake.DigitsN(2), fake.DigitsN(6), fake.DigitsN(2), fake.DigitsN(6), fake.WordsN(1),
		fake.WordsN(1), fake.WordsN(1), fake.WordsN(1), fake.WordsN(1), fake.WordsN(1),
		fake.WordsN(1), fake.DigitsN(2), fake.FullName(), fake.DigitsN(2), fake.FullName(),
		fake.DigitsN(2), fake.FullName(), fake.Sentence(),
		fake.Brand())
	if IsItArray {
		return strings.Replace(jsonData, "\"", "\\\"", -1)
	}
	return jsonData
}

// RandomXML: Random XML Generator
func RandomXML(IsItArray bool) string {
	xmlData := fmt.Sprintf(XMLSkeleton(), fake.Digits(), fake.DomainName(),
		fake.DigitsN(4), fake.WordsN(1), fake.FullName(), fake.FullName(), fake.StreetAddress(), fake.City(),
		fake.Country(), fake.EmailAddress(), fake.Phone(), fake.Title(), fake.Sentences(), fake.Digits(), fake.Color(),
		fake.Digits(), fake.DigitsN(2), fake.Title(), fake.Digits(), fake.Digits(), fake.DigitsN(2))
	if IsItArray {
		return strings.Replace(xmlData, "\"", "\\\"", -1)
	}
	return xmlData
}

// RandomPickerFromArray: Picks Random Value from any array
func RandomPickerFromArray(a []string) string {
	if len(a) == 0 {
		return ""
	}
	return a[RandomValueFromLength(len(a))]
}

// RandomValueFromLength: Random value from length
func RandomValueFromLength(i int) int {
	if i == 0 {
		return 0
	}
	return r.Int() % i
}
