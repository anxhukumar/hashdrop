package ui

import (
	"fmt"
	"strings"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/files"
)

func ShowFile(fileData files.FileDetailedData) {

	downloadURL := fmt.Sprintf(
		"%s/%s",
		config.UrlPrefix,
		fileData.S3Key,
	)

	fmt.Println()
	fmt.Println("================= FILE DETAILS =================")
	fmt.Println()
	fmt.Printf("Name:              %s\n", fileData.FileName)
	fmt.Printf("File ID:           %s\n", fileData.ID.String()[:8])
	fmt.Printf("Status:            %s\n", fileData.Status)
	fmt.Println()

	fmt.Println("Encryption:")
	fmt.Printf("  Mode:            %s\n", fileData.KeyManagementMode)
	fmt.Printf("  Plaintext size:  %s\n", formatBytes(fileData.PlaintextSizeBytes))
	fmt.Printf("  Encrypted size:  %s\n", formatBytes(fileData.EncryptedSizeBytes))
	fmt.Println()

	fmt.Println("Integrity:")
	fmt.Println("  Hash (SHA-256):")
	fmt.Printf("    %s\n", wrapHash(fileData.PlaintextHash))
	fmt.Println()

	fmt.Println("Access:")
	fmt.Printf("  Download URL:    %s\n", downloadURL)
	fmt.Println()
	fmt.Println("------------------------------------------------")
}

func wrapHash(hash string) string {
	const width = 64
	var out string

	for i := 0; i < len(hash); i += width {
		end := i + width
		if end > len(hash) {
			end = len(hash)
		}
		out += hash[i:end] + "\n    "
	}

	return strings.TrimSpace(out)
}
