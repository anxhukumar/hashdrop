/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <file-id>",
	Short:        "Delete an uploaded file",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			return errors.New("<file-id> is required")
		}
		fileID := args[0]

		userConfirmation, err := prompt.ReadLine("Are you sure you want to permanently delete this file? (y/n) ")
		if err != nil {
			return err
		}

		if strings.ToLower(userConfirmation) != "y" {
			return nil
		}

		// Get access token
		token, err := auth.EnsureAccessToken()
		if err != nil {
			if Verbose {
				return fmt.Errorf("error authenticating user: %w", err)
			}
			return errors.New("error authenticating user")
		}

		queryParam := map[string]string{
			"id": fileID,
		}

		err = api.Delete(config.DeleteFileEndpoint, token, queryParam)
		if err != nil {
			if Verbose {
				return fmt.Errorf("error deleting: %w", err)
			}
			return errors.New("failed to delete file")
		}

		return nil

	},
}

func init() {
	filesCmd.AddCommand(deleteCmd)
}
