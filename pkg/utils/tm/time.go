// tm package is used to handle time formatting and parsing.
package tm

import (
	"aio/pkg/utils/str"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

// timeformat is the default time format used to format time.Time objects.
const timeformat = "Mon 02 Jan 2006 15:04 MST"

// dbtimeformat is the default time format used to format time.Time objects for a database.
const dbtimeformat = "2006-01-02 15:04:05"

const validformats = `
	Please provide a valid time string in the following format:
	- "now", "today", "yesterday", "tomorrow"
	- "day", "days", "week", "weeks", "month", "months", "year", "years"
	- "mon", "tue", "wed", "thu", "fri", "sat", "sun"
	- "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"
	- "mondays", "tuesdays", "wednesdays", "thursdays", "fridays", "saturdays", "sundays"
	- "last", "next", "this" + "day", "week", "month", "year" or the day of the week
	- "in" + number + "day", "week", "month", "year" or the day of the week
	- number + "day", "week", "month", "year" or the day of the week + "ago"

	You can also provide a time in the following format:
	- date + "HH:MM"
	- date + "at HH:MM"
	- date + " @ HH:MM"
	- date + " on HH:MM"

	Or you can provide a specific date with the following format:
	- "Mon 02 Jan 2024 15:04"
	- "02 Jan 2024 15:04"
	- "02 Jan 2024"

	Please note that the time string is case insensitive
	If you don't provide a time, it will default to 00:00

	Example:
	- "Next Monday"
	- "In 2 weeks"
	- "Last month"
	- "3 days ago"
	- "in 2 mon 15:00"
	- "3 Thursdays ago at 12:00"
	- "Yesterday at 12:00"
	- "02 Jan 2024 15:04"
	- "02 Jan 2024"
`

// days is a slice of strings containing the days of the week.
var days = []string{
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
	"Sun",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
	"Mondays",
	"Tuesdays",
	"Wednesdays",
	"Thursdays",
	"Fridays",
	"Saturdays",
	"Sundays",
}

func getErr(err error) error {
	e := errors.New("Valid time foramts: " + validformats)
	return errors.Join(err, e)
}

// Format returns a formatted time string from a time.Time object.
// the output will be in the following format: "Mon 02 Jan 06 15:04".
func Format(t time.Time) string {
	ts := t.Format(timeformat)
	ta := strings.Split(ts, " ") // remove timezone
	ta = ta[:len(ta)-1]
	return strings.Join(ta, " ")
}

// ParseAt is a helper function to parse a time string with a specific time.
// the function get a time.Time object and a time string with a specific time (HH:MM).
// the function will return a time.Time object with the specific time.
func ParseAt(t time.Time, at string) (time.Time, error) {
	var err error
	at = strings.ToUpper(at)
	parser := "2006-01-02 MST 15:04"
	ts := t.Format("2006-01-02 MST") + " " + at
	t, err = time.Parse(parser, ts)
	if err != nil {
		return time.Time{}, getErr(errors.New("failed to parse time: " + err.Error()))
	}

	return t, nil
}

// parseWeekDay is a helper function to parse a time string with a specific day of the week.
// the function use 3 parameters, the time string (day of the week with last, next or this directives), the specific time and a multiplier.
// the function will return a time.Time object with the specific day of the week and time.
// the time.Time output will be calculated based on the directives received, starting from the current datetime.
func parseWeekDay(s, at string, multiplier int) (time.Time, error) {
	now := time.Now() // get current datetime

	// catch the directive and parse the time string accordingly
	// if the directive is "Last" we will subtract the time
	// if the directive is "Next" we will add the time
	// if the directive is "This" or default we will keep the current week
	switch {
	case strings.Contains(s, "Last"):
		s = strings.Replace(s, "Last ", "", 1)     // remove the directive
		t, err := parseWeekDay(s, at, multiplier)  // parse the day of the week
		return t.AddDate(0, 0, -7*multiplier), err // subtract the time based on the multiplier
	case strings.Contains(s, "Next"):
		s = strings.Replace(s, "Next ", "", 1)    // remove the directive
		t, err := parseWeekDay(s, at, multiplier) // parse the day of the week
		return t.AddDate(0, 0, 7*multiplier), err // add the time based on the multiplier
	default:
		t := now
		di := slices.Index(days, s)
		wds := days[di] // get the day of the week
		var wd time.Weekday

		// set the day of the week for the filter
		switch wds {
		case "Mon", "Monday", "Mondays":
			wd = time.Monday
		case "Tue", "Tuesday", "Tuesdays":
			wd = time.Tuesday
		case "Wed", "Wednesday", "Wednesdays":
			wd = time.Wednesday
		case "Thu", "Thursday", "Thursdays":
			wd = time.Thursday
		case "Fri", "Friday", "Fridays":
			wd = time.Friday
		case "Sat", "Saturday", "Saturdays":
			wd = time.Saturday
		case "Sun", "Sunday", "Sundays":
			wd = time.Sunday
		default:
			err := errors.New("invalid day of the week")
			return time.Time{}, getErr(err)
		}

		// continue to subtract the day of the week from the current datetime until we reach the correct day
		for t.Weekday() != wd {
			t = t.AddDate(0, 0, -1)
		}

		// check if the day of the week is in the past or in the future
		// if the datetime isn't in the current week, we will add or subtract a week to get the correct datetime
		switch {
		case AfterWeek(t, now):
			t = t.AddDate(0, 0, -7)
		case BeforeWeek(t, now):
			t = t.AddDate(0, 0, 7)
		}

		// return the datetime with the specific time
		return ParseAt(t, at)
	}
}

// Parse is a helper function to parse a time string.
// the function get a string and return a time.Time object.
// the function will parse the time string and return a time.Time object based on the directives received.
// the time.Time output will be calculated based on the directives received, starting from the current datetime.
// in case of an error, the function will return an error.
func Parse(s string) (time.Time, error) {
	now := time.Now() // get current datetime
	at := "00:00"     // set default time

	s = str.CapitalizeAll(s) // capitalize all words

	// remove the day from the timeformat if the format is like "Mon 02 Jan 24 15:04"
	ss := strings.Split(s, " ")
	if len(ss) > 4 && slices.Index(days, ss[0]) != -1 {
		s = strings.Join(ss[1:], " ")
	}

	// check if the datetime string contains a time and remove it, then set the time
	switch {
	case strings.Contains(s, " At "):
		sa := strings.Split(s, " At ")
		s = sa[0]
		at = strings.ToUpper(sa[1])
	case strings.Contains(s, " @ "):
		sa := strings.Split(s, " @ ")
		s = sa[0]
		at = strings.ToUpper(sa[1])
	case strings.Contains(s, " On "):
		sa := strings.Split(s, " On ")
		s = sa[0]
		at = strings.ToUpper(sa[1])
	case strings.Contains(s, ":"):
		sa := strings.Split(s, " ")
		for _, v := range sa {
			if strings.Contains(v, ":") {
				at = strings.ToUpper(v)
				s = strings.Replace(s, " "+v, "", 1)
				break
			}
		}
	}

	// normalize the time string
	switch {
	case strings.Contains(s, "Days"):
		s = strings.Replace(s, "Days", "Day", 1)
	case strings.Contains(s, "Weeks"):
		s = strings.Replace(s, "Weeks", "Week", 1)
	case strings.Contains(s, "Months"):
		s = strings.Replace(s, "Months", "Month", 1)
	case strings.Contains(s, "Years"):
		s = strings.Replace(s, "Years", "Year", 1)
	}

	// if the time string contains " Ago" or " In ", parse it accordingly, we will use a multiplier to multiply the output based on the directive received
	multiplier := 1
	// if the time string contains " Ago" we will parse it as "Last", so we know it's in the past
	if strings.Contains(s, " Ago") {
		sa := strings.Split(s, " ")
		if len(sa) < 3 {
			return time.Time{}, getErr(errors.New("failed to parse time string, invalid format"))
		}

		num, err := strconv.Atoi(sa[slices.Index(sa, "Ago")-2])
		if err != nil {
			return time.Time{}, getErr(err)
		}
		s = "Last " + sa[slices.Index(sa, "Ago")-1]
		multiplier = num
	}
	// if the time string contains " In " we will parse it as "Next", so we know it's in the future
	if strings.Contains(s, "In ") {
		sa := strings.Split(s, " ")
		if len(sa) < 3 {
			return time.Time{}, getErr(errors.New("failed to parse time string, invalid format"))
		}

		num, err := strconv.Atoi(sa[slices.Index(sa, "In")+1])
		if err != nil {
			return time.Time{}, getErr(err)
		}
		s = "Next " + sa[slices.Index(sa, "In")+2]
		multiplier = num
	}

	// check if the datetime string contains a day of the week
	ss = strings.Split(s, " ")
	di := -1
	for _, v := range ss {
		if slices.Index(days, v) != -1 {
			di = slices.Index(days, v)
			break
		}
	}

	// if the datetime string contains a day of the week, parse it accordingly
	if di != -1 {
		return parseWeekDay(s, at, multiplier)
	}

	// set an operator based on the directive received, the operator will be used to multiply the output based on the directive received
	// if the directive is "Last" we will multiply the output by -1 to subtract the time
	// if the directive is "Next" we will multiply the output by 1 to add the time
	// if the directive is "This", or default we will multiply the output by 0 to keep the time
	var op func(int) int
	switch {
	case strings.Contains(s, "Last"):
		op = func(i int) int { return i * -1 }
		s = strings.Replace(s, "Last ", "", 1)
	case strings.Contains(s, "Next"):
		op = func(i int) int { return i }
		s = strings.Replace(s, "Next ", "", 1)
	case strings.Contains(s, "This"):
		op = func(i int) int { return 0 * i }
		s = strings.Replace(s, "This ", "", 1)
	default:
		op = func(i int) int { return 0 * i }
	}

	// search for the directive in the time string and parse it accordingly
	switch s {
	case "Now":
		return now, nil // if the directive is "Now" we will return the current datetime
	case "Today":
		return ParseAt(now, at) // if the directive is "Today" we will return the current datetime with the specific time
	case "Yesterday":
		return ParseAt(now.AddDate(0, 0, -1), at) // if the directive is "Yesterday" we will return the current datetime minus 1 day with the specific time
	case "Tomorrow":
		return ParseAt(now.AddDate(0, 0, 1), at) // if the directive is "Tomorrow" we will return the current datetime plus 1 day with the specific time
	case "Day", "Days":
		return ParseAt(now.AddDate(0, 0, op(1*multiplier)), at) // if the directive is "Day" or "Days" we will return the current datetime plus or minus the number of days with the specific time
	case "Week", "Weeks":
		return ParseAt(now.AddDate(0, 0, op(7*multiplier)), at) // if the directive is "Week" or "Weeks" we will return the current datetime plus or minus the number of weeks with the specific time
	case "Month", "Months":
		return ParseAt(now.AddDate(0, op(1*multiplier), 0), at) // if the directive is "Month" or "Months" we will return the current datetime plus or minus the number of months with the specific time
	case "Year", "Years":
		return ParseAt(now.AddDate(op(1*multiplier), 0, 0), at) // if the directive is "Year" or "Years" we will return the current datetime plus or minus the number of years with the specific time
	default:
		// default we assume the directive is a specific date and parse it accordingly
		tz, _ := now.Zone()                              // get timezone
		tf := strings.Replace(timeformat, "Mon ", "", 1) // remove also the day from the timeformat
		s += " " + at                                    // add time
		s += " " + tz                                    // add timezone

		t, err := time.Parse(tf, s) // if the time string is invalid, we will return an error
		if err != nil {
			return time.Time{}, getErr(errors.New("failed to parse time: " + err.Error()))
		}
		return t, nil
	}
}
