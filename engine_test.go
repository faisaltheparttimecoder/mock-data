package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var (
	reExpInteger                = "(-|)\\d+"
	reExpCharacter              = "[0-9a-zA-Z]+"
	reExpDate                   = "\\d{4}\\-(0?[1-9]|1[012])\\-(0?[1-9]|[12][0-9]|3[01])"
	reExpTimeWithoutTz          = "([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]"
	reExpTimeStamp              = reExpDate + " " + reExpTimeWithoutTz
	reExpTimeStampTz            = reExpTimeStamp + ".\\d+"
	reExpTimeStampTzWithDecimal = reExpTimeStamp + ".\\d{1,%d}"
	reExpTimeWithTz             = reExpTimeWithoutTz + ".\\d+"
	reExpIps                    = "((^\\s*((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.)" +
		"{3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\\s*$)|" +
		"(^\\s*((([0-9a-f]{1,4}:){7}([0-9a-f]{1,4}|:))|(([0-9a-f]{1,4}:){6}(:[0-9a-f]{1,4}|" +
		"((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|" +
		"(([0-9a-f]{1,4}:){5}(((:[0-9a-f]{1,4}){1,2})|" +
		":((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|" +
		"(([0-9a-f]{1,4}:){4}(((:[0-9a-f]{1,4}){1,3})|((:[0-9a-f]{1,4})?:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)" +
		"(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|" +
		"(([0-9a-f]{1,4}:){3}(((:[0-9a-f]{1,4}){1,4})|((:[0-9a-f]{1,4}){0,2}:" +
		"((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|" +
		"(([0-9a-f]{1,4}:){2}(((:[0-9a-f]{1,4}){1,5})|((:[0-9a-f]{1,4}){0,3}:" +
		"((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|" +
		"(([0-9a-f]{1,4}:){1}(((:[0-9a-f]{1,4}){1,6})|((:[0-9a-f]{1,4}){0,4}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)" +
		"(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(:(((:[0-9a-f]{1,4}){1,7})|((:[0-9a-f]{1,4}){0,5}" +
		":((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:)))(%.+)?\\s*$))"
	reExpBoolean            = "(true|false)"
	reExpText               = ".*" // TODO: In text datatype anything goes, so this regex is not going to stop or validate anything
	reExpCiText             = "[A-Z][a-z]*(\\s[A-Z][a-z]*)*"
	reExpFloat              = "[+-]?(\\d*\\.)?\\d+"
	reExpBit                = "(?:[01]+|[01]*\\\\.[01]+)"
	reExpPoint              = "^({)?((\")?[0-9]{1,},[0-9]{1,}(\")?(,)?)*(})?$"
	reExpCircle             = "^({)?((\")?<\\([0-9]{1,},[0-9]{1,}\\),[0-9]{1,}>(\")?(,)?)*(})?$"
	reExpOtherGeometricData = "^({)?((\")?[0-9]{1,},[0-9]{1,},[0-9]{1,},[0-9]{1,}(\")?(,)?)*(})?$"
	reExpUuid               = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	reExpMacAddress         = "([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})"
	reExpTsQuery            = "[a-zA-Z\\s|&!()]+"
	reExpLogSequenceNumber  = "([0-9]|[a-zA-Z0-9]){2}\\/([0-9]|[a-zA-Z0-9]){8}"
	reExpTransactionXID     = "[0-9]{1,8}:[0-9]{1,8}:"
	isItValidArray          = func(s string) bool {
		return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
	} // Check it starts and ends with a delimiter
	quoteCalculation = func(s string) bool {
		return strings.Count(s, "\"")-(strings.Count(s, "\",\"")*2)-2 == strings.Count(s, "\\\"")
	} // from the overall quote, we minus 2 quotes from the beginning and end, and then we also minus x2 quotes
	// for the transition to another json in the array i.e eg o/p ("{\"k\": \"v1\"}","{\"k\": \"v2\"}",......)
)

