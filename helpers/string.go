// helpers package contains helper functions to manipulate strings.
package helpers

import (
	"strings"
	"unicode"
)

// StringCapitalizeFirst function capitalizes the first letter of a string.
// It takes a string as input and returns a string with the first letter capitalized.
func StringCapitalizeFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + strings.ToLower(s[i+1:])
	}
	return ""
}

// StringCapitalizeAll function capitalizes the first letter of each word in a string.
// It takes a string as input and returns a string with the first letter of each word capitalized.
func StringCapitalizeAll(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = StringCapitalizeFirst(word)
	}
	return strings.Join(words, " ")
}
