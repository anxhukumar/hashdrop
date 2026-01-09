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
	"github.com/anxhukumar/hashdrop/cli/internal/files"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
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

		if strings.ToLower(strings.TrimSpace(userConfirmation)) != "y" {
			fmt.Println("Aborted.")
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

		// Check if there are multiple matches of the short FileID
		fileMatches, err := files.CheckMultipleShortFileIDMatch(fileID, queryParam, token)
		if err != nil {
			return err
		}

		if len(fileMatches) > 1 {
			ui.ShowMultipleFileMatches(fileMatches)
			return nil
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
