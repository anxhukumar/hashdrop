package decryptCommand

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DownloadEncryptedFile(fileUrl string, verbose bool) (io.ReadCloser, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		if verbose {
			return nil, fmt.Errorf("download failed: %w", err)
		}
		return nil, errors.New("error while downloading file data")
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if verbose {
			return nil, fmt.Errorf("download error: %s | %s", resp.Status, string(body))
		}
		return nil, errors.New("error while downloading file data")
	}

	return resp.Body, nil
}
