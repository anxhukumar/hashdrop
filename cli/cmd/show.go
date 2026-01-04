/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/files"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:          "show <file-id>",
	Short:        "Show details of an uploaded file",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			return errors.New("<file-id> is required")
		}
		fileID := args[0]

		// Get details of a file
		file, err := files.GetDetailedFile(fileID)
		if err != nil {
			if Verbose {
				return fmt.Errorf("show file: %w", err)
			}
			return errors.New("error getting file (use --verbose for details)")
		}

		// If we end up receiving more than one file then show the full id of all those files
		if len(file) > 1 {
			ui.ShowMultipleFileMatches(file)
			return nil
		}

		ui.ShowFile(file[0])

		return nil
	},
}

func init() {
	filesCmd.AddCommand(showCmd)
}
