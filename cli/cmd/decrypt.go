/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	decryptCommand "github.com/anxhukumar/hashdrop/cli/internal/decrypt_command"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	verifyFlag bool
	vaultFlag  bool
	keyFlag    bool
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt <file-url> [destination]",
	Short: "Decrypt a file from a Hashdrop download link",
	Long: `
Decrypts a file that was shared or downloaded via a Hashdrop URL.

You provide the file's download link and Hashdrop will:
• Retrieve the encrypted file
• Ask for the required decryption secret (vault, passphrase, or raw key)
• Decrypt the file locally
• Optionally verify its integrity using the original hash

By default, the decrypted file is saved to your Downloads directory.
You may optionally provide a destination path or directory.

Examples:
  hashdrop decrypt https://api.hashdrop.dev/... 
  hashdrop decrypt https://api.hashdrop.dev/... ./output.txt
  hashdrop decrypt --key https://api.hashdrop.dev/... --verify
`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			return errors.New("<file-url> is required")
		}
		fileUrl := args[0]

		// Optional argument of destination
		var decryptionFilePath string

		if len(args) == 2 {
			decryptionFilePath = args[1]
		}

		// Get fileID from the url
		u, err := url.Parse(fileUrl)
		if err != nil {
			return fmt.Errorf("invalid file URL: %w", err)
		}
		fileID := path.Base(u.Path)

		var DEK []byte

		// Block if both flags are used
		if vaultFlag && keyFlag {
			return errors.New("choose only one mode: --vault or --key")
		}

		// If user didn't select any mode with flag, show options
		if !vaultFlag && !keyFlag {
			decryptMode, err := decryptCommand.ShowDecryptionOptions()
			if err != nil {
				return err
			}

			switch decryptMode {
			case decryptCommand.VaultDecryptMode:
				vaultFlag = true
			case decryptCommand.KeyDecryptMode:
				keyFlag = true
			}
		}

		// If user selected vault
		// Check in vault if the file exists and return the DEK
		if vaultFlag {
			DEK, err = decryptCommand.CheckVaultForKey(fileID, Verbose)
			if err != nil {
				return err
			}

			if len(DEK) == 0 {
				return errors.New("file not found in vault. Try --key mode instead")
			}
		}

		// If user selects key mode, process passphrase or key directly
		if keyFlag {
			DEK, err = decryptCommand.GetEncryptionKeyFromUser(fileID)
			if err != nil {
				return err
			}
		}

		// Decrypt file

		// Get encrypted file
		encFileData, totalSize, err := decryptCommand.DownloadEncryptedFile(fileUrl, Verbose)
		if err != nil {
			return err
		}
		defer encFileData.Close()

		fmt.Println("📥 downloading")

		// Progress bar
		bar := progressbar.NewOptions64(
			totalSize,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(15),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stderr, "\n")
			}),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[cyan]=[reset]", // Using cyan for download to distinguish from upload
				SaucerHead:    "[cyan]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		proxyReader := progressbar.NewReader(encFileData, bar)

		// If user provided a path use that otherwise decrypt in download directory
		out, finalOutPath, err := decryptCommand.GetOutputFile(fileID, Verbose, decryptionFilePath)
		if err != nil {
			return err
		}
		defer out.Close()

		derivedPlaintextHash, err := encryption.DecryptAndHashStreaming(&proxyReader, out, DEK)
		if err != nil {
			if Verbose {
				return fmt.Errorf("decrypt failed: %w", err)
			}
			return errors.New("error while decrypting file (use --verbose for details)")
		}

		if verifyFlag {
			err := decryptCommand.VerifyHash(fileID, Verbose, derivedPlaintextHash)
			if err != nil {
				return err
			}
		}

		fmt.Printf("📁 Decrypted file written to: %s\n", finalOutPath)

		return nil
	},
}

func init() {
	decryptCmd.Flags().BoolVar(&vaultFlag, "vault", false, "Use vault mode for decryption")
	decryptCmd.Flags().BoolVar(&keyFlag, "key", false, "Use key directly for decryption")
	// (long: --verify, short: -V)
	decryptCmd.Flags().BoolVarP(&verifyFlag, "verify", "V", false, "Verify file integrity while decrypting")
	rootCmd.AddCommand(decryptCmd)
}
