package main

import "testing"

// TODO: add more tests

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
