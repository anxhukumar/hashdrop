package upload

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func CompleteFileUpload(reqBody FileUploadMetadata) (FileUploadSuccessResponse, error) {

	// Struct to receive response
	respBody := FileUploadSuccessResponse{}

	// Fetch access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return FileUploadSuccessResponse{}, err
	}

	// Post data
	err = api.PostJSON(config.CompleteFileUploadEndpoint, reqBody, &respBody, token)
	if err != nil {
		return FileUploadSuccessResponse{}, fmt.Errorf("complete file upload: %w", err)
	}

	return respBody, nil

}
