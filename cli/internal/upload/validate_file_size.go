package upload

import (
	"errors"
	"fmt"
	"os"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Validates if the size of the file is within the limit
func ValidateFileSize(filePath string, Verbose bool) error {

	info, err := os.Stat(filePath)
	if err != nil {
		if Verbose {
			return fmt.Errorf("failed to read file: %w", err)
		}
		return errors.New("failed to read file (use --verbose for details)")
	}

	fileSize := info.Size()
	limitBytes := int64(config.UploadFileSizeLimit) * 1024 * 1024

	if fileSize == 0 {
		return fmt.Errorf("file is empty")
	}

	if fileSize > limitBytes {
		return fmt.Errorf(
			"file too large: %.2f MB (limit %.2f MB)",
			float64(fileSize)/(1024*1024),
			float64(config.UploadFileSizeLimit))
	}

	return nil
}