// Test: BuildData, all the supported datatype should send in some value
// We already have test cases under randomizer and also under build
// function below to check if the data send is valid, so here we just check
// if the data is returned and is not null
func TestBuildData(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_smallint", "smallint", false},
		{"build_smallint_array", "smallint[]", false},
		{"build_integer", "integer", false},
		{"build_integer_array", "integer[]", false},
		{"build_bigint", "bigint", false},
		{"build_bigint_array", "bigint[]", false},
		{"build_character_single", "character", false},
		{"build_character_single_array", "character[]", false},
		{"build_character_varying_single", "character varying", false},
		{"build_character_varying_single_array", "character varying[]", false},
		{"build_character_multiple", "character(10)", false},
		{"build_character_multiple_array", "character(10)[]", false},
		{"build_character_varying_multiple", "character varying(10)", false},
		{"build_character_varying_multiple_array", "character varying(10)[]", false},
		{"build_date", "date", false},
		{"build_date_array", "date[]", false},
		{"build_timestamp_with_time_zone", "timestamp with time zone", false},
		{"build_timestamp_with_time_zone_array", "timestamp with time zone[]", false},
		{"build_timestamp_without_time_zone", "timestamp without time zone", false},
		{"build_timestamp_without_time_zone_array", "timestamp without time zone[]", false},
		{"build_timestamp_without_time_zone_varying", "timestamp(4) without time zone", false},
		{"build_timestamp_without_time_zone_varying_array", "timestamp(4) without time zone[]", false},
		{"build_time_with_time_zone", "time with time zone", false},
		{"build_time_with_time_zone_array", "time with time zone[]", false},
		{"build_time_without_time_zone", "time without time zone", false},
		{"build_time_without_time_zone_array", "time without time zone[]", false},
		{"build_interval", "interval", false},
		{"build_interval_array", "interval[]", false},
		{"build_inet", "inet", false},
		{"build_inet_array", "inet[]", false},
		{"build_cidr", "cidr", false},
		{"build_cidr_array", "cidr[]", false},
		{"build_boolean", "boolean", false},
		{"build_boolean_array", "boolean[]", false},
		{"build_text", "text", false},
		{"build_text_array", "text[]", false},
		{"build_citext", "citext", false},
		{"build_citext_array", "citext[]", false},
		{"build_bytea", "bytea", false},
		{"build_money", "money", false},
		{"build_money_array", "money[]", false},
		{"build_real", "real", false},
		{"build_real_array", "real[]", false},
		{"build_double_precision", "double precision", false},
		{"build_double_precision_array", "double precision[]", false},
		{"build_numeric", "numeric(4,2)", false},
		{"build_numeric_array", "numeric(4,2)[]", false},
		{"build_bit", "bit", false},
		{"build_bit_array", "bit[]", false},
		{"build_bit_varying", "bit varying(4)", false},
		{"build_bit_varying_array", "bit varying(4)[]", false},
		{"build_uuid", "uuid", false},
		{"build_uuid_array", "uuid[]", false},
		{"build_macaddr", "macaddr", false},
		{"build_macaddr_array", "macaddr[]", false},
		{"build_json", "json", false},
		{"build_json_array", "json[]", false},
		{"build_xml", "xml", false},
		{"build_xml_array", "xml[]", false},
		{"build_tsquery", "tsquery", false},
		{"build_tsquery_array", "tsquery[]", false},
		{"build_tsvector", "tsvector", false},
		{"build_tsvector_array", "tsvector[]", false},
		{"build_pg_lsn", "pg_lsn", false},
		{"build_pg_lsn_array", "pg_lsn[]", false},
		{"build_txid_snapshot", "txid_snapshot", false},
		{"build_txid_snapshot_array", "txid_snapshot[]", false},
		{"build_path", "path", false},
		{"build_path_array", "path[]", false},
		{"build_polygon", "polygon", false},
		{"build_polygon_array", "polygon[]", false},
		{"build_line", "line", false},
		{"build_line_array", "line[]", false},
		{"build_lseg", "lseg", false},
		{"build_lseg_array", "lseg[]", false},
		{"build_box", "box", false},
		{"build_box_array", "box[]", false},
		{"build_circle", "circle", false},
		{"build_circle_array", "circle[]", false},
		{"build_point", "point", false},
		{"build_point_array", "point[]", false},
		{"build_non_existence_datatype", "unknowndatatype", true},
		{"build_non_existence_datatype_array", "unknowndatatype[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := BuildData(tt.dt); IsStringEmpty(fmt.Sprintf("%v", got)) != tt.want {
				t.Errorf("TestBuildData = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Check if the data send matches the regex
func doesDataMatchDataType(data interface{}, re string) bool {
	match := func(d string) bool {
		reExp := fmt.Sprintf("^%s$", re)
		return regexp.MustCompile(reExp).MatchString(d)
	}
	s := fmt.Sprintf("%v", data)
	stringSplit := strings.Split(s, ",")
	if len(stringSplit) > 0 {
		for _, i := range stringSplit {
			if !match(TrimPrefixNSuffix(i, "{", "}")) {
				return false
			}
		}
		return true // Matched all of them
	}
	return match(s)
}

// Test: doesDataMatchDataType, build a test to verify this function returns the appropriate data
func TestDoesDataMatchDataType(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		re   string
		want bool
	}{
		{"valid_data_that_matches_the_regex", "123456", reExpInteger, true},
		{"valid_negative_data_that_matches_the_regex", "-123456", reExpInteger, true},
		{"valid_data_array_that_matches_the_regex", "{1234,5678,90098}", reExpInteger, true},
		{"valid_negative_data_array_that_matches_the_regex", "{-1234,5678,90098}", reExpInteger, true},
		{"invalid_data_that_matches_the_regex", "abcd", reExpInteger, false},
		{"invalid_data_array_that_matches_the_regex", "{abcd,efgh,ijkl}", reExpInteger, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doesDataMatchDataType(tt.data, tt.re); got != tt.want {
				t.Errorf("TestDoesDataMatchDataType = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: buildInteger, check if the integer value is send
// for array, it should be {1, -1, -2, .....}
func TestBuildInteger(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_smallint", "smallint", true},
		{"build_smallint_array", "smallint[]", true},
		{"build_integer", "integer", true},
		{"build_integer_array", "integer[]", true},
		{"build_bigint", "bigint", true},
		{"build_bigint_array", "bigint[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildInteger(tt.dt); doesDataMatchDataType(got, reExpInteger) != tt.want {
				t.Errorf("TestBuildInteger = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildCharacter, check if the string value is send i.e it should be only a-z A-Z or 0-9
// for array it should be "{text,text,text, ....}"
func TestBuildCharacter(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_character_single", "character", true},
		{"build_character_single_array", "character[]", true},
		{"build_character_varying_single", "character varying", true},
		{"build_character_varying_single_array", "character varying[]", true},
		{"build_character_multiple", "character(10)", true},
		{"build_character_multiple_array", "character(10)[]", true},
		{"build_character_varying_multiple", "character varying(10)", true},
		{"build_character_varying_multiple_array", "character varying(10)[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildCharacter(tt.dt); doesDataMatchDataType(got, reExpCharacter) != tt.want {
				t.Errorf("TestBuildCharacter = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildDate, check if date send is in valid format
// for array it should be "{2018-05-11,2013-07-16,2017-02-05,....}"
func TestBuildDate(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_date", "date", true},
		{"build_date_array", "date[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildDate(tt.dt); doesDataMatchDataType(got, reExpDate) != tt.want {
				t.Errorf("TestBuildDate = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildTimeStamp, check if timestamp send with timezone is in valid format
// for array with timestamp without time zone it should be "{2018-05-15 09:32:30,2020-10-06 16:45:31,....}"
// for array with timestamp with time zone it should be "{2018-05-15 09:32:30.00000,2020-10-06 16:45:31.00000,....}"
// for array with timestamp without time zone varying it should be "{2011-09-12 01:58:47.4360,2014-06-11 16:35:00.1348,...}"
func TestBuildTimeStamp(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_timestamp_with_time_zone", "timestamp with time zone", true},
		{"build_timestamp_with_time_zone_array", "timestamp with time zone[]", true},
		{"build_timestamp_without_time_zone", "timestamp without time zone", true},
		{"build_timestamp_without_time_zone_array", "timestamp without time zone[]", true},
		{"build_timestamp_without_time_zone_varying", "timestamp(4) without time zone", true},
		{"build_timestamp_without_time_zone_varying_array", "timestamp(4) without time zone[]", true},
		{"build_null_timestamp_without_time_zone_varying", "timestamp() without time zone", false},
		{"build_null_timestamp_without_time_zone_varying_array", "timestamp() without time zone[]", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var re string
			if strings.HasPrefix(tt.dt, "timestamp without time zone") {
				re = reExpTimeStamp
			} else if strings.HasPrefix(tt.dt, "timestamp with time zone") {
				re = reExpTimeStampTz
			} else {
				re = fmt.Sprintf(reExpTimeStampTzWithDecimal, 4)
			}
			if got, _ := buildTimeStamp(tt.dt); doesDataMatchDataType(got, re) != tt.want {
				t.Errorf("TestBuildTimeStamp = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildTimeWithTz, check if time with timezone send is in valid format
// for array it should be {10:58:33.000000,18:10:24.000000,20:58:06.000000,....}"
func TestBuildTimeWithTz(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_time_with_time_zone", "time with time zone", true},
		{"build_time_with_time_zone_array", "time with time zone[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildTimeWithTz(tt.dt); doesDataMatchDataType(got, reExpTimeWithTz) != tt.want {
				t.Errorf("TestBuildTimeWithTz = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildInterval, check if time without timezone send is in valid format
// for array it should be {10:58:33,18:10:24,20:58:06,....}"
func TestBuildInterval(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_time_without_time_zone", "time without time zone", true},
		{"build_time_without_time_zone_array", "time without time zone[]", true},
		{"build_interval", "interval", true},
		{"build_interval_array", "interval[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildInterval(tt.dt); doesDataMatchDataType(got, reExpTimeWithoutTz) != tt.want {
				t.Errorf("TestBuildInterval = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildIps, check if the ip send is in valid format
// for array it should be {187.162.172.157,56bf:85a8:6d58:a9ee:ce47:9e50:e2a6:4a06,,....}"
func TestBuildIps(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_inet", "inet", true},
		{"build_inet_array", "inet[]", true},
		{"build_cidr", "cidr", true},
		{"build_cidr_array", "cidr[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildIps(tt.dt); doesDataMatchDataType(got, reExpIps) != tt.want {
				t.Errorf("TestBuildIps = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildBoolean, check if the boolean send is in valid format
// for array it should be {true,true,false,....}"
func TestBuildBoolean(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_boolean", "boolean", true},
		{"build_boolean_array", "boolean[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildBoolean(tt.dt); doesDataMatchDataType(got, reExpBoolean) != tt.want {
				t.Errorf("TestBuildBoolean = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildText, check if the text that is send is valid
// for array it should be "{text,text,text, ....}"
func TestBuildText(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_text", "text", true},
		{"build_text_array", "text[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildText(tt.dt); doesDataMatchDataType(got, reExpText) != tt.want {
				t.Errorf("TestBuildText = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildCiText, check if the Ci text that is send is valid
// for array it should be "{Text1,Text2,Text3, ....}"
func TestBuildCiText(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_citext", "citext", true},
		{"build_citext_array", "citext[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildCiText(tt.dt); doesDataMatchDataType(got, reExpCiText) != tt.want {
				t.Errorf("TestBuildCiText = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildBytea, check if the bytea send is valid
func TestBuildBytea(t *testing.T) {
	// We will pass this one, since TestRandomBytea() already validated it
	// and there is no array in bytes
}

// Test: buildFloat, check if the float that is send is valid
// for array it should be "{1.222,1222.222,1222.221,....}"
func TestBuildFloat(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_money", "money", true},
		{"build_money_array", "money[]", true},
		{"build_real", "real", true},
		{"build_real_array", "real[]", true},
		{"build_double_precision", "double precision", true},
		{"build_double_precision_array", "double precision[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildFloat(tt.dt); doesDataMatchDataType(got, reExpFloat) != tt.want {
				t.Errorf("TestBuildFloat = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildNumeric, check if the number that is send is valid
// for array it should be "{1.222,1222.222,1222.221,....}"
func TestBuildNumeric(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_numeric", "numeric(4,2)", true},
		{"build_numeric_array", "numeric(4,2)[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildNumeric(tt.dt); doesDataMatchDataType(got, reExpFloat) != tt.want {
				t.Errorf("TestBuildNumeric = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildBit, check if the bit that is send is valid
// for array it should be "{1010101,1000101,111001,....}"
func TestBuildBit(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_bit", "bit", true},
		{"build_bit_array", "bit[]", true},
		{"build_bit_varying", "bit varying(4)", true},
		{"build_bit_varying_array", "bit varying(4)[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildBit(tt.dt); doesDataMatchDataType(got, reExpBit) != tt.want {
				t.Errorf("TestBuildBit = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildUuid, check if the uuid that is send is valid
// for array it should be "{878b47bb-0d8a-46e2-b8bc-66b1c1e318d2,e48d09d3-a41f-494b-b698-1006f76d936b,....}"
func TestBuildUuid(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_uuid", "uuid", true},
		{"build_uuid_array", "uuid)[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildUuid(tt.dt); doesDataMatchDataType(got, reExpUuid) != tt.want {
				t.Errorf("TestBuildUuid = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildMacAddr, check if the mac address that is send is valid
// for array it should be "{56:37:65:75:4a:71,37:78:39:38:75:6d,4c:4d:73:41:74:75,34:58:43:66:37:4e,....}"
func TestBuildMacAddr(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_macaddr", "macaddr", true},
		{"build_macaddr_array", "macaddr[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildMacAddr(tt.dt); doesDataMatchDataType(got, reExpMacAddress) != tt.want {
				t.Errorf("TestBuildMacAddr = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildJson, check if the json is valid
func TestBuildJson(t *testing.T) {
	//{"build_xml", "xml", false},
	//{"build_xml_array", "xml[]", false},
	t.Run("build_json", func(t *testing.T) {
		var js json.RawMessage
		if got, _ := buildJson("json"); json.Unmarshal([]byte(fmt.Sprintf("%v", got)), &js) != nil {
			t.Errorf("TestBuildJson = %v, want valid JSON format", got)
		}
	})
	// We have tested the JSON creation above, now lets test all the quote are escaped for postgres to add them
	t.Run("build_json_array", func(t *testing.T) {
		got, _ := buildJson("json[]")
		j := fmt.Sprintf("%v", got)
		if !quoteCalculation(j) {
			t.Errorf("TestBuildJson = %v, want valid array JSON format with proper quote escape", got)
		}
	})
}

// Test: buildXml, check if the json is valid
func TestBuildXml(t *testing.T) {
	t.Run("build_xml", func(t *testing.T) {
		if got, _ := buildXml("xml"); xml.Unmarshal([]byte(fmt.Sprintf("%v", got)), new(interface{})) != nil {
			t.Errorf("TestBuildJson = %v, want valid XML format", got)
		}
	})
	// We have tested the XML creation above, now lets test all the quote are escaped for postgres to add them
	t.Run("build_xml_array", func(t *testing.T) {
		got, _ := buildXml("xml[]")
		x := fmt.Sprintf("%v", got)
		if !quoteCalculation(x) {
			t.Errorf("TestBuildJson = %v, want valid array JSON format with proper quote escape", got)
		}
	})
}

// Test: buildTsQuery, check if the text search query that is send is valid
// for array it should be "{animi | dolor,quos & hic,temporibus & tenetur  & ! et,....}"
func TestBuildTsQuery(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_tsquery", "tsquery", true},
		{"build_tsquery_array", "tsquery[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildTsQuery(tt.dt); doesDataMatchDataType(got, reExpTsQuery) != tt.want {
				t.Errorf("TestBuildMacAddr = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildTsVector, check if the text search vector that is send is valid
// for array it should be "{text,text,....}"
func TestBuildTsVector(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_tsvector", "tsvector", true},
		{"build_tsvector_array", "tsvector[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildTsVector(tt.dt); doesDataMatchDataType(got, reExpText) != tt.want { // can be some random text, so this will pass regardless
				t.Errorf("TestBuildTsVector = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildLseg, check if the log sequence number that is send is valid
// for array it should be "{74/52735449,52/68754178,....}"
func TestBuildLseg(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_pg_lsn", "pg_lsn", true},
		{"build_pg_lsn_array", "pg_lsn[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildLseg(tt.dt); doesDataMatchDataType(got, reExpLogSequenceNumber) != tt.want {
				t.Errorf("TestBuildLseg = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildTxidSnapShot, check if the transaction xid that is send is valid
// for array it should be "{18456477:79831243:,2801399:86330894:,....}"
func TestBuildTxidSnapShot(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_txid_snapshot", "txid_snapshot", true},
		{"build_txid_snapshot_array", "txid_snapshot[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := buildTxidSnapShot(tt.dt); doesDataMatchDataType(got, reExpTransactionXID) != tt.want {
				t.Errorf("TestBuildTxidSnapShot = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: buildGeometry, check if the geometric data that is send is valid
// for point array it should be "{"101,200","102,300",....}"
// for circle array it should be "{"<101,200>100","<102,300>1000",....}"
// for others array it should be "{"101,200,104,220","102,300,101,200",....}"
func TestBuildGeometry(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"build_path", "path", true},
		{"build_path_array", "path[]", true},
		{"build_polygon", "polygon", true},
		{"build_polygon_array", "polygon[]", true},
		{"build_line", "line", true},
		{"build_line_array", "line[]", true},
		{"build_lseg", "lseg", true},
		{"build_lseg_array", "lseg[]", true},
		{"build_box", "box", true},
		{"build_box_array", "box[]", true},
		{"build_circle", "circle", true},
		{"build_circle_array", "circle[]", true},
		{"build_point", "point", true},
		{"build_point_array", "point[]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var re string
			got, _ := buildGeometry(tt.dt)
			if strings.HasPrefix(tt.dt, "point") {
				re = reExpPoint
			} else if strings.HasPrefix(tt.dt, "circle") {
				re = reExpCircle
			} else {
				re = reExpOtherGeometricData
			}
			if regexp.MustCompile(re).MatchString(fmt.Sprintf("%v", got)) != tt.want {
				t.Errorf("TestBuildGeometry = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: findTimeStampDecimal, check if the function returns the correct decimal
func TestFindTimeStampDecimal(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want int
	}{
		{"valid_timestamp_with_decimal", "timestamp(4) without time zone", 4},
		{"valid_timestamp_with_decimal_array", "timestamp(4) without time zone[]", 4},
		{"invalid_timestamp_null_decimal", "timestamp() without time zone[]", 0},
		{"invalid_timestamp_null_decimal_array", "timestamp() without time zone[]", 0},
		{"invalid_timestamp_with_decimal", "timestamp(a) without time zone", 0},
		{"invalid_timestamp_with_decimal_array", "timestamp(a) without time zone[]", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findTimeStampDecimal(tt.dt); got != tt.want {
				t.Errorf("TestFindTimeStampDecimal = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: findNumberPrecision, check if the function returns the correct interger and precision
func TestFindNumberPrecision(t *testing.T) {
	tests := []struct {
		name          string
		dt            string
		wantInteger   int
		wantPrecision int
	}{
		{"valid_numeric", "numeric(5,3)", 5, 3},
		{"valid_numeric_array", "numeric(5,3)[]", 5, 3},
		{"invalid_numeric", "numeric(a,b)", 0, 0},
		{"invalid_numeric_array", "numeric(a,b)[]", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotInt, gotPrec, _ := findNumberPrecision(tt.dt); gotInt != tt.wantInteger && gotPrec != tt.wantPrecision {
				t.Errorf("TestFindNumberPrecision = (%v,%v), want %v, %v",
					gotInt, gotPrec, tt.wantInteger, tt.wantPrecision)
			}
		})
	}
}

// Test: isDataTypeAnArray, check if the function returns the correct decimal
func TestIsDataTypeAnArray(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"valid_array", "timestamp(4) without time zone[]", true},
		{"not_a_array", "timestamp(4) without time zone", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := isDataTypeAnArray(tt.dt); got != tt.want {
				t.Errorf("TestIsDataTypeAnArray = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: ArrayGenerator, we will just check here the format here is correct,
// for all the datatypes that it supported we are already doing it at randomDataByDataTypeForArray test
func TestArrayGenerator(t *testing.T) {
	tests := []struct {
		name       string
		dt         string
		originaldt string
	}{
		{"valid_int_array_generated", "int", "integer[]"},
		{"valid_timestamp_array_generated", "timestamp", "timestamp with time zone[]"},
		{"valid_character_array_generated", "string", "character varying[]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ArrayGenerator(tt.dt, tt.originaldt, 3, 5); !isItValidArray(got) {
				t.Errorf("TestIsDataTypeAnArray = %v, want valid array format", got)
			}
		})
	}
}

// Test: randomDataByDataTypeForArray, all the supported datatype should send in some value
func TestRandomDataByDataTypeForArray(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		oDt  string
		want bool
	}{
		{"random_int", "int", "", false},
		{"random_string", "string", "", false},
		{"random_date", "date", "", false},
		{"random_timestamp", "timestamp", "", false},
		{"random_timestamptz", "timestamptz", "", false},
		{"random_timestamptzWithDecimals", "string", "timestamp(6) with time zone", false},
		{"random_float", "float", "", false},
		{"random_numericFloat", "numericFloat", "", false},
		{"random_bit", "bit", "", false},
		{"random_text", "text", "", false},
		{"random_citext", "citext", "", false},
		{"random_timetz", "timetz", "", false},
		{"random_bool", "bool", "", false},
		{"random_IP", "IP", "", false},
		{"random_macaddr", "macaddr", "", false},
		{"random_uuid", "uuid", "", false},
		{"random_txid_snapshot", "txid_snapshot", "", false},
		{"random_pg_lsn", "pg_lsn", "", false},
		{"random_tsquery", "tsquery", "", false},
		{"random_tsvector", "tsvector", "", false},
		{"random_doesnt_exists", "doesnt_exists", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := randomDataByDataTypeForArray(tt.dt, tt.oDt, 4, 5); IsStringEmpty(fmt.Sprintf("%v", got)) != tt.want {
				t.Errorf("TestRandomDataByDataTypeForArray = %v, want = %v", got, tt.want)
			}
		})
	}
}

// Test: GeometricArrayGenerator, test the data send is in valid format
func TestGeometricArrayGenerator(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"validate_build_path_array", "path[]", true},
		{"validate_build_polygon_array", "polygon[]", true},
		{"validate_build_line_array", "line[]", true},
		{"validate_build_lseg_array", "lseg[]", true},
		{"validate_build_box_array", "box", true},
		{"validate_build_circle_array", "circle[]", true},
		{"validate_build_point_array", "point[]", true},
	}
	for _, tt := range tests {
		re := strings.Replace(reExpOtherGeometricData, "*", "", -1)
		format := regexp.MustCompile(re)
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasPrefix(tt.dt, "box") { // since it can have only one value
				if got := GeometricArrayGenerator(5, tt.dt); format.MatchString(got) != tt.want {
					t.Errorf("TestGeometricArrayGenerator = %v, want %v", got, tt.want)
				}
			} else {
				if got := GeometricArrayGenerator(5, tt.dt); isItValidArray(got) != tt.want {
					t.Errorf("TestGeometricArrayGenerator = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// Test: JsonXmlArrayGenerator, test the data send is in valid format
func TestJsonXmlArrayGenerator(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want bool
	}{
		{"valid_json_array", "json[]", true},
		{"invalid_json_array", "json", false},
		{"valid_xml_array", "xml[]", true},
		{"invalid_xml_array", "xml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonXmlArrayGenerator(tt.dt); !(quoteCalculation(got) && isItValidArray(got)) != tt.want {
				t.Errorf("TestJsonXmlArrayGenerator = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: buildEnumDataTypes, check if the code is able to pick random value from the emum datatype
func TestBuildEnumDataTypes(t *testing.T) {
	data := []string{"good", "ok", "bad"}
	// create the table in the database
	setDatabaseConfigForTest()
	_, err := ExecuteDB(fmt.Sprintf("CREATE TYPE rating as ENUM ('good', 'ok', 'bad');"))
	if err != nil {
		t.Errorf("TestBuildEnumDataTypes: failed with error, %v", err)
	}
	t.Run("check_for_valid_data_from_emum_type", func(t *testing.T) {
		if got, _ := buildEnumDataTypes("rating"); !StringContains(got, data) {
			t.Errorf("TestBuildEnumDataTypes = %v, want that matches one of these %v", got, data)
		}
	})
}
