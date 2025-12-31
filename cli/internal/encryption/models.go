package encryption

// Vault struct to store users Data encryption keys
type Vault struct {
	Version int               `json:"version"`
	Entries map[string]string `json:"entries"`
}

// Vault key meta data
type VaultKeyMetadata struct {
	Version int         `json:"version"`
	Argon   ArgonParams `json:"argon"`
	Salt    []byte      `json:"vault_salt"`
}
type ArgonParams struct {
	Time    uint32 `json:"time"`
	Memory  uint32 `json:"memory"`
	Threads uint8  `json:"threads"`
	KeyLen  uint32 `json:"key_len"`
}
