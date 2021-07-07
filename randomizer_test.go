package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"net"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Date unit test cases
type dateTest []struct {
	name     string
	fromYear int
	toYear   int
	fromDate time.Time
	toDate   time.Time
	decimal  int
}

var (
	dateByYear = func(n int) time.Time {
		return time.Now().AddDate(n, 0, 0)
	}
	dateTests = dateTest{
		{"date_from_ten_years_ago_to_now", -10, 0, dateByYear(-10),
			dateByYear(0), 2},
		{"date_from_now_to_future", 0, 10, dateByYear(0),
			dateByYear(10), 4},
		{"date_between_past_and_future", -10, 10, dateByYear(-10),
			dateByYear(10), 6},
	}
)

// Test: RandomString, should generate a random string of the requested size
func TestRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero_length_string", 0},
		{"one_random_character_length_string", 1},
		{"ten_random_character_length_string", 10},
		{"hundred_random_character_length_string", 100},
		{"thousand_random_character_length_string", 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(RandomString(tt.length)); got != tt.length && reflect.TypeOf(got).Kind() != reflect.String {
				t.Errorf("TestRandomString() = %v, want string and length = %v", got, tt.length)
			}
		})
	}
}

// Test: RandomInt, should generate a random integer between a max and minimum value
func TestRandomInt(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{"single_digit_integer", 1, 10},
		{"double_digit_integer", 100, 200},
		{"minimum_greater_than_max", 300, 200},
		{"negative_min_value", -100, 200},
		{"negative_max_value", 100, -200},
		{"both_negative_value", -100, -200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomInt(tt.min, tt.max); got < tt.min && reflect.TypeOf(got).Kind() != reflect.Int {
				t.Errorf("TestRandomInt() = %v, want int and >= %v", got, tt.min)
			}
		})
	}
}

// Test: RandomBytea, should generate a random bytea based on the size provided
func TestRandomBytea(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"single_length_bytea", 1},
		{"ten_length_bytea", 10},
		{"hundred_length_bytea", 100},
		{"thousand_length_bytea", 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(RandomBytea(tt.length)); got > tt.length {
				t.Errorf("TestRandomBytea() = %v, want length <= %v", got, tt.length)
			}
		})
	}
}

// Test: RandomFloat, should generate a random float number based on the min, max, precision provided
func TestRandomFloat(t *testing.T) {
	tests := []struct {
		name      string
		min       int
		max       int
		precision int
	}{
		{"single_digit_float_value", 1, 10, 2},
		{"double_digit_float_value", 10, 100, 2},
		{"triple_digit_float_value", 100, 1000, 2},
		{"bigger_precision_float_value", 100, 1000, 24},
		{"negative_float_value", -100, 1, 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomFloat(tt.min, tt.max, tt.precision); got > float64(tt.max) &&
				reflect.TypeOf(got).Kind() != reflect.Float64 {
				t.Errorf("TestRandomFloat() = %v, want float64 and length <= %v", got, tt.max)
			}
		})
	}
}

// Test: RandomCalenderDateTime, get a unix time in seconds between two years
func TestRandomCalenderDateTime(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomCalenderDateTime(tt.fromYear, tt.toYear); !got.After(tt.fromDate) &&
				!got.Before(tt.toDate) {
				t.Errorf("TestRandomCalenderDateTime = %v, want dates >= %v < %v", got, tt.fromDate, tt.toDate)
			}
		})
	}
}

// Test: RandomDate, check if the dates provided is in the given format,
// RandomDate using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomDate(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomDate(tt.fromYear, tt.toYear); !doesDataMatchDataType(got, reExpDate) {
				t.Errorf("TestRandomDate = %v, want in format = YYYY-MM-DD", got)
			}
		})
	}
}

// Test: RandomTimestamp, check if the timestamp provided is in the given format,
// RandomTimestamp using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomTimestamp(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomTimestamp(tt.fromYear, tt.toYear); !doesDataMatchDataType(got, reExpTimeStamp) {
				t.Errorf("TestRandomTimestamp = %v, want in format = YYYY-MM-DD hh:mm:ss", got)
			}
		})
	}
}

// Test: RandomTimeStampTz, check if the timestamp with tz provided is in the given format,
// RandomTimeStampTz using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomTimeStampTz(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomTimeStampTz(tt.fromYear, tt.toYear); !doesDataMatchDataType(got, reExpTimeStampTz) {
				t.Errorf("TestRandomTimeStampTz = %v, want in format = YYYY-MM-DD hh:mm:ss.000000", got)
			}
		})
	}
}

