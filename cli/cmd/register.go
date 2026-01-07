/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// RegisterCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new Hashdrop account",
	Long: `
Registers a new Hashdrop account using your email and password.

You will be prompted to enter and confirm your password.
Once registration succeeds, you can log in and start uploading encrypted files.
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, err := prompt.ReadLine("Email: ")
		if err != nil {
			return err
		}

		fmt.Println("\nPassword requirements:")
		fmt.Printf("• Minimum length: %d characters\n\n", config.MinPasswordLen)

		password, err := prompt.ReadPassword("Password: ")
		if err != nil {
			return err
		}

		// Check if the user inserted correct length of the password
		if len(password) < config.MinPasswordLen {
			return fmt.Errorf("password must be at least %d characters long", config.MinPasswordLen)
		}

		confirm, err := prompt.ReadPassword("Confirm Password: ")
		if err != nil {
			return err
		}

		if password != confirm {
			return errors.New("passwords do not match")
		}

		// Call Register api function
		if err := auth.Register(email, password); err != nil {
			if Verbose {
				return err
			}
			return errors.New("registration failed (use --verbose for details)")
		}

		fmt.Println("✓ Registration successful")

		return nil

	},
}

func init() {
	authCmd.AddCommand(registerCmd)
}
