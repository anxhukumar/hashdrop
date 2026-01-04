/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/files"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	revealKey bool
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:          "show <file-id>",
	Short:        "Show details of an uploaded file",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			return errors.New("<file-id> is required")
		}
		fileID := args[0]

		// Get details of a file
		file, err := files.GetDetailedFile(fileID)
		if err != nil {
			if Verbose {
				return fmt.Errorf("show file: %w", err)
			}
			return errors.New("error getting file (use --verbose for details)")
		}

		// If we end up receiving more than one file then show the full id of all those files
		if len(file) > 1 {
			ui.ShowMultipleFileMatches(file)
			return nil
		}

		// If reveal key flag is provided show the reveal key of the file from vault after decrypting the vault
		if revealKey {
			var vaultMasterKey []byte
			var vaultData encryption.Vault

			for {
				// Derive vault master key
				pass, err := prompt.ReadPassword("Enter vault password: ")
				if err != nil {
					return err
				}

				if strings.TrimSpace(pass) == "" {
					fmt.Println("vault password cannot be empty or whitespace")
					continue
				}

				vaultMasterKey, err = encryption.DeriveVaultMasterKey(pass)
				if err != nil {
					return fmt.Errorf("Error deriving vault master key: %w", err)
				}

				// Load vault and decrypt it using vault key
				vaultData, err = encryption.LoadVault(vaultMasterKey)
				if err != nil {
					if errors.Is(err, encryption.ErrInvalidVaultKeyOrCorrupted) {
						fmt.Println("Failed to unlock vault. The password may be incorrect or the vault file may be corrupted.")
						continue
					}
					if errors.Is(err, encryption.ErrVaultNotFound) {
						return err
					}
				}

				break
			}

			key, ok := vaultData.Entries[file[0].ID.String()]
			if !ok {
				ui.NoEncryptionKey()
				return nil
			}
			ui.ShowFile(file[0], key)
			return nil
		}

		ui.ShowFile(file[0], "")

		return nil
	},
}

func init() {

	// Key flag (long: --reveal-key, short: -R)
	showCmd.Flags().BoolVarP(&revealKey, "reveal-key", "R", false, "Reveal encryption key of a file")

	filesCmd.AddCommand(showCmd)
}
