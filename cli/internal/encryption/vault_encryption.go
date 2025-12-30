package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Encrypt vault data using AES-GCM standard
// returns [nonce][ciphertext+tag]
func EncryptVault(vaultBytes []byte, vaultMasterKey []byte) ([]byte, error) {

	// Init AES-GCM
	block, err := aes.NewCipher(vaultMasterKey)
	if err != nil {
		return nil, fmt.Errorf("cipher init: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm init: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce: %w", err)
	}

	// Encrypt (GCM automatically appends auth tag)
	ciphertext := gcm.Seal(nil, nonce, vaultBytes, nil)

	// Append nonce in encrypted bytes
	out := append(nonce, ciphertext...)
	return out, nil
}

// Decrypt vault data
func DecryptVault(encData []byte, vaultMasterKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(vaultMasterKey)
	if err != nil {
		return nil, fmt.Errorf("cipher init: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm init: %w", err)
	}

	// Checks if the data atleast contains nonce
	nonceSize := gcm.NonceSize()
	if len(encData) < nonceSize {
		return nil, fmt.Errorf("invalid vault data")
	}

	nonce := encData[:nonceSize]
	ciphertext := encData[nonceSize:]

	// GCM Verifies integrity + password correctness here
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt vault: %w", err)
	}

	return plain, nil
}
