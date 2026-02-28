package ui

import "fmt"

func PrintVaultCreationInfo() {
	msg := `
================= VAULT SETUP =================

A secure local vault will now be created to store your file encryption keys.

• This vault is protected by a password
• This is NOT your account password
• You must remember this vault password to unlock your files
• If you forget it, your vault cannot be opened and your files cannot be decrypted
• Your vault password must be at least 12 characters long

Use a strong passphrase you can reliably remember,
or store it securely in a trusted password manager.

Press Enter to continue.
Press Ctrl+C to cancel.
------------------------------------------------
`
	fmt.Print(msg)
}
