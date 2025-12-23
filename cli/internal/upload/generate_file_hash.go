package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Returns the SHA-256 hash of a file as a hex string
func GenerateFileHash(filePath string) (string, error) {

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	// Create a hasher
	h := sha256.New()

	// Stream file data into the hash
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("copy: %w", err)
	}

	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil

}
