package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Test: TimeNow, check if its on the correct timestamp format
func TestTimeNow(t *testing.T) {
	format := regexp.MustCompile("((19|20|21)\\d\\d)(0?[1-9]|1[012])(0?[1-9]|[12][0-9]|3[01])" +
		"[0-9]{1,2}[0-9]{1,2}[0-9]{1,2}")
	t.Run("validate_current_time_format", func(t *testing.T) {
		if got := TimeNow(); !format.MatchString(got) {
			t.Errorf("TestTimeNow = %v in format %v, want in format YYYYMMDDhhmmss", got, timestampLayout)
		}
	})
}

// Test: ConnectDB, check if are able to make successful connection
func TestConnectDB(t *testing.T) {
	t.Run("validate_connection_uri", func(t *testing.T) {
		cmdOptions.Uri = "postgres://user:pass@host:5432/db?sslmode=disable"
		db := ConnectDB()
		if db.Options().Database != "db" &&
			db.Options().Password != "pass" &&
			db.Options().User != "user" &&
			db.Options().Addr != fmt.Sprintf("%s:%d", "host", 5432) {
			t.Errorf("TestConnectDB = The connection parameters via uri are not valid")
		}
		db.Close()
		cmdOptions.Uri = "" // unset this to prevent other scripts from failing
	})
	t.Run("validate_non_connection_uri", func(t *testing.T) {
		db := ConnectDB()
		setDatabaseConfigForTest()
		if db.Options().Database != cmdOptions.Database &&
			db.Options().Password != cmdOptions.Password &&
			db.Options().User != cmdOptions.Username &&
			db.Options().Addr != fmt.Sprintf("%s:%d", cmdOptions.Hostname, cmdOptions.Port) {
			t.Errorf("TestConnectDB = The connection parameters are not correctly set")
		}
		db.Close()
	})
}

// Test: ExecuteDB, check if the database connection is able to execute any statement
func TestExecuteDB(t *testing.T) {
	t.Run("validate_database_connection_with_query", func(t *testing.T) {
		setDatabaseConfigForTest()
		if _, err := ExecuteDB("select 1"); err != nil {
			t.Errorf("TestExecuteDB expect to return select, got = %v", err)
		}
	})
}

// Test: setDBDefaults, set default environment variables
func TestSetDBDefaults(t *testing.T) {
	tests := []struct {
		name     string
		database string
		username string
		password string
		port     int
		hostname string
		want     []string
	}{
		{"empty_environment_variables", "", "", "", 0, "",
			[]string{"postgres", "postgres", "postgres", "5432", "localhost"}},
		{"set_database_environment_variables", "bigdata", "",
			"", 0, "",
			[]string{"bigdata", "postgres", "postgres", "5432", "localhost"}},
		{"set_username_environment_variables", "", "bigusername",
			"", 0, "",
			[]string{"postgres", "bigusername", "postgres", "5432", "localhost"}},
		{"set_password_environment_variables", "", "",
			"bigpassword", 0, "",
			[]string{"postgres", "postgres", "bigpassword", "5432", "localhost"}},
		{"set_port_environment_variables", "", "", "", 8999, "",
			[]string{"postgres", "postgres", "postgres", "8999", "localhost"}},
		{"set_hostname_environment_variables", "", "",
			"", 0, "bighostname",
			[]string{"postgres", "postgres", "postgres", "5432", "bighostname"}},
		{"set_all_environment_variables", "bigdb", "biguser",
			"bigpass", 8900, "bighost",
			[]string{"bigdb", "biguser", "bigpass", "8900", "bighost"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdOptions.Database = tt.database
			cmdOptions.Username = tt.username
			cmdOptions.Password = tt.password
			cmdOptions.Port = tt.port
			cmdOptions.Hostname = tt.hostname
			setDBDefaults()
			d := []string{cmdOptions.Database, cmdOptions.Username, cmdOptions.Password,
				strconv.Itoa(cmdOptions.Port), cmdOptions.Hostname}
			if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("TestSetDBDefaults() = %v, want %v", d, tt.want)
			}
		})
	}
}

