// Package excel: file with required data structures
package excel

import "time"

type Employee struct {
	Name string // parsed employee name
	Row  int    // row number of employee in roster
}

// Employees encapsulate employee data in a roster. Column is 0-based column number with
// employee names, List is the list of employees found in roster (as `Employee` structs)
type Employees struct {
	List   []Employee // list of employees found as `Employee` structs
	Column int        // column number with employees in roster
}

type Date struct {
	Value  time.Time // actual date in roster's date row
	Column int       // column number in roster where date is at
}

type DateRow struct {
	Row      int // row number where dates are
	NumDates int // total number of dates in date row
}

type Shift struct {
	Date  time.Time // date of the shift
	Label string    // string describing shift in roster
}

// Coordinates: row, column in roster
