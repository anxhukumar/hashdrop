package upload

import (
	"errors"
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/auth"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func CompleteFileUpload(reqBody FileUploadMetadata) error {

	// Struct to receive response
	respBody := FileUploadStatus{}

	// Fetch access token
	token, err := auth.EnsureAccessToken()
	if err != nil {
		return err
	}

	// Post data
	err = api.PostJSON(config.CompleteFileUploadEndpoint, reqBody, &respBody, token)
	if err != nil {
		return fmt.Errorf("complete file upload: %w", err)
	}

	if !respBody.Successful {
		return errors.New("server rejected upload (likely file too large)")
	}

	return nil

}
