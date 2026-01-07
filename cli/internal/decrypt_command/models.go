package decryptCommand

// Incoming: passphrase salt
type PassphraseSalt struct {
	Salt string `json:"salt"`
}
