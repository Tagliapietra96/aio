// str package provides utility functions for string manipulation.
package str

import (
	"strings"
	"unicode"
)

// CapitalizeFirst function capitalizes the first letter of a string.
// It takes a string as input and returns a string with the first letter capitalized.
func CapitalizeFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + strings.ToLower(s[i+1:])
	}
	return ""
}

// CapitalizeAll function capitalizes the first letter of each word in a string.
// It takes a string as input and returns a string with the first letter of each word capitalized.
func CapitalizeAll(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = CapitalizeFirst(word)
	}
	return strings.Join(words, " ")
}
