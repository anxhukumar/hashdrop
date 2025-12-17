/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your hashdrop account",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, err := prompt.ReadLine("Email: ")
		if err != nil {
			return err
		}

		password, err := prompt.ReadPassword("Password: ")
		if err != nil {
			return err
		}

		return auth.Login(email, password)
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
