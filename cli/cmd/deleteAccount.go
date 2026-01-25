/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// deleteAccountCmd represents the deleteAccount command
var deleteAccountCmd = &cobra.Command{
	Use:   "delete-account",
	Short: "Permanently delete your account and all associated data",
	Long: `
Permanently deletes your Hashdrop account and all associated data.

This action will:
• Delete your user account
• Remove all uploaded files from secure storage
• Revoke all access and refresh tokens
• Invalidate all shared or download links
• Remove your local Hashdrop configuration

This operation is irreversible.
Once completed, your data cannot be recovered.
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Confirm deletion from user
		confirmMsg, err := prompt.ReadLine(
			`⚠️  PERMANENT ACCOUNT DELETION

This will:
• Delete your account
• Permanently remove all uploaded files
• Revoke all access tokens
• Make all shared links invalid

This action CANNOT be undone.

Type [DELETE ALL MY DATA] to confirm: `)
		if err != nil {
			return err
		}

		if strings.ToLower(strings.TrimSpace(confirmMsg)) != "delete all my data" {
			fmt.Println("Aborted.")
			return nil
		}

		// Get access token
		token, err := auth.EnsureAccessToken()
		if err != nil {
			if Verbose {
				return fmt.Errorf("error authenticating user: %w", err)
			}
			return errors.New("error authenticating user (use --verbose for details)")
		}

		err = api.Delete(config.DeleteUserEndpoint, token, nil)
		if err != nil {
			if Verbose {
				return fmt.Errorf("error deleting account: %w", err)
			}
			return errors.New("error deleting account (use --verbose for details)")
		}

		// Delete local .hashdrop directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("✔ Account deleted, but could not locate your home directory to remove local data.")
			fmt.Println("    If present, please delete the folder '~/.hashdrop' manually.")
			return nil
		} else {
			hashdropDir := filepath.Join(homeDir, ".hashdrop")
			if err := os.RemoveAll(hashdropDir); err != nil {
				fmt.Println("✔ Account deleted, but failed to remove local directory (~/.hashdrop).")
				fmt.Println("    You may delete it manually.")
				return nil
			}
		}

		fmt.Println("✔ Your account and all associated data have been permanently deleted.")

		return nil
	},
}

func init() {
	authCmd.AddCommand(deleteAccountCmd)
}
