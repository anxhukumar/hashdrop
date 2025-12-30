package encryption

// Vault struct to store users Data encryption keys
type Vault struct {
	Version int               `json:"version"`
	Entries map[string]string `json:"entries"`
}

// Vault meta data
type VaultMetaData struct {
	Version int `json:"version"`
	Argon   struct {
		Time    uint32 `json:"time"`
		Memory  uint32 `json:"memory"`
		Threads uint8  `json:"threads"`
		KeyLen  uint32 `json:"key_len"`
	} `json:"argon"`
	VaultSalt string `json:"vault_salt"`
}
