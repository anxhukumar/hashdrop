package decryptCommand

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/verify"
)

func VerifyHash(fileID string, verbose bool, derivedPlaintextHash []byte) error {
	originalHashStr, err := verify.GetFileHash(fileID)
	if err != nil {
		if verbose {
			return fmt.Errorf("get original hash: %w", err)
		}
		return errors.New("error while getting original file hash")
	}

	originalHashBytes, err := hex.DecodeString(originalHashStr)
	if err != nil {
		return fmt.Errorf("invalid hash format from server: %w", err)
	}

	// Verify hash
	if !bytes.Equal(derivedPlaintextHash, originalHashBytes) {
		fmt.Println("❌ Hash mismatch — file may be corrupted")
		return nil
	}

	fmt.Println("✅ Hash verified — file is intact")
	return nil
}
