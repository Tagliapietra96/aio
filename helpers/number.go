// helpers package contains helper functions to manipulate numbers.
package helpers

import (
	"strconv"

	"errors"

	"github.com/charmbracelet/log"
)

// NumberValidate function checks if a string is a valid number.
// It takes a string as input and returns an error if the string is not a valid number.
func NumberValidate(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("not a valid number")
	}
	return nil
}

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
