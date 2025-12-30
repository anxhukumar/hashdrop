package encryption

// Vault struct to store users Data encryption keys
type Vault struct {
	Version int `json:"version"`
	Argon   struct {
		Time    uint32 `json:"time"`
		Memory  uint32 `json:"memory"`
		Threads uint8  `json:"threads"`
	} `json:"argon"`
	VaultSalt string            `json:"vault_salt"`
	Entries   map[string]string `json:"entries"`
}
