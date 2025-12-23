/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	key  string
	name string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload <file-path>",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		filePath := args[0]
		if filePath == "" {
			return errors.New("<file-path> is required")
		}

	},
}

func init() {
	// Key flag (long: --key, short: -k)
	uploadCmd.Flags().StringVarP(&key, "key", "k", "", "Encryption key / passphrase")
	// Name flag (long: --name, short: -n)
	uploadCmd.Flags().StringVarP(&name, "name", "n", "", "Optional name for the file")

	rootCmd.AddCommand(uploadCmd)
}
