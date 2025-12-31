package encryption

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"golang.org/x/crypto/argon2"
)

// Generate Vault Master Key from password string
func GenerateVaultMasterKey(password string) (key []byte, err error) {
	password_bytes := []byte(password)

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key = argon2.IDKey(
		password_bytes,
		salt,
		config.ArgonTime,
		config.ArgonMemory,
		config.ArgonThreads,
		config.ArgonKeyLen,
	)

	vaultMetaData := VaultKeyMetadata{
		Version: 1,
		Argon: ArgonParams{
			Time:    config.ArgonTime,
			Memory:  config.ArgonMemory,
			Threads: config.ArgonThreads,
			KeyLen:  config.ArgonKeyLen,
		},
		Salt: salt,
	}

	// create vault_meta.json file
	if err := StoreVaultMetadata(vaultMetaData); err != nil {
		return nil, err
	}

	return key, nil
}

// Derive Vault Master Key
func DeriveVaultMasterKey(password string) ([]byte, error) {
	password_bytes := []byte(password)

	meta, err := LoadVaultMetadata()
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey(
		password_bytes,
		meta.Salt,
		meta.Argon.Time,
		meta.Argon.Memory,
		meta.Argon.Threads,
		meta.Argon.KeyLen,
	)

	return key, nil
}

// Get vault meta data path
func vaultMetadataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, config.ConfigDirName, config.VaultMetadataFileName), nil
}

// Store the meta data in the users home directory ~/.hashdrop/vault_meta.json
func StoreVaultMetadata(metadata VaultKeyMetadata) error {

	// Get path
	vaultMetaPath, err := vaultMetadataPath()
	if err != nil {
		return err
	}

	// Marshal the data with proper indentation
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault metadata: %w", err)
	}

	// Create the vault_meta.json
	if err := os.WriteFile(vaultMetaPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write vault metadata file: %w", err)
	}

	return nil
}

// Load vault metadata
func LoadVaultMetadata() (VaultKeyMetadata, error) {
	var res VaultKeyMetadata

	path, err := vaultMetadataPath()
	if err != nil {
		return res, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return res, fmt.Errorf("failed to read vault metadata file: %w", err)
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return res, fmt.Errorf("failed to decode vault metadata json: %w", err)
	}

	return res, nil
}
