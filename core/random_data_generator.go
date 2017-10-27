package core

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/icrowley/fake"
)

// Random text generator based on the length of string needed
var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Random Bytea data
func RandomBytea(maxlen int) []byte {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, r.Intn(maxlen)+1)
	for i := range result {
		result[i] = byte(r.Intn(255))
	}
	return result
}

// Random String
func RandomString(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// Random Number generator based on the min and max specified
func RandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, errors.New("Max value is greater or equal to Min value, cannot generate data within these ranges")
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + max, nil
}

// Random Float generator based on precision specified
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RandomFloat(min, max, precision int) (float64, error) {
	output := math.Pow(10, float64(precision))
	randNumber, err := RandomInt(min, max)
	if err != nil {
		return 0.0, err
	}
	return float64(round(float64(randNumber)/rand.Float64()*output)) / output, nil
}

// Random calender date time generator
func RandomCalenderDateTime(fromyear, toyear int) (time.Time, error) {
	if fromyear > toyear {
		return time.Now(), errors.New("Number of years behind is greater than number of years in future")
	}
	min := time.Now().AddDate(fromyear, 0, 0).Unix()
	max := time.Now().AddDate(toyear, 0, 0).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0), nil
}

// Random date
func RandomDate(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02"), nil
}

// Random Timestamp without time zone
func RandomTimestamp(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02 15:04:05"), nil
}

// Random Timestamp with time zone
func RandomTimestamptz(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("2006-01-02 15:04:05.000000"), nil
}

// Random Time without time zone
func RandomTime(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("15:04:05"), nil
}

// Random Timestamp without time zone
func RandomTimetz(fromyear, toyear int) (string, error) {
	timestamp, err := RandomCalenderDateTime(fromyear, toyear)
	if err != nil {
		return "", err
	}
	return timestamp.Format("15:04:05.000000"), nil
}

// Random bool generator based on if number is even or not
func RandomBoolean() bool {
	number, _ := RandomInt(1, 9999)
	if number%2 == 0 {
		return true
	} else {
		return false
	}
}

// Random Paragraphs
func RandomParagraphs() string {
	n, _ := strconv.Atoi(fake.DigitsN(1))
	return fake.ParagraphsN(n)
}

// Random IPv6 & IPv4 Address
func RandomIP() string {
	number, _ := RandomInt(1, 9999)
	if number%2 == 0 {
		return fake.IPv4()
	} else {
		return fake.IPv6()
	}
}

// Random bit
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

// Random UUID
func RandomUUID() (string, error) {
	// To generate random UUID, we will use unix tool "uuidgen" (unix utility)
	uuidString, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", fmt.Errorf("Unable to run uuidgen to generate UUID data: %v", err)
	}
	return strings.TrimSpace(string(uuidString)), nil
}

// Random Mac Address
func RandomMacAddress() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1),
		RandomString(1), RandomString(1))
}

// Random Text Search Query
func RandomTSQuery() string {
	number, _ := RandomInt(1, 9999)
	number = number % 5
	if number == 0 {
		return fake.WordsN(1) + " & " + fake.WordsN(1)
	} else if number == 1 {
		return fake.WordsN(1) + " | " + fake.WordsN(1)
	} else if number == 2 {
		return " ! " + fake.WordsN(1) + " & " + fake.WordsN(1)
	} else if number == 3 {
		return fake.WordsN(1) + " & " + fake.WordsN(1) + "  & ! " + fake.WordsN(1)
	} else {
		return fake.WordsN(1) + " & ( " + fake.WordsN(1) + " | " + fake.WordsN(1) + " )"
	}
	return ""
}

// Random Text Search Query
func RandomTSVector() string {
	return fake.SentencesN(fake.Day())
}

// Random Geometric data
func RandomGeometricData(randomInt int, GeoMetry string) string {
	var geometry []string
	if GeoMetry == "point" {
		return "(" + fake.DigitsN(2) + "," + fake.DigitsN(3) + ")"
	} else if GeoMetry == "circle" {
		return "(" + fake.DigitsN(2) + "," + fake.DigitsN(3) + ")," + fake.DigitsN(2) + ")"
	} else {
		for i := 0; i < randomInt; i++ {
			x, _ := RandomFloat(1, 10, 2)
			y, _ := RandomFloat(1, 10, 2)
			geometry = append(geometry, "("+fmt.Sprintf("%v", x)+","+fmt.Sprintf("%v", y)+")")
		}

		return "(" + strings.Join(geometry, ",") + ")"
	}
	return ""
}


// Random Log Sequence Number
func RandomLSN() string {
	return fmt.Sprintf("%02x/%02x",
		RandomString(1), RandomString(4))
}

// Random transaction XID
func RandomTXID() string {
	x, _ := strconv.Atoi(fake.DigitsN(8))
	y, _ := strconv.Atoi(fake.DigitsN(8))
	if x > y { // left side of ":" should be always less than right side
		return fmt.Sprintf("%v:%v:", y, x)
	} else {
		return fmt.Sprintf("%v:%v:", x, y)
	}
	return ""
}

