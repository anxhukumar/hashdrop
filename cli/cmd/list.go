/*
Copyright © 2026 Anshu Kumar

Licensed under the Apache License, Version 2.0.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/files"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your uploaded files",
	Long: `
Displays all files you have uploaded to your Hashdrop account.

The list includes basic metadata such as:
• File name
• Encrypted size
• Upload status
• Key management mode (vault or passphrase)
• Creation date

Use this command to find file IDs or quickly inspect your stored files.
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get all the files of user
		files, err := files.GetAllFiles()
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
