package decryptCommand

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Get output writer for decrypted file
func GetOutputFile(fileID string, verbose bool, decryptionFilePath string) (io.WriteCloser, string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get home directory: %w", err)
	}

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
		if verbose {
			return nil, finalOutPath, fmt.Errorf("create output: %w", err)
		}
		return nil, finalOutPath, errors.New("error while creating output path")
	}

	return out, finalOutPath, nil
}
