package upload

import (
	"fmt"
	"os"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Validates if the size of the file is within the limit
func ValidateFileSize(filePath string) error {

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fileSize := info.Size()
	limitBytes := int64(config.UploadFileSizeLimit) * 1024 * 1024

	if fileSize > limitBytes {
		return fmt.Errorf("file too large: %v bytes (limit %v bytes)", fileSize, limitBytes)
	}

	return nil
}
