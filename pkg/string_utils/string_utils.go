package string_utils

import (
	"strings"
)

// FirstLetterUppercase makes the first letter of string uppercase
func FirstLetterUppercase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string([]rune(s)[0])) + string([]rune(s)[1:])
}
