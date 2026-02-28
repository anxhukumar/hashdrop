package upload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/schollz/progressbar/v3"
)

// These are byte values to predict chunk size
const nonceSizeVal = 12
const gcmTagVal = 16
const lenFieldVal = 4

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

	// Streaming data so that http could listen simultaneously and send it
	go func() {
		defer pw.Close()

		// Encrypt and stream data to destination
		if err := encryption.EncryptFileStreaming(f, pw, dek); err != nil {
			pw.CloseWithError(err)
		}
	}()

	// Calculate total size
	info, err := os.Stat(localFilePath)
	if err != nil {
		return err
	}

	plainSize := info.Size()

	const chunkSize = config.FileStreamingChunkSize
	nonceSize := nonceSizeVal
	gcmTag := gcmTagVal
	lenField := lenFieldVal

	numChunks := (plainSize + chunkSize - 1) / chunkSize
	overheadPerChunk := int64(nonceSize + lenField + gcmTag)
	encryptedSize := plainSize + numChunks*overheadPerChunk

	// Progress bar configuration
	bar := progressbar.NewOptions64(
		encryptedSize,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	proxyReader := progressbar.NewReader(pr, bar)

	// Send http request while taking data from the stream
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		presign.UploadResource.URL,
		&proxyReader,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	req.ContentLength = encryptedSize

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s | %s", resp.Status, string(body))
	}

	return nil
}
