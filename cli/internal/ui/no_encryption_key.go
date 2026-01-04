package ui

import (
	"fmt"
)

func NoEncryptionKey() {
	msg := fmt.Sprintln(`
❌ Unable to find the encryption key for this file.

This file’s key is not present in your local vault.
Possible reasons:
• The file was uploaded using no-vault mode
• The vault was recreated or reset after upload
• This file was uploaded on a different machine

Without the correct key, Hashdrop cannot decrypt this file.	
	`)
	fmt.Print(msg)
}
