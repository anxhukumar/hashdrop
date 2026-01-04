package files

import (
	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Get all files of a user
func GetAllFiles() ([]FilesMetadata, error) {

	// Struct to receive decoded json response
	respBody := []FilesMetadata{}

	// Get access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return nil, err
	}

	// Get data
	err = api.GetJSON(config.GetAllFilesEndpoint, &respBody, token, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil

}
