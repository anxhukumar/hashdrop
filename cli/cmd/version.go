/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the current Hashdrop CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.CurrentCliVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
