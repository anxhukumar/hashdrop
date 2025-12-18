/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// LoginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:          "login",
	Short:        "Login to your hashdrop account",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, err := prompt.ReadLine("Email: ")
		if err != nil {
			return err
		}

		password, err := prompt.ReadPassword("Password: ")
		if err != nil {
			return err
		}

		// Call Login api function
		if err := auth.Login(email, password); err != nil {
			if Verbose {
				return err
			}
			return errors.New("login failed (use --verbose for details)")
		}

		fmt.Println("✓ Logged in successfully")

		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
}
