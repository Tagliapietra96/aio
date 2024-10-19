// num validation functions
package num

import (
	"errors"
	"strconv"
)

// Validate function checks if a string is a valid number.
// It takes a string as input and returns an error if the string is not a valid number.
func Validate(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("not a valid number")
	}
	return nil
}
