/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/upload"
	"github.com/spf13/cobra"
)

var (
	key  string
	name string
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

		// Prompt the user and ask for encryption key if they didn't add or is invalid
		if strings.TrimSpace(key) == "" || len(key) < config.MinEncryptionKeyLen {

			// Check if user is here only for the key length
			if len(key) > 0 {
				fmt.Printf("key must be at least %d characters long\n", config.MinEncryptionKeyLen)
			}

			for {

				k, err := prompt.ReadPassword("Enter encryption key: ")
				if err != nil {
					return err
				}
				key = k

				if strings.TrimSpace(key) == "" {
					fmt.Println("key cannot be empty or whitespace")
					continue
				}

				// Check if the key length is valid
				if len(key) < config.MinEncryptionKeyLen {
					fmt.Printf("key must be at least %d characters long\n", config.MinEncryptionKeyLen)
					continue
				}

				confirmKey, err := prompt.ReadPassword("Confirm key: ")
				if err != nil {
					return err
				}

				if key != confirmKey {
					fmt.Println("keys do not match, try again")
					continue
				}

				break
			}

		}

		// Prompt the user and show the default file name to user if they didn't use the flag
		// They can edit it here if they want
		if name == "" {
			defaultFileName := filepath.Base(filePath)
			n, err := prompt.ReadLine(fmt.Sprintf("File name (press Enter to keep %q): ", defaultFileName))
			if err != nil {
				return err
			}
			if n == "" {
				name = defaultFileName
			} else {
				name = n
			}

		}

		errChPhase2 := make(chan error, 2)

		wg.Add(2)

		// Send mime type and file name and get presign link
		var presignResource upload.PresignResponse
		go func() {
			defer wg.Done()
			err := upload.GetPresignedLink(name, mimeType, &presignResource)
			if err != nil {
				if Verbose {
					errChPhase2 <- err
					return
				}
				errChPhase2 <- errors.New("error getting upload link from server (use --verbose for details)")
				return
			}

		}()

		// Encrypt data

		// wait here

		// prmpt the user to keep the key safe
		// Upload the encrypted data
		// once the data is uploaded successfully send teh full metadata to backend for storage

		// get cloudfront url at the end in response to download encrypted data

	},
}

func init() {
	// Key flag (long: --key, short: -k)
	uploadCmd.Flags().StringVarP(&key, "key", "k", "", "Encryption key / passphrase")
	// Name flag (long: --name, short: -n)
	uploadCmd.Flags().StringVarP(&name, "name", "n", "", "Optional name for the file")

	rootCmd.AddCommand(uploadCmd)
}