// Test: RandomTimeStampTzWithDecimals, check if the timestamp with tz provided is in the given format,
// RandomTimeStampTzWithDecimals using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomTimeStampTzWithDecimals(t *testing.T) {
	for _, tt := range dateTests {
		re := fmt.Sprintf(reExpTimeStampTzWithDecimal, tt.decimal)
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomTimeStampTzWithDecimals(tt.fromYear, tt.toYear, tt.decimal); !doesDataMatchDataType(got, re) {
				t.Errorf("TestRandomTimeStampTz = %v, want in format = YYYY-MM-DD hh:mm:ss.%s", got,
					strings.Repeat("0", tt.decimal))
			}
		})
	}
}

// Test: RandomTime, check if the time without tz provided is in the given format,
// RandomTime using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomTime(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomTime(tt.fromYear, tt.toYear); !doesDataMatchDataType(got, reExpTimeWithoutTz) {
				t.Errorf("TestRandomTime = %v, want in format = hh:mm:ss", got)
			}
		})
	}
}

// Test: RandomTimeTz, check if the time without tz provided is in the given format,
// RandomTimeTz using the RandomCalenderDateTime, so we have already checked from above that they are
// in range, so lets check the format here
func TestRandomTimeTz(t *testing.T) {
	for _, tt := range dateTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := RandomTimeTz(tt.fromYear, tt.toYear); !doesDataMatchDataType(got, reExpTimeWithTz) {
				t.Errorf("TestRandomTimeTz = %v, want in format = hh:mm:ss.000000", got)
			}
		})
	}
}

// Test: RandomBoolean, check if it return only true or false
func TestRandomBoolean(t *testing.T) {
	t.Run("boolean_value", func(t *testing.T) {
		if got := RandomBoolean(); reflect.TypeOf(got).Kind() != reflect.Bool {
			t.Errorf("TestRandomBoolean = %v, want boolean (true or false)", got)
		}
	})
}

// Test: RandomParagraphs, check if the value returned is paragraph
func TestRandomParagraphs(t *testing.T) {
	t.Run("random_paragraph", func(t *testing.T) {
		if got := RandomParagraphs(); len(got) <= 0 {
			t.Errorf("TestRandomParagraphs = %v, want > 0", len(got))
		}
	})
}

// Test: RandomCiText, check if the words emptied is capitalized
func TestRandomCiText(t *testing.T) {
	t.Run("check_capitalize_words", func(t *testing.T) {
		if got := RandomCiText(); !doesDataMatchDataType(got, reExpCiText) {
			t.Errorf("TestRandomCiText = %v, want = %v", got, strings.Title(got))
		}
	})
}

// Test: RandomIP, check if the ip address is valid
func TestRandomIP(t *testing.T) {
	t.Run("validate_ip", func(t *testing.T) {
		if got := net.ParseIP(RandomIP()); got == nil {
			t.Errorf("TestRandomIP = %v, want = %v", got, "valid ip address")
		}
	})
}

// Test: RandomBit, check if the bit produced is valid and is of the length requested
func TestRandomBit(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"one_size_bit", 1},
		{"ten_size_bit", 10},
		{"fifty_size_bit", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomBit(tt.length);
			!doesDataMatchDataType(got, reExpBit) || !(len(got) > 0 && len(got) <= tt.length) {
				t.Errorf("TestRandomBit = %v, want valid bit format or length <= %v", got, tt.length)
			}
		})
	}
}

// Test: RandomUUID, check the uuid provided is valid
func TestRandomUUID(t *testing.T) {
	t.Run("validate_uuid", func(t *testing.T) {
		if got, err := uuid.Parse(RandomUUID()); err != nil {
			t.Errorf("TestRandomUUID = %v, want valid UUID err = %v", got, err)
		}
	})
}

// Test: RandomMacAddress, check if its a valid mac address
func TestRandomMacAddress(t *testing.T) {
	t.Run("validate_mac_address", func(t *testing.T) {
		if got, err := net.ParseMAC(RandomMacAddress()); err != nil {
			t.Errorf("TestRandomMacAddress = %v, want valid mac address, err %v", got.String(), err)
		}
	})
}

// Test: RandomTSQuery, check for valid text search query, there is no pattern here
// we randomly generate some text search, so we just check here if they are a text that's all
func TestRandomTSQuery(t *testing.T) {
	t.Run("validate_text_search", func(t *testing.T) {
		if got := RandomTSQuery(); reflect.TypeOf(got).Kind() != reflect.String {
			t.Errorf("TestRandomTSQuery = %v, want valid text", got)
		}
	})
}

