/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/list"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all uploaded files",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get all the files of user
		files, err := list.GetAllFiles()
		if err != nil {
			if Verbose {
				return fmt.Errorf("get files: %w", err)
			}

			return errors.New("error getting files (use --verbose for details)")
		}

		ui.ListFiles(files)

		return nil
	},
}

func init() {
	filesCmd.AddCommand(listCmd)
}
