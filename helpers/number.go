/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

// File Name: number.go
// Created by: Matteo Tagliapietra 2024-10-15
// Last Update: 2024-10-15

// This file contains helper functions to manipulate numbers.
// It contains functions to format numbers in a fancy way.

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

// helpers package contains helper functions to manipulate numbers.
package helpers

import (
	"strconv"

	"errors"

	"github.com/charmbracelet/log"
)

// imports the necessary packages
// fmt package is used to format strings

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

//
// validation functions
//

// NumberValidate function checks if a string is a valid number.
// It takes a string as input and returns an error if the string is not a valid number.
func NumberValidate(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("not a valid number")
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

//
// Number functions
//

// NumberFormat function formats a number in a fancy way.
// It takes a string as input and returns a string with the number formatted.
func NumberParse(s string) float64 {
	err := NumberValidate(s)
	if err != nil {
		log.Fatal("Trying to parse a non-number string")
	}
	n, _ := strconv.ParseFloat(s, 64)
	return n
}
