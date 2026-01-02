package upload

import "github.com/google/uuid"

// Outgoing: Sent to receive a presigned S3 POST URL
// It carries basic metadata for the upload.
type FileUploadRequest struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
}

// Incoming: Presigned url response
type PresignResponse struct {
	FileID         uuid.UUID      `json:"file_id"`
	UploadResource S3PostResponse `json:"upload_resource"`
}

type S3PostResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

// Outgoing: To the server to complete upload
type FileUploadMetadata struct {
	FileID             uuid.UUID `json:"file_id"`
	PlaintextHash      string    `json:"plaintext_hash"`
	PlaintextSizeBytes int64     `json:"plaintext_size_bytes"`
	PassphraseSalt     string    `json:"passphrase_salt"`
}

// Incoming: Success response if file upload is successfull
type FileUploadSuccessResponse struct {
	S3ObjectKey      string `json:"s3_object_key"`
	UploadedFileSize int64  `json:"uploaded_file_size"`
}
