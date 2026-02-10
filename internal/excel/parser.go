// Package excel: Excelize logic to read roster
package excel

import (
	"fmt"

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
