// Package excel: Excelize logic to read roster
package excel

import "github.com/xuri/excelize/v2"

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
