package upload

import "github.com/google/uuid"

// Outgoing: Sent to receive a presigned S3 POST URL
// It carries basic metadata for the upload.
type FileUploadRequest struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
}

// Incoming: presigned url response
type PresignResponse struct {
	FileID         uuid.UUID      `json:"file_id"`
	UploadResource S3PostResponse `json:"upload_resource"`
}

type S3PostResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
