package config

const (
	// ================================
	// App & Storage Paths
	// ================================
	ConfigDirName         = ".hashdrop"       // directory under user home
	TokensFileName        = "tokens.json"     // authentication tokens file
	VaultFileName         = "vault.enc"       // encrypted vault file
	VaultMetadataFileName = "vault_meta.json" // vault metadata (argon params, salt, etc.)
	VaultVersion          = 1                 // vault schema version

	// ================================
	// Security / Password Policy
	// ================================
	MinPasswordLen            = 8  // account password minimum
	MinCustomEncryptionKeyLen = 12 // no-vault passphrase minimum
	MinVaultPasswordLen       = 12 // vault password minimum

	// ================================
	// Upload & File Limits
	// ================================
	UploadFileSizeLimit        = 50        // max upload file size (MB)
	MaxTimeAllowedToUploadFile = 30        // max upload time (minutes)
	FileStreamingChunkSize     = 64 * 1024 // streaming chunk size (64KB)

	// ================================
	// API Configuration
	// ================================
	BaseURL                    = "http://localhost:8080"
	RegisterEndpoint           = "/api/register"
	LoginEndpoint              = "/api/login"
	RefreshTokenEndpoint       = "/api/refresh"
	RevokeRefreshTokenEndpoint = "/api/revoke"
	GetPresignedLinkEndpoint   = "/api/files/presign"
	CompleteFileUploadEndpoint = "/api/files/complete"

	// ================================
	// Cryptography — Argon2 Parameters
	// Used for:
	//  - Vault Master Key derivation
	//  - No-vault file passphrase key derivation
	// ================================
	ArgonTime    = 3         // iterations
	ArgonMemory  = 64 * 1024 // memory cost (KB) → 64MB
	ArgonThreads = 1         // parallelism
	ArgonKeyLen  = 32        // derived key length (bytes)

	// ================================
	// DOWNLOAD URL PREFIX
	// ================================
	UrlPrefix = "https://cdn.hashdrop.dev/"
)
