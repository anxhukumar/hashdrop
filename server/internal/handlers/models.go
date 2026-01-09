package handlers

import (
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/google/uuid"
)

// Incoming: struct to receive from the user
type UserIncoming struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Outgoing: User struct to send a response after creation
type UserOutgoing struct {
	ID        any       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// Incoming: Login struct to receive from the user
type UserLoginIncoming struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Outgoing: struct to send the user once they are logged in
type UserLoginOutgoing struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Incoming: refresh token
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// Outgoing: new access token
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Incoming: sent by the client before we generate
// a presigned S3 POST URL. It carries basic metadata for the upload.
type FileUploadRequest struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
}

// Outgoing: presigned url response to the client
type PresignResponse struct {
	FileID         uuid.UUID         `json:"file_id"`
	UploadResource aws.S3PutResponse `json:"upload_resource"`
}

// Incoming: send by the client after successful file upload
type FileUploadMetadata struct {
	FileID             uuid.UUID `json:"file_id"`
	PlaintextHash      string    `json:"plaintext_hash"`
	PlaintextSizeBytes int64     `json:"plaintext_size_bytes"`
	PassphraseSalt     string    `json:"passphrase_salt"`
}

// Outgoing: Send status if file upload is successfull
type FileUploadSuccessResponse struct {
	S3ObjectKey      string `json:"s3_object_key"`
	UploadedFileSize int64  `json:"uploaded_file_size"`
}

// Outgoing: Send all files of a user
type FilesMetadata struct {
	FileName           string    `json:"file_name"`
	EncryptedSizeBytes int64     `json:"encrypted_size_bytes"`
	Status             string    `json:"status"`
	KeyManagementMode  string    `json:"key_management_mode"`
	CreatedAt          time.Time `json:"created_at"`
	ID                 uuid.UUID `json:"file_id"`
}

// Outgoing: Send deatils of one file
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

// Outgoing: Send passphrase salt
type PassphraseSaltRes struct {
	Salt string `json:"salt"`
}

// Outgoing: send file hash
type FileHash struct {
	Hash string `json:"hash"`
}

// Outgoing: send multiple file matches to resolve file id conflict
type FileIDConflictMatches struct {
	FileName string    `json:"file_name"`
	FileID   uuid.UUID `json:"file_id"`
}
