package string_utils

import (
	"fmt"
	"testing"
)

func TestFirstLetterUppercase(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{"hello world", "Hello world"},
		{"GO", "GO"},
		{"capital", "Capital"},
		{"go", "Go"},
		{"", ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input: %s, expected: %s", tc.input, tc.expected), func(t *testing.T) {
			if result := FirstLetterUppercase(tc.input); result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
		})
	}
}
