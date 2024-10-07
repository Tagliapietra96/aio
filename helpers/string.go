////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// File Name: string.go
// Created by: Matteo Tagliapietra 2024-09-01
// Last Update: 2024-10-05

// This file contains helper functions to manipulate strings.

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// helpers package contains helper functions to manipulate strings.
package helpers

// imports the necessary packages
// strings package is used to manipulate strings
// unicode package is used to manipulate unicode characters
import (
	"strings"
	"unicode"
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

//
// Capitalization functions
//

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

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
