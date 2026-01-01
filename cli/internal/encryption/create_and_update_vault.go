package encryption

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
	"github.com/google/uuid"
)

// Create vault if it doesn't exist, updates it and also generate user relevant errors
func CreateAndUpdateVault(fileDEK []byte, fileID uuid.UUID, Verbose bool) error {

	// Check if vault exists
	exists, err := VaultExists()
	if err != nil {
		if Verbose {
			return fmt.Errorf("vault exists : %w", err)
		}
		return errors.New("can't check if vault exists (use --verbose for details)")
	}

	// If user chooses vault mode and vault does not exist then create it.
	// If it exists update the vault with new (fileID -> fileDEK).
	if !exists {
		ui.PrintVaultCreationInfo()
		fmt.Scanln() // waits until Enter is pressed to continue
		var vaultMasterKey []byte
		var vaultData Vault
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

			vaultMasterKey, err = GenerateVaultMasterKey(pass)
			if err != nil {
				if Verbose {
					return fmt.Errorf("generate vault key: %w", err)
				}
				return errors.New("error generating vault master key (use --verbose for details)")
			}

			// Create vault data and store it (fileID -> FileDEK)
			vaultData = Vault{
				Version: config.VaultVersion,
				Entries: make(map[string]string),
			}

			vaultData.Entries[fileID.String()] = base64.StdEncoding.EncodeToString(fileDEK)

			if err = EncryptAndStoreVault(vaultData, vaultMasterKey); err != nil {
				if Verbose {
					return fmt.Errorf("creating vault: %w", err)
				}
				return errors.New("error while creating vault (use --verbose for details)")
			}

			break
		}
	} else {
		var vaultMasterKey []byte
		var vaultData Vault

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

			vaultMasterKey, err = DeriveVaultMasterKey(pass)
			if err != nil {
				return fmt.Errorf("Error deriving vault master key: %w", err)
			}

			// Load vault and decrypt it using vault key
			vaultData, err = LoadVault(vaultMasterKey)
			if err != nil {
				if errors.Is(err, ErrInvalidVaultKeyOrCorrupted) {
					fmt.Println("Failed to unlock vault. The password may be incorrect or the vault file may be corrupted.")
					continue
				}
				if errors.Is(err, ErrVaultNotFound) {
					return err
				}
			}

			break
		}

		// Update the vault with new (fileID -> fileDEK) and store
		vaultData.Entries[fileID.String()] = base64.StdEncoding.EncodeToString(fileDEK)

		if err = EncryptAndStoreVault(vaultData, vaultMasterKey); err != nil {
			if Verbose {
				return fmt.Errorf("updating vault: %w", err)
			}
			return errors.New("error while updating vault")
		}
	}

	return nil
}
