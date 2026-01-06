package ui

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/anxhukumar/hashdrop/cli/internal/files"
)

func ShowFile(fileData files.FileDetailedData, encryptionKey string) {

	if encryptionKey == "" {
		encryptionKey = "hidden (use --reveal-key)"
	}

	downloadURL := fmt.Sprintf(
		"%s%s",
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
	fmt.Printf("  Key:             %s\n", encryptionKey)
	fmt.Printf("  Mode:            %s\n", fileData.KeyManagementMode)
	fmt.Printf("  Plaintext size:  %s\n", formatBytes(fileData.PlaintextSizeBytes))
	fmt.Printf("  Encrypted size:  %s\n", formatBytes(fileData.EncryptedSizeBytes))
	fmt.Println()

	fmt.Println("Integrity:")
	fmt.Printf("  Hash (SHA-256): %s\n", fileData.PlaintextHash)
	fmt.Println()

	fmt.Println("Access:")
	fmt.Printf("  Download URL:    %s\n", downloadURL)
	fmt.Println()
	fmt.Println("------------------------------------------------")
}
