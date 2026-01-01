package upload

import (
	"errors"
	"sync"
)

// Get file hash and mime type in upload command and also generate user relevant errors
func GetFileInfo(filePath string, verbose bool) (fileHash string, mimeType string, err error) {

	var wg sync.WaitGroup

	errChPhase1 := make(chan error, 2)

	wg.Add(2)

	// Generate file hash from data
	go func() {
		defer wg.Done()
		hash, err := GenerateFileHash(filePath)
		if err != nil {
			if verbose {
				errChPhase1 <- err
				return
			}
			errChPhase1 <- errors.New("error generating file hash (use --verbose for details)")
			return
		}
		fileHash = hash
	}()

	// Get the mime type of data
	go func() {
		defer wg.Done()
		mime, err := GetMime(filePath)
		if err != nil {
			if verbose {
				errChPhase1 <- err
				return
			}
			errChPhase1 <- errors.New("error generating mime type (use --verbose for details)")
			return
		}
		mimeType = mime
	}()

	wg.Wait()
	close(errChPhase1)

	// Return if any error is received
	for err := range errChPhase1 {
		if err != nil {
			return "", "", err
		}
	}

	return fileHash, mimeType, nil
}
