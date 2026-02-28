package decryptCommand

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DownloadEncryptedFile(fileUrl string, verbose bool) (io.ReadCloser, int64, error) {
	resp, err := http.Get(fileUrl)
	if err != nil {
		if verbose {
			return nil, 0, fmt.Errorf("download failed: %w", err)
		}
		return nil, 0, errors.New("error while downloading file data")
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if verbose {
			return nil, 0, fmt.Errorf("download error: %s | %s", resp.Status, string(body))
		}
		return nil, 0, errors.New("error while downloading file data")
	}

	return resp.Body, resp.ContentLength, nil
}
