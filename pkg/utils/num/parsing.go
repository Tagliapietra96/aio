// num package provides utility functions for number manipulation.
package num

import (
	"errors"
	"strconv"
)

// NumberFormat function formats a number in a fancy way.
// It takes a string as input and returns a string with the number formatted.
func ParseFloat(s string) (float64, error) {
	err := Validate(s)
	if err != nil {
		return 0, errors.New("trying to parse a non-number string")
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("failed to parse the number: " + err.Error())
	}
	return n, nil
}
