package verify

import (
	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Incoming: get file hash
type FileHash struct {
	Hash string `json:"hash"`
}

// Get file hash
func GetFileHash(fileID string) (string, error) {

	// Struct to receive decoded json response
	respBody := FileHash{}

	// Get access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return "", err
	}

	// Get data

	queryParam := map[string]string{
		"id": fileID,
	}

	err = api.GetJSON(config.GetFileHashEndpoint, &respBody, token, queryParam)
	if err != nil {
		return "", err
	}

	return respBody.Hash, nil
}
