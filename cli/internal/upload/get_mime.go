package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Returns the mime type of the file
func GetMime(filepath string) (string, error) {

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read file: %w", err)
	}

	return http.DetectContentType(buf[:n]), nil
}
