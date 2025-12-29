package upload

import (
	"fmt"

	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
)

func GetPresignedLink(fileName, mimeType string, respBody *PresignResponse) error {

	reqBody := FileUploadRequest{
		FileName: fileName,
		MimeType: mimeType,
	}

	// Post data
	err := api.PostJSON(config.GetPresignedLinkEndpoint, reqBody, respBody)
	if err != nil {
		return fmt.Errorf("get presigned link: %w", err)
	}

	return nil
}
