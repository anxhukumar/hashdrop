package ui

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/files"
)

func ShowMultipleFileMatches(matches []files.FileDetailedData) {

	fmt.Println()
	fmt.Println("──────────────────────────────────────────────")
	fmt.Println("The prefix you entered matches more than one file.")
	fmt.Println("Please choose the full file ID.")
	fmt.Println()
	fmt.Println("Matching files:")

	for _, f := range matches {
		fmt.Printf("  • %s\n", f.FileName)
		fmt.Printf("    ID: %s\n\n", f.ID.String())
	}

	fmt.Println("──────────────────────────────────────────────")
	fmt.Println()
}
