// tm db formatting and parsing functions
package tm

import (
	"errors"
	"time"
)

// DBParse is a helper function to parse a time string from a database.
// the function get a string and return a time.Time object.
func DBParse(s string) (time.Time, error) {
	t, err := time.Parse(dbtimeformat, s)
	if err != nil {
		return time.Time{}, errors.New("failed to parse time: " + err.Error())
	}
	return t, nil
}

// DBFormat is a helper function to format a time.Time object for a database.
// the function get a time.Time object and return a string.
func DBFormat(t time.Time) string {
	return t.Format(dbtimeformat)
}

// DBReformat is a helper function to reformat a time string from a database.
// the function get a string and return a string.
func DBReformat(s string) (string, error) {
	t, err := Parse(s)
	if err != nil {
		return "", err
	}
	return DBFormat(t), nil
}
