// Package tui executes program in terminal
package tui

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/akalankae/shift-weaver-go/internal/core"
	"github.com/akalankae/shift-weaver-go/internal/excel"
)

// Run function executes the program functionality using TUI
// Get user credentials & path for roster excel file
// Parse the roster file and build relevant data structures
// Write shifts for the user to iCloud CalDAV server
func Run() {
	// Credentials: username (email), password & roster file
	creds := getCredentials()
	fmt.Printf("Username: %s\nPassword: %s\n", creds.Username, creds.Password)

	// Roster file
	const DataDir string = "../data"
	rosterFileName, err := getRosterFile(DataDir)
	if err != nil {
		panic(err)
	}
	rosterAbsPath, err := filepath.Abs(rosterFileName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Roster: %s\n", rosterAbsPath)

	// Print list of worksheets in roster file
	sheets, err := excel.GetWorksheetList(rosterAbsPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nFound %d worksheets in %s\n", len(sheets), rosterAbsPath)

	rosterName := selectRoster(sheets)
	fmt.Println("You selected roster:", rosterName)

	// Parse the selected roster and put all the shift data together
	roster, err := excel.NewRoster(rosterAbsPath, rosterName)
	if err != nil {
		panic(err)
	}
	for emp, shifts := range roster {
		fmt.Println(emp)
		sort.Slice(shifts, func(i, j int) bool {
			return shifts[i].Date.Before(shifts[j].Date)
		})
		for _, shift := range shifts {
			fmt.Printf("%v : %s\n", shift.Date, shift.Label)
		}
	}
}

func getCredentials() (credentials core.Credentials) {
	fmt.Println("Enter your Credentials for iCloud")
	fmt.Print("Username: ")
	_, err := fmt.Scanf("%s", &credentials.Username)
	if err != nil {
		panic(err)
	}

	fmt.Print("Password: ")
	_, err = fmt.Scanf("%s", &credentials.Password)
	if err != nil {
		panic(err)
	}

	return
}

func getRosterFile(rosterFileDir string) (rosterFile string, err error) {
	files, err := filepath.Glob(filepath.Join(rosterFileDir, "*.xlsx"))
	if err != nil {
		return // "", error
	}

	for {
		var fileNumber int

		for i, file := range files {
			fmt.Println(i, filepath.Base(file))
		}
		fmt.Print("Enter index for roster file you want to read: ")
		n, _ := fmt.Scanf("%d", &fileNumber) // need fix! non digit input fucks this up
		if (fileNumber >= 0 && fileNumber < len(files)) && n == 1 {
			rosterFile = files[fileNumber]
			break
		}
		fmt.Printf("Please enter a file number between 0 and %d\n", len(files))
	}
	return // path_string, nil
}

// selectRoster function gets the user to pick one of available rosters
func selectRoster(rosters []string) string {
	fmt.Println("Select one of", len(rosters), "available rosters")

	for {
		var rosterNumber int
		for i, sheet := range rosters {
			fmt.Printf("%.2d) %s\n", i, sheet)
		}

		fmt.Print("Enter index for roster you want to read: ")
		n, _ := fmt.Scanf("%d", &rosterNumber) // need fix! non digit input fucks this up
		if (rosterNumber >= 0 && rosterNumber < len(rosters)) && n == 1 {
			return rosters[rosterNumber]
		}
		fmt.Printf("Please enter a file number between 0 and %d\n", len(rosters))
	}
}
