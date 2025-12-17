/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new hashdrop account",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, err := prompt.ReadLine("Email: ")
		if err != nil {
			return err
		}

		password, err := prompt.ReadPassword("Password: ")
		if err != nil {
			return err
		}

		confirm, err := prompt.ReadPassword("Confirm Password: ")
		if err != nil {
			return err
		}

		if password != confirm {
			return errors.New("passwords do not match")
		}

		return auth.Register(email, password)
	},
}

func init() {
	authCmd.AddCommand(registerCmd)
}
