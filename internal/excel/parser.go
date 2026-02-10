// Package excel: Excelize logic to read roster
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

// NewRoster constructs a Roster, that maps employees in the roster (as Employee type) to
// their shifts (as Shift type) for that time period.
func NewRoster(filePath, rosterName string) (roster map[string][]Shift, err error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return roster, err
	}
	defer file.Close()
	dates := getDateRow(file, rosterName)
	if dates.NumDates == 0 {
		return roster, fmt.Errorf("could not find date row in %q", rosterName)
	}
	fmt.Println(rosterName, "  Date Row:", dates.Row, ", Number of dates:", dates.NumDates)
	emps, err := getEmployees(file, rosterName)
	if err != nil {
		return roster, err
	}

	dateToColumnNumberMap := getDateToColumnNumberMap(file, rosterName, dates.Row)
	roster = make(map[string][]Shift, len(emps.List))
	for _, emp := range emps.List {
		for date, dateCol := range dateToColumnNumberMap {
			shift, err := newShift(file, rosterName, date, emp.Row, dateCol)
			if err == nil {
				roster[emp.Name] = append(roster[emp.Name], shift)
			}
		}
	}
	return
}

func GetWorksheetList(excelFile string) (worksheetList []string, err error) {
	file, err := excelize.OpenFile(excelFile)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()
	worksheetList = file.GetSheetList()
	return
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

// getDateList gets the list of Date structures in the roster for the term.
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

// Prompt to clear screen before moving on to next sheet
func clearScreen() bool {
	var reply string
	fmt.Print("Press <enter> to move on OR anything else to Quit: ")
	n, _ := fmt.Scanf("%s", &reply)
	fmt.Fprint(os.Stdout, "\033[2J\033[H")
	return n > 0
}
