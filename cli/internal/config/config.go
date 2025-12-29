package config

const (
	MinPasswordLen             = 8
	MinEncryptionKeyLen        = 12
	BaseURL                    = "http://localhost:8080"
	RegisterEndpoint           = "/api/register"
	LoginEndpoint              = "/api/login"
	RefreshTokenEndpoint       = "/api/refresh"
	RevokeRefreshTokenEndpoint = "/api/revoke"
	GetPresignedLinkEndpoint   = "/api/files/presign"
	UploadFileSizeLimit        = 50
)
