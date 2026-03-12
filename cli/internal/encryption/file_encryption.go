package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Encrypt file in chunks. [ nonce ][ 4-byte length ][ ciphertext+tag ]
func EncryptFileStreaming(src io.Reader, dst io.Writer, dek []byte) error {

	if len(dek) != 32 {
		return fmt.Errorf("invalid DEK length: %d", len(dek))
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	buf := make([]byte, config.FileStreamingChunkSize)

	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			plaintext := buf[:n]

			// random nonce per chunk
			nonce := make([]byte, gcm.NonceSize())
			if _, err := rand.Read(nonce); err != nil {
				return err
			}

			ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

			if _, err := dst.Write(nonce); err != nil {
				return err
			}

			// Add length of ciphertext
			cipherLen := len(ciphertext)
			if cipherLen > math.MaxUint32 {
				return fmt.Errorf("ciphertext too large")
			}

			buflen := make([]byte, 4)
			binary.BigEndian.PutUint32(buflen, uint32(cipherLen))
			if _, err := dst.Write(buflen); err != nil {
				return err
			}

			if _, err := dst.Write(ciphertext); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}
