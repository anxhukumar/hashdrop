package decryptCommand

import (
	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

// Get salt of file
func GetSalt(fileID string) (string, error) {

	// Struct to receive decoded json response
	respBody := PassphraseSalt{}

	// Get access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return "", err
	}

	// Get data

	queryParam := map[string]string{
		"id": fileID,
	}

	err = api.GetJSON(config.GetSaltEndpoint, &respBody, token, queryParam)
	if err != nil {
		return "", err
	}

	return respBody.Salt, nil

}
