/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/decrypt"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/verify"
	"github.com/spf13/cobra"
)

var (
	verifyFlag bool
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:          "decrypt <file-url> [destination]",
	Short:        "",
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

		// Check if vault exists
		ok, err := encryption.VaultExists()
		if err != nil {
			if Verbose {
				return fmt.Errorf("vault exists: %w", err)
			}
			return errors.New("failed to check if vault exists locally")
		}

		var DEK []byte

		// If it exists
		if ok {
			fmt.Println("🔐 Vault detected")
			// Check if the fileId exists in the vault
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

				dekString, ok := vaultData.Entries[fileID]
				if !ok {
					fmt.Println("ℹ️ This file is not stored in your vault. Switching to passphrase mode.")
					break
				}

				// Fetch data encryption key
				decoded, err := base64.StdEncoding.DecodeString(dekString)
				if err != nil {
					return fmt.Errorf("invalid DEK encoding in vault: %w", err)
				}

				DEK = decoded

				break
			}
		}

		// If vault doesn't exist or the key doesn't exist in vault
		if len(DEK) == 0 {
			fmt.Println("🔁 Using passphrase mode")

			for {

				pass, err := prompt.ReadPassword("Enter passphrase: ")
				if err != nil {
					return err
				}

				if strings.TrimSpace(pass) == "" {
					fmt.Println("passphrase cannot be empty or whitespace")
					continue
				}

				// Get salt from backend
				salt, err := decrypt.GetSalt(fileID)
				if err != nil {
					if Verbose {
						return fmt.Errorf("get salt: %w", err)
					}
					return errors.New("this file was not uploaded using passphrase mode and is not present in your vault")
				}

				// Parse string salt to bytes
				decodedSalt, err := base64.StdEncoding.DecodeString(salt)
				if err != nil {
					return fmt.Errorf("invalid salt encoding: %w", err)
				}

				fileDEK := encryption.DeriveDEK(pass, decodedSalt)

				DEK = fileDEK

				break
			}
		}

		// Decrypt file

		// Get encrypted file
		encFileData, err := http.Get(fileUrl)
		if err != nil {
			if Verbose {
				return fmt.Errorf("download failed: %w", err)
			}
			return errors.New("error while downloading file data")
		}
		defer encFileData.Body.Close()

		if encFileData.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(encFileData.Body)
			if Verbose {
				return fmt.Errorf("download error: %s | %s", encFileData.Status, string(body))
			}
			return errors.New("error while downloading file data")
		}

		// Get output file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// If user provided a path use that otherwise decrypt in download directory
		finalOutPath := ""
		if len(decryptionFilePath) != 0 {
			info, err := os.Stat(decryptionFilePath)
			if err == nil && info.IsDir() {
				// User provided a directory
				finalOutPath = filepath.Join(decryptionFilePath, fileID)
			} else {
				// User provided a file path (may not exist yet)
				finalOutPath = decryptionFilePath
			}

		} else {
			finalOutPath = filepath.Join(homeDir, "Downloads", fileID)
		}

		out, err := os.Create(finalOutPath)
		if err != nil {
			if Verbose {
				return fmt.Errorf("create output: %w", err)
			}
			return errors.New("error while creating output path")
		}
		defer out.Close()

		derivedPlaintextHash, err := encryption.DecryptAndHashStreaming(encFileData.Body, out, DEK)
		if err != nil {
			if Verbose {
				return fmt.Errorf("decrypt failed: %w", err)
			}
			return errors.New("error while decrypting file")
		}

		if verifyFlag {
			originalHashStr, err := verify.GetFileHash(fileID)
			if err != nil {
				if Verbose {
					return fmt.Errorf("get original hash: %w", err)
				}
				return errors.New("error while getting original file hash")
			}

			originalHashBytes, err := hex.DecodeString(originalHashStr)
			if err != nil {
				return fmt.Errorf("invalid hash format from server: %w", err)
			}

			// Verify hash
			if !bytes.Equal(derivedPlaintextHash, originalHashBytes) {
				fmt.Println("❌ Hash mismatch — file may be corrupted")
				return nil
			}

			fmt.Println("✅ Hash verified — file is intact")
		}

		fmt.Printf("📁 Decrypted file written to: %s\n", finalOutPath)

		return nil
	},
}

func init() {
	// (long: --verify, short: -V)
	decryptCmd.Flags().BoolVarP(&verifyFlag, "verify", "V", false, "Verify file integrity while decrypting")
	rootCmd.AddCommand(decryptCmd)
}
