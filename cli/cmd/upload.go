/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/anxhukumar/hashdrop/cli/internal/upload"
	"github.com/spf13/cobra"
)

var (
	noVault  bool
	fileName string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload <file-path>",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			return errors.New("<file-path> is required")
		}
		filePath := args[0]

		// Validate file size
		if err := upload.ValidateFileSize(filePath, Verbose); err != nil {
			return err
		}

		// Concurrently handle hash and mime generation
		var wg sync.WaitGroup

		errChPhase1 := make(chan error, 2)

		wg.Add(2)

		var fileHash string
		var mimeType string

		// Generate file hash from data
		go func() {
			defer wg.Done()
			hash, err := upload.GenerateFileHash(filePath)
			if err != nil {
				if Verbose {
					errChPhase1 <- err
					return
				}
				errChPhase1 <- errors.New("error generating file hash (use --verbose for details)")
				return
			}
			fileHash = hash
		}()

		// Get the mime type of data
		go func() {
			defer wg.Done()
			mime, err := upload.GetMime(filePath)
			if err != nil {
				if Verbose {
					errChPhase1 <- err
					return
				}
				errChPhase1 <- errors.New("error generating mime type (use --verbose for details)")
				return
			}
			mimeType = mime
		}()

		wg.Wait()
		close(errChPhase1)

		// Return if any error is received
		for err := range errChPhase1 {
			if err != nil {
				return err
			}
		}

		// File encryption key and salt
		var fileDEK []byte  // no-vault + vault
		var fileSalt []byte // only for no-vault users

		// Create DEK with and without vault
		if noVault {

			ui.PrintNoVaultWarning()
			fmt.Scanln() // waits until Enter is pressed to continue
			for {

				pass, err := prompt.ReadPassword("Enter passphrase: ")
				if err != nil {
					return err
				}

				if strings.TrimSpace(pass) == "" {
					fmt.Println("passphrase cannot be empty or whitespace")
					continue
				}

				// Check if the key length is valid
				if len(pass) < config.MinCustomEncryptionKeyLen {
					fmt.Printf("passphrase must be at least %d characters long\n", config.MinCustomEncryptionKeyLen)
					continue
				}

				confirmPass, err := prompt.ReadPassword("Confirm: ")
				if err != nil {
					return err
				}

				if pass != confirmPass {
					fmt.Println("passphrase do not match, try again")
					continue
				}

				fileDEK, fileSalt, err = encryption.GenerateDEKfromPassphrase(pass)
				if err != nil {
					if Verbose {
						return fmt.Errorf("generate key from passphrase: %w", err)
					}
					return errors.New("error generating key from passphrase (use --verbose for details)")
				}

				break
			}

		} else {
			var err error

			fileDEK, err = encryption.GenerateRandomDEK()
			if err != nil {
				if Verbose {
					return fmt.Errorf("generate random key: %w", err)
				}
				return errors.New("error generating random key (use --verbose for details)")
			}
		}

		// Prompt the user and show the default file name to user if they didn't use the flag
		// They can edit it here if they want
		if fileName == "" {
			defaultFileName := filepath.Base(filePath)
			n, err := prompt.ReadLine(fmt.Sprintf("File name (press Enter to keep %q): ", defaultFileName))
			if err != nil {
				return err
			}
			if n == "" {
				fileName = defaultFileName
			} else {
				fileName = n
			}
		}

		// Send mime type and file name and get presign link
		var presignResource upload.PresignResponse

		if err := upload.GetPresignedLink(fileName, mimeType, &presignResource); err != nil {
			if Verbose {
				return err
			}
			return errors.New("error getting presigned link (use --verbose for details)")
		}

		// Check if vault exists
		exists, err := encryption.VaultExists()
		if err != nil {
			if Verbose {
				return fmt.Errorf("vault exists : %w", err)
			}
			return errors.New("can't check if vault exists (use --verbose for details)")
		}

		// If user chooses vault mode and vault does not exist then create it.
		// If it exists update the vault with new (fileID -> fileDEK).
		if !noVault {
			if !exists {
				ui.PrintVaultCreationInfo()
				fmt.Scanln() // waits until Enter is pressed to continue
				var vaultMasterKey []byte
				var vaultData encryption.Vault
				for {

					pass, err := prompt.ReadPassword("Enter vault password: ")
					if err != nil {
						return err
					}

					if strings.TrimSpace(pass) == "" {
						fmt.Println("vault password cannot be empty or whitespace")
						continue
					}

					// Check if the key length is valid
					if len(pass) < config.MinVaultPasswordLen {
						fmt.Printf("vault password must be at least %d characters long\n", config.MinVaultPasswordLen)
						continue
					}

					confirmPass, err := prompt.ReadPassword("Confirm: ")
					if err != nil {
						return err
					}

					if pass != confirmPass {
						fmt.Println("passwords do not match, try again")
						continue
					}

					vaultMasterKey, err = encryption.GenerateVaultMasterKey(pass)
					if err != nil {
						if Verbose {
							return fmt.Errorf("generate vault key: %w", err)
						}
						return errors.New("error generating vault master key (use --verbose for details)")
					}

					// Create vault data and store it (fileID -> FileDEK)
					vaultData = encryption.Vault{
						Version: config.VaultVersion,
						Entries: make(map[string]string),
					}

					vaultData.Entries[presignResource.FileID.String()] = base64.StdEncoding.EncodeToString(fileDEK)

					if err = encryption.EncryptAndStoreVault(vaultData, vaultMasterKey); err != nil {
						if Verbose {
							return fmt.Errorf("creating vault: %w", err)
						}
						return errors.New("error while creating vault (use --verbose for details)")
					}

					break
				}
			} else {
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

				// Update the vault with new (fileID -> fileDEK) and store
				vaultData.Entries[presignResource.FileID.String()] = base64.StdEncoding.EncodeToString(fileDEK)

				if err = encryption.EncryptAndStoreVault(vaultData, vaultMasterKey); err != nil {
					if Verbose {
						return fmt.Errorf("updating vault: %w", err)
					}
					return errors.New("error while updating vault")
				}
			}
		}

		// Encrypt and upload file

		// Cancel if user presses ctrl+C
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Enforce max upload time
		ctx, cancel := context.WithTimeout(ctx, config.MaxTimeAllowedToUploadFile*time.Minute)
		defer cancel()

		if err := upload.UploadFileToS3(
			ctx,
			presignResource,
			filePath,
			fileDEK,
		); err != nil {
			if Verbose {
				return fmt.Errorf("upload file: %w", err)
			}
			return errors.New("error while uploading file (use --verbose for details)")
		}

		// once the data is uploaded successfully send teh full metadata to backend for storage

		// get cloudfront url at the end in response to download encrypted data

	},
}

func init() {
	// Key flag (long: --key, short: -k)
	uploadCmd.Flags().BoolVarP(&noVault, "no-vault", "N", false, "Disable local key vault. Use a self-managed encryption passphrase. If lost, the file cannot be decrypted.")
	// Name flag (long: --name, short: -n)
	uploadCmd.Flags().StringVarP(&fileName, "name", "n", "", "Optional name for the file")

	rootCmd.AddCommand(uploadCmd)
}
