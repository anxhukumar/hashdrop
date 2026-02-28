/*
Copyright © 2026 Anshu Kumar

Licensed under the Apache License, Version 2.0.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