// Test: RandomTSVector, check if the text search vector is valid, again its some random text
// we can just check if its a valid text
func TestRandomTSVector(t *testing.T) {
	t.Run("validate_text_search", func(t *testing.T) {
		if got := RandomTSVector(); reflect.TypeOf(got).Kind() != reflect.String {
			t.Errorf("TestRandomTSVector = %v, want valid text", got)
		}
	})
}

// Test: RandomGeometricData, check if the data send is valid geometric data
func TestRandomGeometricData(t *testing.T) {
	tests := []struct {
		name    string
		g       string
		isArray bool
		re      string
	}{
		{"validate_point_non_array", "point", false, reExpPoint},
		{"validate_point_array", "point", true, reExpPoint},
		{"validate_circle_non_array", "circle", false, reExpCircle},
		{"validate_circle_array", "circle", true, reExpCircle},
		{"validate_other_geometric_data_non_array", "others", false, reExpOtherGeometricData},
		{"validate_other_geometric_data_array", "others", true, reExpOtherGeometricData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := regexp.MustCompile(tt.re)
			if got := RandomGeometricData(10, tt.g, tt.isArray); !format.MatchString(got) {
				t.Errorf("TestRandomGeometricData = %v, want valid %s data", got, tt.g)
			}
		})
	}
}

// Test: RandomLSN, check is the send LSN is in valid format
func TestRandomLSN(t *testing.T) {
	t.Run("validate_lsn", func(t *testing.T) {
		if got := RandomLSN(); !doesDataMatchDataType(got, reExpLogSequenceNumber) {
			t.Errorf("TestRandomLSN = %v, want valid LSN data in the format eg 43/4e584273", got)
		}
	})
}

// Test: RandomTXID, check is the send TXID is in valid format
func TestRandomTXID(t *testing.T) {
	t.Run("validate_txid", func(t *testing.T) {
		if got := RandomTXID(); !doesDataMatchDataType(got, reExpTransactionXID) {
			t.Errorf("TestRandomTXID = %v, want valid LSN data in the format eg 77552990:99910540:", got)
		}
	})
}

// Test: RandomJSON, check if the json send is valid
func TestRandomJSON(t *testing.T) {
	t.Run("valid_json_non_array", func(t *testing.T) {
		var js json.RawMessage
		if got := RandomJSON(false); json.Unmarshal([]byte(got), &js) != nil {
			t.Errorf("TestRandomJSON = %v, want valid JSON format", got)
		}
	})
	// We have tested the JSON creation above, now lets test all the quote are escaped for postgres to add them
	t.Run("valid_json_array", func(t *testing.T) {
		got := RandomJSON(true)
		if strings.Count(got, "\"") != strings.Count(got, "\\\"") {
			t.Errorf("TestRandomJSON = %v, want valid array JSON format with proper quote escape", got)
		}
	})
}

// Test: RandomXML, check if the xml send is valid
func TestRandomXML(t *testing.T) {
	t.Run("valid_xml_non_array", func(t *testing.T) {
		if got := RandomXML(false); xml.Unmarshal([]byte(got), new(interface{})) != nil {
			t.Errorf("TestRandomXML = %v, want valid XML format", got)
		}
	})
	t.Run("valid_xml_array", func(t *testing.T) {
		got := RandomXML(true)
		if strings.Count(got, "\"") != strings.Count(got, "\\\"") {
			t.Errorf("TestRandomXML = %v, want valid XML format with proper quote escape", got)
		}
	})
}

// Test: RandomValueFromLength, the value should be with the length
func TestRandomValueFromLength(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero_length", 0},
		{"one_length", 1},
		{"five_length", 5},
		{"ten_length", 10},
		{"fifty_length", 50},
	}
	for _, tt := range tests {
		if got := RandomValueFromLength(tt.length); got < 0 && got >= tt.length {
			t.Errorf("TestRandomValueFromLength = %v, want value within the range of 0 to %v", got, tt.length)
		}
	}
}

// Test: RandomPickerFromArray, check if the picked value is in the array of elements
func TestRandomPickerFromArray(t *testing.T) {
	tests := []struct {
		name string
		ele  []string
	}{
		{"null_array_of_elements", []string{}},
		{"two_array_of_elements", []string{"a", "b"}},
		{"five_array_of_elements", []string{"a", "b", "c", "d", "e"}},
		{"ten_array_of_elements", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomPickerFromArray(tt.ele); got != "" && !StringContains(got, tt.ele) {
				t.Errorf("TestRandomPickerFromArray = %v, want element in the array = %v", got, tt.ele)
			}
		})
	}
}
