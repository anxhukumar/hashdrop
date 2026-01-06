package decrypt

// Incoming: passphrase salt
type PassphraseSalt struct {
	Salt string `json:"salt"`
}
