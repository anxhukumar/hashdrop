package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// -Decrypts streaming
// -Always computes SHA-256 of plaintext
// -Optionally writes plaintext (if dsa != nil)
// Returns computed hash bytes
func DecryptAndHashStreaming(src io.Reader, dst io.Writer, dek []byte) ([]byte, error) {
	if len(dek) != config.ArgonKeyLen {
		return nil, fmt.Errorf("invalid DEK length: %d", len(dek))
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if nonceSize == 0 {
		return nil, errors.New("invalid nonce size")
	}

	const chunksize = 64 * 1024
	cipherBuf := make([]byte, chunksize+gcm.Overhead())

	if dst == nil {
		dst = io.Discard
	}

	hasher := sha256.New()

	for {
		// nonce
		nonce := make([]byte, nonceSize)
		_, err := io.ReadFull(src, nonce)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed reading nonce: %w", err)
		}

		// ciphertext
		n, err := io.ReadFull(src, cipherBuf)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("failed reading ciphertext: %w", err)
		}
		if n == 0 {
			return nil, errors.New("unexpected empty ciphertext chunk")
		}

		cipherChunk := cipherBuf[:n]

		// decrypt (AES-GCM also check integrity)
		plain, err := gcm.Open(nil, nonce, cipherChunk, nil)
		if err != nil {
			return nil, fmt.Errorf("decryption failed (wrong key or corrupted data): %w", err)
		}

		// hash plaintext
		if _, err := hasher.Write(plain); err != nil {
			return nil, fmt.Errorf("hash write failed: %w", err)
		}

		// optionally write plaintext out
		if _, err := dst.Write(plain); err != nil {
			return nil, fmt.Errorf("failed writing plaintext: %w", err)
		}

		if err == io.ErrUnexpectedEOF {
			break
		}
	}

	return hasher.Sum(nil), nil
}
