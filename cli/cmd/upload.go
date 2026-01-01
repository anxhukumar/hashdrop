/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/upload"
	"github.com/spf13/cobra"
)

var (
	noVault  bool
	fileName string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:          "upload <file-path>",
	Short:        "",
	Long:         ``,
	SilenceUsage: true,
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
		fileHash, mimeType, err := upload.GetFileInfo(filePath, Verbose)
		if err != nil {
			return err
		}

		// Create DEK for vault and no-vault users
		fileDEK, fileSalt, err := encryption.ObtainFileEncryptionKey(noVault, Verbose)
		if err != nil {
			return err
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

		// Check if vault exists, create if it doesn't and update it
		if err = encryption.CreateAndUpdateVault(fileDEK, presignResource.FileID, noVault, Verbose); err != nil {
			return err
		}

		// Encrypt and upload file (streaming)

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
