package ui

import (
	"fmt"
)

func PrintNoVaultWarning() {
	msg := `
================= NO VAULT MODE ENABLED =================

You have chosen to disable the local key vault.

In this mode:
• You must manually provide and remember your encryption passphrase
• The passphrase is never stored by HashDrop or backed up anywhere
• Without your passphrase, your files cannot be decrypted
• Loss of the passphrase results in permanent and irreversible data loss

This option is intended for users who explicitly want to self-manage encryption secrets.

If you wish to proceed, press Enter.
If you do not wish to proceed, press Ctrl+C to cancel.
---------------------------------------------------------
`
	fmt.Print(msg)
}
