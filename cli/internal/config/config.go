package config

const (
	MinPasswordLen            = 8
	MinCustomEncryptionKeyLen = 12
	UploadFileSizeLimit       = 50

	// API
	BaseURL                    = "http://localhost:8080"
	RegisterEndpoint           = "/api/register"
	LoginEndpoint              = "/api/login"
	RefreshTokenEndpoint       = "/api/refresh"
	RevokeRefreshTokenEndpoint = "/api/revoke"
	GetPresignedLinkEndpoint   = "/api/files/presign"

	// No-vault DEK argon2 params
	ArgonTime    = 3
	ArgonMemory  = 64 * 1024
	ArgonThreads = 1
	ArgonKeyLen  = 32
)
