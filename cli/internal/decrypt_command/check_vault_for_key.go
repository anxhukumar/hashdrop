package decryptCommand

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
)

// Check vault for DEK of a file
func CheckVaultForKey(fileID string, verbose bool) ([]byte, error) {

	// Check if vault exists
	ok, err := encryption.VaultExists()
	if err != nil {
		if verbose {
			return nil, fmt.Errorf("vault exists: %w", err)
		}
		return nil, errors.New("failed to check if vault exists locally (use --verbose for details)")
	}

	if !ok {
		return nil, errors.New("vault does not exist")
	}

	// Check if the fileId exists in the vault
	var vaultMasterKey []byte
	var vaultData encryption.Vault

	for {
		// Derive vault master key
		pass, err := prompt.ReadPassword("Enter vault password: ")
		if err != nil {
			return nil, err
		}

		if strings.TrimSpace(pass) == "" {
			fmt.Println("vault password cannot be empty or whitespace")
			continue
		}

		vaultMasterKey, err = encryption.DeriveVaultMasterKey(pass)
		if err != nil {
			return nil, fmt.Errorf("Error deriving vault master key: %w", err)
		}

		// Load vault and decrypt it using vault key
		vaultData, err = encryption.LoadVault(vaultMasterKey)
		if err != nil {
			if errors.Is(err, encryption.ErrInvalidVaultKeyOrCorrupted) {
				fmt.Println("Failed to unlock vault. The password may be incorrect or the vault file may be corrupted.")
				continue
			}
			if errors.Is(err, encryption.ErrVaultNotFound) {
				return nil, err
			}
		}

		dekString, ok := vaultData.Entries[fileID]
		if !ok {
			fmt.Println("ℹ️ This file is not stored in your vault.")
			return nil, nil
		}

		// Fetch data encryption key
		decoded, err := base64.StdEncoding.DecodeString(dekString)
		if err != nil {
			return nil, fmt.Errorf("invalid DEK encoding in vault: %w", err)
		}
		return decoded, nil
	}
}
