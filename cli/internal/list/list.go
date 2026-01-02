package list

import (
	"time"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
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

func GetAllFiles() ([]FilesMetadata, error) {

	// Struct to receive decoded json response
	respBody := []FilesMetadata{}

	// Get access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return nil, err
	}

	// Get data
	err = api.GetJSON(config.GetAllFiles, &respBody, token, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil

}
