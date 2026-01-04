package files

import (
	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Get details of one specific file
func GetDetailedFile(fileID string) ([]FileDetailedData, error) {

	// Struct to receive decoded json response
	respBody := []FileDetailedData{}

	// Get access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return nil, err
	}

	// Get data
	queryParam := map[string]string{
		"id": fileID,
	}

	err = api.GetJSON(config.GetDetailedFileEndpoint, &respBody, token, queryParam)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
