package encryption

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
	"github.com/anxhukumar/hashdrop/cli/internal/ui"
)

// Create Data Encryption Key for vault user and no-vault users and also generate user relevant errors
func ObtainFileEncryptionKey(noVault, verbose bool) (fileDEK []byte, fileSalt []byte, err error) {

	if !noVault {
		// If vault is being used
		var err error

		fileDEK, err = GenerateRandomDEK()
		if err != nil {
			if verbose {
				return nil, nil, fmt.Errorf("generate random key: %w", err)
			}
			return nil, nil, errors.New("error generating random key (use --verbose for details)")
		}

	} else {
		// If vault is not being used
		ui.PrintNoVaultWarning()
		fmt.Scanln() // waits until Enter is pressed to continue
		for {

			pass, err := prompt.ReadPassword("Enter passphrase: ")
			if err != nil {
				return nil, nil, err
			}

			if strings.TrimSpace(pass) == "" {
				fmt.Println("passphrase cannot be empty or whitespace")
				continue
			}

			// Check if the key length is valid
			if len(pass) < config.MinCustomEncryptionKeyLen {
				fmt.Printf("passphrase must be at least %d characters long\n", config.MinCustomEncryptionKeyLen)
				continue
			}

			confirmPass, err := prompt.ReadPassword("Confirm: ")
			if err != nil {
				return nil, nil, err
			}

			if pass != confirmPass {
				fmt.Println("passphrase do not match, try again")
				continue
			}

			fileDEK, fileSalt, err = GenerateDEKfromPassphrase(pass)
			if err != nil {
				if verbose {
					return nil, nil, fmt.Errorf("generate key from passphrase: %w", err)
				}
				return nil, nil, errors.New("error generating key from passphrase (use --verbose for details)")
			}

			break
		}
	}

	return fileDEK, fileSalt, nil
}
