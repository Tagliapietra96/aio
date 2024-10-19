// tm package check functions
package tm

import "time"

// SameWeek checks if two time.Time objects are in the same week.
func SameWeek(t1, t2 time.Time) bool {
	_, w1 := t1.ISOWeek()
	_, w2 := t2.ISOWeek()
	return t1.Year() == t2.Year() && w1 == w2
}

// BeforeWeek checks if a time.Time object is weeks before another time.Time object
func BeforeWeek(t1, t2 time.Time) bool {
	_, w1 := t1.ISOWeek()
	_, w2 := t2.ISOWeek()
	return t1.Year() < t2.Year() || (t1.Year() == t2.Year() && w1 < w2)
}

// AfterWeek checks if a time.Time object is weeks after another time.Time object
func AfterWeek(t1, t2 time.Time) bool {
	_, w1 := t1.ISOWeek()
	_, w2 := t2.ISOWeek()
	return t1.Year() > t2.Year() || (t1.Year() == t2.Year() && w1 > w2)
}
