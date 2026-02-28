/*
Copyright © 2026 Anshu Kumar

Licensed under the Apache License, Version 2.0.
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
