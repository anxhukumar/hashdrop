/*
Copyright © 2026 Anshu Kumar

Licensed under the Apache License, Version 2.0.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:          "logout",
	Short:        "Log out of Hashdrop",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Revoke refresh token
		if err := auth.RevokeRefreshToken(); err != nil && Verbose {
			fmt.Println("Warning:", err)
		}

		// Call Delete tokens function
		if err := auth.DeleteTokens(); err != nil {
			if Verbose {
				return err
			}
			return errors.New("logout failed (use --verbose for details)")
		}

		fmt.Println("Logged out successfully.")

		return nil
	},
}

func init() {
	authCmd.AddCommand(logoutCmd)
}