// Random JSON generator
func RandomJson() string {
	return "{" +
		"    \"_id\": \"" + RandomString(24) + "\"," +
		"    \"index\": \"" + fake.DigitsN(10) + "\"," +
		"    \"guid\": \"" + RandomString(8) + "-" + RandomString(4) + "-" + RandomString(4) + "-" + RandomString(4) + "-" + RandomString(12) + "\"," +
		"    \"isActive\": \"" + strconv.FormatBool(RandomBoolean()) + "\"," +
		"    \"balance\": \"$" + fake.Digits() + "." + fake.DigitsN(2) + "\"," +
		"    \"website\": \"https://" + fake.DomainName() + "/" + fake.WordsN(1) + "\"," +
		"    \"age\": \"" + fake.DigitsN(2) + "\"," +
		"    \"username\": \"" + fake.UserName() + "\"," +
		"    \"eyeColor\": \"" + fake.Color() + "\"," +
		"    \"name\": \"" + fake.FullName() + "\"," +
		"    \"gender\": \"" + fake.Gender() + "\"," +
		"    \"company\": \"" + fake.Company() + "\"," +
		"    \"email\": \"" + fake.EmailAddress() + "\"," +
		"    \"phone\": \"" + fake.Phone() + "\"," +
		"    \"address\": \"" + fake.StreetAddress() + "\"," +
		"    \"zipcode\": \"" + fake.Zip() + "\"," +
		"    \"state\": \"" + fake.State() + "\"," +
		"    \"country\": \"" + fake.Country() + "\"," +
		"    \"about\": \"" + fake.WordsN(12) + "\"," +
		"    \"Machine IP\": \"" + RandomIP() + "\"," +
		"    \"job title\": \"" + fake.JobTitle() + "\"," +
		"    \"registered\": \"" + strconv.Itoa(fake.Year(2000, 2050)) + "-" + strconv.Itoa(fake.MonthNum()) + "-" + strconv.Itoa(fake.Day()) + "T" + fake.DigitsN(2) + ":" + fake.DigitsN(2) + ":" + fake.DigitsN(2) + " -" + fake.DigitsN(1) + ":" + fake.DigitsN(2) + "\"," +
		"    \"latitude\": \"" + fake.DigitsN(2) + "." + fake.DigitsN(6) + "\"," +
		"    \"longitude\": \"" + fake.DigitsN(2) + "." + fake.DigitsN(6) + "\"," +
		"    \"tags\": [" +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"," +
		"      \"" + fake.WordsN(1) + "\"" +
		"    ]," +
		"    \"friends\": [" +
		"      {" +
		"        \"id\": \"" + fake.DigitsN(2) + "\"," +
		"        \"name\": \"" + fake.FullName() + "\"" +
		"      }," +
		"      {" +
		"        \"id\": \"" + fake.DigitsN(2) + "\"," +
		"        \"name\": \"" + fake.FullName() + "\"" +
		"      }," +
		"      {" +
		"        \"id\": \"" + fake.DigitsN(2) + "\"," +
		"        \"name\": \"" + fake.FullName() + "\"" +
		"      }" +
		"    ]," +
		"    \"greeting\": \"" + fake.Sentence() + "\"," +
		"    \"favoriteBrand\": \"" + fake.Brand() + "\"" +
		"  }"
}

// Random XML Generator
func RandomXML() string {
	return "<?xml version=\"" + fake.DigitsN(1) + "." + fake.DigitsN(1) + "\" encoding=\"UTF-8\"?>" +
		"<shiporder orderid=\"" + fake.Digits() + "\"" +
		" xmlns:xsi=\"http://" + fake.DomainName() + "/" + fake.DigitsN(4) + "/" + fake.WordsN(1) + "\" " +
		"xsi:noNamespaceSchemaLocation=\"shiporder.xsd\">" +
		"  <orderperson>" + fake.FullName() + "</orderperson>" +
		"  <shipto>" +
		"    <name>" + fake.FullName() + "</name>" +
		"    <address>" + fake.StreetAddress() + "</address>" +
		"    <city>" + fake.City() + "</city>" +
		"    <country>" + fake.Country() + "</country>" +
		"    <email>" + fake.EmailAddress() + "</email>" +
		"    <phone>" + fake.Phone() + "</phone>" +
		"  </shipto>" +
		"  <item>" +
		"    <title>" + fake.Title() + "</title>" +
		"    <note>" + fake.Sentences() + "</note>" +
		"    <quantity>" + fake.Digits() + "</quantity>" +
		"    <color>" + fake.Color() + "</color>" +
		"    <price>" + fake.Digits() + "." + fake.DigitsN(2) + "</price>" +
		"  </item>" +
		"  <item>" +
		"    <title>" + fake.Title() + "</title>" +
		"    <quantity>" + fake.Digits() + "</quantity>" +
		"    <price>" + fake.Digits() + "." + fake.DigitsN(2) + "</price>" +
		"  </item>" +
		"</shiporder>"
}