package excel

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// getDateList gives the list of Date structures in the term roster.
func getDateList(f *excelize.File, worksheetName string, dateRow DateRow) []Date {
	dateList := make([]Date, 0, dateRow.NumDates)
	rows, err := f.GetRows(worksheetName, excelize.Options{RawCellValue: true})
	if err != nil {
		panic(err)
	}
	for cIdx, cell := range rows[dateRow.Row] {
		serial, err := strconv.ParseFloat(cell, 64)
		if err == nil && serial > 59 {
			date, err := excelize.ExcelDateToTime(serial, false)
			if err == nil {
				dateList = append(dateList, Date{Value: date, Column: cIdx})
			}
		}

	}
	return dateList
}

// getDateRow function: Read named worksheet from given open excel file and return
// the row number with the dates PLUS total number of dates in this row.
// Row number and number of dates in that row will be fields of a custom data structure:
// `DateRow`
func getDateRow(f *excelize.File, worksheetName string) DateRow {
	dateRow := DateRow{Row: 0, NumDates: 0}
	rows, err := f.GetRows(worksheetName, excelize.Options{RawCellValue: true})
	if err == nil {
		for rIdx, row := range rows {
			datesInRow := 0
			for _, cell := range row {
				serial, err := strconv.ParseFloat(cell, 64)
				if err == nil && serial > 59 {
					_, err := excelize.ExcelDateToTime(serial, false)
					if err == nil {
						datesInRow++
					}
				}
			}
			if datesInRow > dateRow.NumDates {
				dateRow.Row = rIdx
				dateRow.NumDates = datesInRow
			}
		}
	}
	return dateRow
}

// getEmployees gives list of employees in roster along with the column number where
// employee names are found. They are found in `List` and `Column` fields of `Employees`
// structure, respectively. This is an idiomatic error handling function, with an error
// return value.
func getEmployees(f *excelize.File, roster string) (Employees, error) {
	possibleEmployeeColsToNames, err := getMapOfNameMatches(f, roster)
	if err != nil {
		return Employees{}, err
	}
	maxLen := 0
	var empCol int
	var empList []Employee
	for col, list := range possibleEmployeeColsToNames {
		if len(list) > maxLen {
			empCol = col
			empList = list
			maxLen = len(list)
		}
	}
	return Employees{empList, empCol}, nil
}

// getMapOfNameMatches function helps find out the column where employee names are found
// in the roster. It gives a mapping of column numbers to a list of strings found in that
// column that looks like people's names. Keys in this map are integers (0-based), values
// are lists of `Employee` structures. This is an error handling function.
func getMapOfNameMatches(f *excelize.File, roster string) (map[int][]Employee, error) {
	namePattern, err := regexp.Compile(`[A-Z](?:[A-Za-z]|['-][A-Z])*(?:\s+[A-Z](?:[A-Za-z]|['-][A-Z])*)+`)
	if err != nil {
		return nil, err
	}
	cols, err := f.GetCols(roster, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, err
	}
	matches := make(map[int][]Employee)
	for cIdx, col := range cols {
		for rIdx, cell := range col {
			if namePattern.MatchString(cell) {
				matches[cIdx] = append(matches[cIdx], Employee{namePattern.FindString(cell), rIdx})
			}
		}
	}
	return matches, nil
}

// getDateToColumnNumberMap gives a mapping of dates in the roster, with `time.Time`
// objects as keys, to row number of that date, as 0-based index values. This is NOT an
// error handling function.
func getDateToColumnNumberMap(f *excelize.File, worksheetName string, dateRow int) map[time.Time]int {
	dateToCol := make(map[time.Time]int)
	rows, err := f.GetRows(worksheetName, excelize.Options{RawCellValue: true})
	if err != nil {
		panic(err)
	}
	targetRow := rows[dateRow]
	for colIdx, cell := range targetRow {
		if serial, err := strconv.ParseFloat(cell, 64); err == nil {
			if date, err := excelize.ExcelDateToTime(serial, false); err == nil {
				dateToCol[date] = colIdx
			}
		}
	}
	return dateToCol
}

// clearScreen function clears screen before moving on to next set of lines to display.
func clearScreen() bool {
	var reply string
	fmt.Print("Press <enter> to move on OR anything else to Quit: ")
	n, _ := fmt.Scanf("%s", &reply)
	fmt.Fprint(os.Stdout, "\033[2J\033[H")
	return n > 0
}

// newShift gives a newly created shift that encapsulates the date of the shift (as a
// time.Time object in `Date` field) and and string value of the shift (e.g. "M", "N",
// "E", ... as `Label` field). This is an error handling function.
func newShift(f *excelize.File, sheet string, date time.Time, row int, col int) (Shift, error) {
	cell, err := excelize.CoordinatesToCellName(col+1, row+1)
	if err != nil {
		return Shift{}, err
	}
	cellType, err := f.GetCellType(sheet, cell)
	if err != nil {
		return Shift{}, err
	}
	if cellType == excelize.CellTypeUnset {
		return Shift{}, fmt.Errorf("wrong cell type")
	}
	cellValue, err := f.GetCellValue(sheet, cell, excelize.Options{RawCellValue: true})
	if err != nil {
		return Shift{}, err
	}
	if cellValue == "" || strings.TrimSpace(cellValue) == "" {
		return Shift{}, fmt.Errorf("empty entry")
	}
	return Shift{date, cellValue}, nil
}
