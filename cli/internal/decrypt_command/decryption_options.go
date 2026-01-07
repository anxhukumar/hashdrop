package decryptCommand

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/prompt"
)

type DecryptMode int

const (
	VaultDecryptMode DecryptMode = iota
	KeyDecryptMode
)

// Options to display if the user didn't use flags to select mode
func ShowDecryptionOptions() (DecryptMode, error) {
	fmt.Println("How would you like to decrypt this file?")
	fmt.Println()
	fmt.Println("1) Use local vault")
	fmt.Println("2) Use encryption secret")
	fmt.Println()
	choice, err := prompt.ReadLine("Choose [1/2]: ")
	if err != nil {
		return 0, err
	}
	switch choice {
	case "1":
		return VaultDecryptMode, nil
	case "2":
		return KeyDecryptMode, nil
	default:
		return 0, fmt.Errorf("invalid choice: %q (expected 1 or 2)", choice)
	}
}
