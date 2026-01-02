package ui

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/list"
)

type filesDataString struct {
	fileName           string
	encryptedSizeBytes string
	status             string
	keyManagementMode  string
	createdAt          string
	shortFileId        string
}

func ListFiles(filesData []list.FilesMetadata) {

	fileDataString := []filesDataString{}

	for _, d := range filesData {

		fileDataString = append(
			fileDataString,
			filesDataString{
				fileName:           d.FileName,
				encryptedSizeBytes: formatBytes(d.EncryptedSizeBytes),
				status:             d.Status,
				keyManagementMode:  d.KeyManagementMode,
				createdAt:          d.CreatedAt.Format("2006-01-02"),
				shortFileId:        d.ID.String()[:8],
			},
		)

	}

	msg := `
================= NO VAULT MODE ENABLED =================

You have chosen to disable the local key vault.

In this mode:
• You must manually provide and remember your encryption passphrase
• The passphrase is never stored by HashDrop or backed up anywhere
• Without your passphrase, your files cannot be decrypted
• Loss of the passphrase results in permanent and irreversible data loss
• Your passphrase must be at least 12 characters long

This option is intended for users who explicitly want to self-manage encryption secrets.

If you wish to proceed, press Enter.
If you do not wish to proceed, press Ctrl+C to cancel.
---------------------------------------------------------
`
	fmt.Print(msg)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
