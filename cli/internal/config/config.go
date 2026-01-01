package config

const (
	ConfigDirName         = ".hashdrop"       // name of the tokens directory
	TokensFileName        = "tokens.json"     // name of the tokens file
	VaultFileName         = "vault.enc"       // name of the vault file
	VaultMetadataFileName = "vault_meta.json" // name of the vault file
	VaultVersion          = 1

	MinPasswordLen             = 8
	MinCustomEncryptionKeyLen  = 12
	MinVaultPasswordLen        = 12
	UploadFileSizeLimit        = 50
	MaxTimeAllowedToUploadFile = 30 // minutes

	// API
	BaseURL                    = "http://localhost:8080"
	RegisterEndpoint           = "/api/register"
	LoginEndpoint              = "/api/login"
	RefreshTokenEndpoint       = "/api/refresh"
	RevokeRefreshTokenEndpoint = "/api/revoke"
	GetPresignedLinkEndpoint   = "/api/files/presign"

	// argon2 params for Data encryption key generation
	// in case of no-vault user and Vault Master key generation
	ArgonTime    = 3
	ArgonMemory  = 64 * 1024
	ArgonThreads = 1
	ArgonKeyLen  = 32
)
