/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File commands",
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
