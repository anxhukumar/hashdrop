package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt file in chunks. [nonce][ciphertext+tag]
func EncryptFileStreaming(src io.Reader, dst io.Writer, dek []byte) error {

	block, err := aes.NewCipher(dek)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	buf := make([]byte, 64*1024) // 64KB chunks

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
