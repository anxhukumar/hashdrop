package encryption

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Get vault path
func vaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, config.ConfigDirName, config.VaultFileName), nil
}

// Stores the encrypted vault data in the users home directory ~/.hashdrop/vault.enc
func EncryptAndStoreVault(vaultData Vault, vaultMasterKey []byte) error {

	if len(vaultMasterKey) != 32 {
		return fmt.Errorf("invalid vault master key length: %d", len(vaultMasterKey))
	}

	// Marshal the data with proper indentation
	data, err := json.Marshal(vaultData)
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	// Encrypt json bytes
	encData, err := EncryptVault(data, vaultMasterKey)
	if err != nil {
		return fmt.Errorf("encrypt vault: %w", err)
	}

	path, err := vaultPath()
	if err != nil {
		return err
	}

	// Store encrypted data in vault file
	if err := os.WriteFile(path, encData, 0600); err != nil {
		return fmt.Errorf("failed to write encrypted vault file: %w", err)
	}

	return nil
}

// Load vault
func LoadVault(vaultMasterKey []byte) (Vault, error) {

	var vaultData Vault

	if len(vaultMasterKey) != 32 {
		return vaultData, fmt.Errorf("invalid vault master key length: %d", len(vaultMasterKey))
	}

	path, err := vaultPath()
	if err != nil {
		return vaultData, err
	}

	encData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return vaultData, ErrVaultNotFound
		}
		return vaultData, fmt.Errorf("failed to read vault file: %w", err)
	}

	// Decrypt vault data
	decData, err := DecryptVault(encData, vaultMasterKey)
	if err != nil {
		if errors.Is(err, ErrInvalidVaultKeyOrCorrupted) {
			return vaultData, ErrInvalidVaultKeyOrCorrupted
		}
		return vaultData, fmt.Errorf("decrypt vault: %w", err)
	}

	// Decode decrypted json data
	if err := json.Unmarshal(decData, &vaultData); err != nil {
		return vaultData, fmt.Errorf("failed to decode vault json: %w", err)
	}

	return vaultData, nil
}

// Delete vault
func DeleteVault() error {
	path, err := vaultPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove vault file: %w", err)
	}
	return nil
}

// Check if vault exists
func VaultExists() (bool, error) {
	path, err := vaultPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("failed to stat vault: %w", err)
}
