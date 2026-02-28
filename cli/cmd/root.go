/*
Copyright © 2026 Anshu Kumar

Licensed under the Apache License, Version 2.0.
*/
package cmd

import (
	"os"

	"github.com/anxhukumar/hashdrop/cli/internal/cliversion"
	"github.com/spf13/cobra"
)

var (
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "hashdrop",
	Short: "Secure file drop with end-to-end encryption",
	Long: `Hashdrop is a zero-trust CLI tool for sharing sensitive files.

Files are encrypted client-side, stored as unreadable blobs,
and shared via links with integrity verification.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch cmd.Name() {
		case "help", "completion", "version":
			return nil
		}
		return cliversion.CheckCliVersion(Verbose)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Verbose flag (long: --verbose, short: -v)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")

}