// Test: Test to check if the string is empty
func TestIsStringEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", true},
		{"not_empty", "some_text", false},
		{"not_empty_spec_simbols_$", "$", false},
		{"not_empty_spec_simbols_%", "%", false},
		{"not_empty_spec_simbols_`", "`", false},
		{"not_empty_世", string('\u4e16'), false},
		{"not_empty_korean_프로그램", "프로그램", false},
		{"not_empty_hindi_कार्यक्रम", "कार्यक्रम", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStringEmpty(tt.s); got != tt.want {
				t.Errorf("IsStringEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: StartProgressBar, check if the progress bar is started with correct arguments
func TestStartProgressBar(t *testing.T) {
	t.Run("checking_the_progress_bar_status", func(t *testing.T) {
		viper.Set("MOCK_DATA_TEST_RUNNER", "false")
		// Start Progress bar
		bar := StartProgressBar("Running go progress bar unit test case", 100)
		defer bar.Close()

		// Reset
		bar.Reset()

		// Sleep for a second
		time.Sleep(1 * time.Second)

		// Increment the bar by 10
		bar.Add(10)

		// Now check the state and see if the progress bar matches our calculation
		state := bar.State()
		if state.CurrentBytes != 10.0 {
			t.Errorf("TestStartProgressBar gotBytes = %f, wantBytes (Number of bytes mismatched) = %f",
				state.CurrentBytes, 10.0)
		}
		if state.CurrentPercent != 0.1 {
			t.Errorf("TestStartProgressBar got = %f, want (Percent of bar mismatched) = %f",
				state.CurrentPercent, 0.1)
		}
		kbPerSec := fmt.Sprintf("%2.2f", state.KBsPerSecond)
		if kbPerSec != "0.01" {
			t.Errorf("TestStartProgressBar, got = %s, want (Speed mismatched) = %s",
				kbPerSec, "0.01")
		}
		viper.Set("MOCK_DATA_TEST_RUNNER", "true")
	})
}

// Test: Removal all the special characters from the string
func TestRemoveSpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"remove_spaces", "This one has too many spaces", "Thisonehastoomanyspaces"},
		{"remove_special_characters", "Thî$ ôñ hæš tøœ måñÿ ßpęčíåł čhäråçtėr", "Thhtmphrtr"},
		{"preserve_numbers", "Th15 0n3 ha$ numb3r5 1n t3xt", "Th150n3hanumb3r51nt3xt"},
		{"allow_underscores", "This-one-allow_underscores", "Thisoneallow_underscores"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveSpecialCharacters(tt.s); got != tt.want {
				t.Errorf("TestRemoveSpecialCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: Format the data for array data types
func TestFormatForArray(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		array bool
		want  string
	}{
		{"basic_string", "Basic string", false, "Basic string"},
		{"basic_string_for_array", "Basic string", true, "\"Basic string\""},
		{"data_for_circle_datatype", "<(997,970),588>", false, "<(997,970),588>"},
		{"data_for_circle_array_datatype", "<(997,970),588>", true, "\"<(997,970),588>\""},
		{"data_for_point_datatype", "436,894", false, "436,894"},
		{"data_for_point_array_datatype", "436,894", true, "\"436,894\""},
		{"data_for_other_geometrical_datatype", "72,442,172,949", false, "72,442,172,949"},
		{"data_for_other_geometrical_array_datatype", "72,442,172,949", true, "\"72,442,172,949\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatForArray(tt.s, tt.array); got != tt.want {
				t.Errorf("TestFormatForArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: YesOrNoConfirmation, check if the function on accept yes
// We ignore no and other characters, since with no it exits the program
// leaving the rest of the test cases to not to run, with other characters the
// function goes into never ending loop, so just test with "Y" for now
func TestYesOrNoConfirmation(t *testing.T) {
	t.Run("yes_confirmation_check", func(t *testing.T) {
		// Create a temp file
		content := []byte("y")
		tmpFile, err := ioutil.TempFile("", "testingYesOrNoConfirmation.tmp")
		if err != nil {
			t.Errorf("TestYesOrNoConfirmation, failed to create temp file, err:%v", err)
		}

		defer os.Remove(tmpFile.Name()) // clean up

		// Write the content to the file
		if _, err := tmpFile.Write(content); err != nil {
			t.Errorf("TestYesOrNoConfirmation, failed to write to file, err: %v", err)
		}
		if _, err := tmpFile.Seek(0, 0); err != nil {
			t.Errorf("TestYesOrNoConfirmation, failed to seek tempfile, err: %v", err)
		}

		// Pass the data as standard input
		oldStdin := os.Stdin
		defer func() {
			os.Stdin = oldStdin
		}() // Restore original Stdin

		os.Stdin = tmpFile
		if uc := YesOrNoConfirmation(); strings.ToLower(uc) != "y" {
			t.Errorf("userInput failed: %v", uc)
		}

		// Close the temp files
		if err := tmpFile.Close(); err != nil {
			t.Errorf("TestYesOrNoConfirmation, failed to close the temp file, err: %v", err)
		}
	})

}

// Test: IgnoreError, check if the function allow's us to proceed if the
// string doesn't match, we dont want to match it, since it will cause fatal error
func TestIgnoreError(t *testing.T) {
	t.Run("the_string_matches_the_error", func(t *testing.T) {
		IgnoreError("the last line of the string matches", "string matches",
			"no failure, we escaped")
	})
}

// Test: TruncateFloat, formatting the float to prevent numeric data flow
// i.e if the datatype is number(2, 4), then values like 123.4 is greater than 2
// on the left side of the "." and it should be trimmed down to something like 1.234
func TestTruncateFloat(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		m    int
		p    int
		want float64
	}{
		{"float_value_is_less_than_max", 1.9999, 4, 2, 1.9999},
		{"float_value_is_greater_than_max", 12345.9999, 4, 2, 4.0915262691133965},
		{"float_value_is_equal_to_max", 1234, 7, 2, 1234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateFloat(tt.f, tt.m, tt.p); got != tt.want {
				t.Errorf("TestTruncateFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: FloatPrecision, extract the float max and precision from the datatype
func TestFloatPrecision(t *testing.T) {
	tests := []struct {
		name          string
		dt            string
		wantMax       int
		wantPrecision int
	}{
		{"no_datatype_length", "", 5, 3},
		{"no_datatype_length_array", "[]", 5, 3},
		{"numeric_datatype", "(4,3)", 4, 3},
		{"numeric_datatype_array", "(4,3)[]", 4, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			max, precision, _ := FloatPrecision(tt.dt)
			if max != tt.wantMax {
				t.Errorf("TestTruncateFloat() = %v, want max %v", max, tt.wantMax)
			}
			if precision != tt.wantPrecision {
				t.Errorf("TestTruncateFloat() = %v, want precision %v", precision, tt.wantPrecision)
			}
		})
	}
}

// Test: Column Extractor based on regex
func TestColExtractor(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		regex string
		want  string
	}{
		{"column_extraction_alter_statement",
			"ALTER TABLE testme ADD CONSTRAINT something UNIQUE (dist_id, zipcode);",
			`ALTER TABLE(.*)ADD CONSTRAINT`, "ALTER TABLE testme ADD CONSTRAINT"},
		{"column_extraction_on_using_clause",
			"CREATE INDEX film_fulltext_idx ON public.film USING gist (fulltext);",
			`ON(.*)USING`, "ON public.film USING"},
		{"column_extraction_reference_clause",
			"ALTER TABLE ONLY public.customer ADD CONSTRAINT customer_address_id_fkey " +
				"FOREIGN KEY (address_id) REFERENCES public.address(address_id) ON UPDATE " +
				"CASCADE ON DELETE RESTRICT;", `REFERENCES[ \t]*([^\n\r]*\))`,
			"REFERENCES public.address(address_id)"},
		{"column_extraction_reference_clause_column",
			"REFERENCES public.address(address_id)" +
				"CASCADE ON DELETE RESTRICT;", `\(.*?\)`,
			"(address_id)"},
		{"column_extraction_brackets_with_brackets",
			"CREATE UNIQUE INDEX ir_translation_code_unique ON public.ir_translation USING " +
				"btree (type, lang, md5(src)) WHERE ((type)::text = 'code'::text);", `\(([^\[\]]*)\)`,
			"(type, lang, md5(src)) WHERE ((type)::text = 'code'::text)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ColExtractor(tt.s, tt.regex); got != tt.want {
				t.Errorf("TestColExtractor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: TrimPrefixNSuffix, trim the prefix and suffix of a string
func TestTrimPrefixNSuffix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		suffix string
		prefix string
		want   string
	}{
		{"trim_none", "good text", "", "", "good text"},
		{"trim_alter_statement", "ALTER TABLE testme ADD CONSTRAINT", "ALTER TABLE",
			"ADD CONSTRAINT", " testme "},
		{"trim_brackets", "(inside bracket)", "(", ")", "inside bracket"},
		{"trim_string_on_using", "ON (inside bracket) USING", "ON", "USING",
			" (inside bracket) "},
		{"trim_none_on_using", "CREATE INDEX film_fulltext_idx ON public.film USING gist (fulltext);",
			"ON", "USING", "CREATE INDEX film_fulltext_idx ON public.film USING gist (fulltext);"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimPrefixNSuffix(tt.s, tt.suffix, tt.prefix); got != tt.want {
				t.Errorf("TestTrimPrefixNSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: RemoveEverySuffixAfterADelimiter, remove everything after a delimiter
func TestRemoveEverySuffixAfterADelimiter(t *testing.T) {
	tests := []struct {
		name string
		s    string
		d    string
		want string
	}{
		{"nothing_to_remove", "this is a text, that has nothing to remove", "exists",
			"this is a text, that has nothing to remove"},
		{"delimiter_exists", "this is a text, that has nothing to remove", "text", "this is a "},
		{"capital_delimiter_exists", "this is a text, that has nothing to remove", "TEXT",
			"this is a "},
		{"where_clause_exists", "btree (type, lang, md5(src)) WHERE ((type)::text = 'code'::text);",
			"where", "btree (type, lang, md5(src)) "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveEverySuffixAfterADelimiter(tt.s, tt.d); got != tt.want {
				t.Errorf("TestRemoveEverySuffixAfterADelimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TEST: BracketsExists, check for existence of bracket
func TestBracketsExists(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", false},
		{"not_empty_no_brackets", "some_text", false},
		{"not_empty_with_brackets", "(some_text)", true},
		{"not_empty_with_multiple_brackets", "(some_text(inner brackets))", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BracketsExists(tt.s); got != tt.want {
				t.Errorf("TestBracketsExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: IsSubStringAvailableOnString, check for the existence of a string on a text
func TestIsSubStringAvailableOnString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		c    string
		want bool
	}{
		{"no_match", "some text", "nothing", false},
		{"text_exists_without_regex", "some text", "text", true},
		{"text_exists_with_regex", "some text on this sentence", ".text*", true},
		{"multiple_condition_checker",
			"CREATE UNIQUE INDEX film_fulltext_idx ON public.film USING gist (fulltext);",
			"ADD CONSTRAINT.*PRIMARY KEY|ADD CONSTRAINT.*UNIQUE|CREATE UNIQUE INDEX", true},
		{"multiple_condition_checker_case_insensitive",
			"create unique index film_fulltext_idx ON public.film USING gist (fulltext);",
			"ADD CONSTRAINT.*PRIMARY KEY|ADD CONSTRAINT.*UNIQUE|CREATE UNIQUE INDEX", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSubStringAvailableOnString(tt.s, tt.c); got != tt.want {
				t.Errorf("TestIsSubStringAvailableOnString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: StringContains, check if the given string is available on the array
func TestStringContains(t *testing.T) {
	tests := []struct {
		name string
		s    string
		a    []string
		want bool
	}{
		{"string_exists_in_array", "a", []string{"a", "b", "c", "d"}, true},
		{"string_exists_in_array_in_middle", "c", []string{"a", "b", "c", "d"}, true},
		{"string_doesnt_exists_in_array", "a", []string{"b", "c", "d"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringContains(tt.s, tt.a); got != tt.want {
				t.Errorf("TestStringContains = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: StringHasPrefix, check if the given string has a match with elements in the array
// i.e does the element in the array has starting characters of the string
func TestStringHasPrefix(t *testing.T) {
	tests := []struct {
		name string
		s    string
		a    []string
		want bool
	}{
		{"string_matches_with_ele_in_array", "topple", []string{"true", "lame", "top", "tame"}, true},
		{"string_doesnt_match_with_ele_in_array", "top", []string{"trans", "driver", "tomorrow"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringHasPrefix(tt.s, tt.a); got != tt.want {
				t.Errorf("TestStringHasPrefix = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test: CharLen, extract the length from the given string
func TestCharLen(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"no_value", "char", 1},
		{"character_length_available", "char(7)", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CharLen(tt.s); got != tt.want || err != nil {
				t.Errorf("TestCharLen = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}
