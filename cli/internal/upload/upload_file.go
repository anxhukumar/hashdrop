package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
)

// Encrypts and uploads file on the fly
func UploadFileToS3(
	ctx context.Context,
	presign PresignResponse,
	localFilePath string,
	dek []byte,
) error {

	f, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	// Streaming multipart data so that http could listen simultaneously and send it
	go func() {
		defer pw.Close()

		// Write AWS required form fields
		for k, v := range presign.UploadResource.Fields {
			if err := mw.WriteField(k, v); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		// The field name must be "file" for s3 by convention
		filePart, err := mw.CreateFormFile("file", "encrypted.bin")
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// Encrypt and stream data here to filePart destination
		if err := encryption.EncryptFileStreaming(f, filePart, dek); err != nil {
			pw.CloseWithError(err)
			return
		}

		// Finish multipart
		if err := mw.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	// Send http request while taking data from the stream
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		presign.UploadResource.URL,
		pr,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s | %s", resp.Status, string(body))
	}

	return nil
}
