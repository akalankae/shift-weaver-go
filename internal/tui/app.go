// Package tui executes program in terminal
package tui

import (
	"fmt"
	"path/filepath"

	"github.com/akalankae/shift-weaver-go/internal/core"
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
	roster, err := getRosterFile(DataDir)
	if err != nil {
		panic(err)
	}
	rosterAbsPath, err := filepath.Abs(roster)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Roster: %s\n", rosterAbsPath)
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
		fmt.Print("Enter index for roster you want to read: ")
		n, _ := fmt.Scanf("%d", &fileNumber) // need fix! non digit input fucks this up
		if (fileNumber >= 0 && fileNumber < len(files)) && n == 1 {
			rosterFile = files[fileNumber]
			break
		}
		fmt.Printf("Please enter a file number between 0 and %d\n", len(files))
	}
	return // path_string, nil
}
