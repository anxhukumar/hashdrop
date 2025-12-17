/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hashdrop",
	Short: "Secure file drop with end-to-end encryption",
	Long: `Hashdrop is a zero-trust CLI tool for sharing sensitive files.

Files are encrypted client-side, stored as unreadable blobs,
and shared via expiring links with integrity verification.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

}
