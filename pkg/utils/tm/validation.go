// tm package validation functions
package tm

import (
	"errors"
	"strings"
	"time"
)

// ValidateDate checks if a string matches the base date time format (02 Jan 2006).
func ValidateDate(s string) error {
	now := time.Now()                                // get current datetime
	tz, _ := now.Zone()                              // get timezone
	tf := strings.Replace(timeformat, "Mon ", "", 1) // remove also the day from the timeformat
	s += " 00:00"                                    // add time
	s += " " + tz
	_, err := time.Parse(tf, s)
	if err != nil {
		return errors.New("invalid time format, use the following format: 02 Jan 2006 (day month year)")
	}
	return nil
}
