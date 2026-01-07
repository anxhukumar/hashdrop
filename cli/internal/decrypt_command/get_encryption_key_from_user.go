package decryptCommand

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/encryption"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
)

// Error handling specially if salt is not found
func IsNotFound(err error) bool {
	return errors.Is(err, api.ErrNotFound)
}

func GetEncryptionKeyFromUser(fileID string) ([]byte, error) {

	for {

		pass, err := prompt.ReadPassword("Enter encryption secret: ")
		if err != nil {
			return nil, err
		}

		if strings.TrimSpace(pass) == "" {
			fmt.Println("secret cannot be empty or whitespace")
			continue
		}

		// Get salt from backend and check if we have to process passphrase or key directly
		salt, err := GetSalt(fileID)
		if err == nil {
			// Nil ensures passphrase mode is used

			// Parse string salt to bytes
			decodedSalt, err := base64.StdEncoding.DecodeString(salt)
			if err != nil {
				return nil, fmt.Errorf("invalid salt encoding: %w", err)
			}

			fileDEK := encryption.DeriveDEK(pass, decodedSalt)

			return fileDEK, nil
		} else if IsNotFound(err) {
			fileDEK, err := base64.StdEncoding.DecodeString(pass)
			if err != nil {
				return nil, errors.New("invalid encryption key")
			}

			return fileDEK, nil
		} else {
			return nil, err
		}

	}
}
