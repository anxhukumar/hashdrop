/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
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

		// Confirm otp to verify account
		for {
			otpInput, err := prompt.ReadLine("Enter your email verification code: ")
			if err != nil {
				return err
			}

			verifyRequestParams := &struct {
				Email string `json:"email"`
				OTP   string `json:"otp"`
			}{
				Email: email,
				OTP:   otpInput,
			}

			status, err := api.PatchJSON(config.VerifyUserEndpoint, verifyRequestParams, nil, "")
			if err != nil {
				switch status {
				case 401:
					fmt.Println("❌ Invalid verification code. Please try again.")
					continue
				case 400:
					return errors.New("verification code expired. Please register again or request a new code")
				default:
					if Verbose {
						return err
					}
					return errors.New("verification failed (use --verbose for details)")
				}
			}

			switch status {
			case 204:
				fmt.Println("✓ Registration successful")
				return nil
			default:
				return fmt.Errorf("unexpected server response: %d", status)
			}

		}

	},
}

func init() {
	authCmd.AddCommand(registerCmd)
}
