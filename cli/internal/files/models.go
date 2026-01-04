package files

import (
	"time"

	"github.com/google/uuid"
)

// Incoming: Get all files
type FilesMetadata struct {
	FileName           string    `json:"file_name"`
	EncryptedSizeBytes int64     `json:"encrypted_size_bytes"`
	Status             string    `json:"status"`
	KeyManagementMode  string    `json:"key_management_mode"`
	CreatedAt          time.Time `json:"created_at"`
	ID                 uuid.UUID `json:"file_id"`
}

// Incoming: Get details of one file
type FileDetailedData struct {
	FileName           string    `json:"file_name"`
	ID                 uuid.UUID `json:"file_id"`
	Status             string    `json:"status"`
	PlaintextSizeBytes int64     `json:"plaintext_size_bytes"`
	EncryptedSizeBytes int64     `json:"encrypted_size_bytes"`
	S3Key              string    `json:"s3_key"`
	KeyManagementMode  string    `json:"key_management_mode"`
	PlaintextHash      string    `json:"plaintext_hash"`
}
