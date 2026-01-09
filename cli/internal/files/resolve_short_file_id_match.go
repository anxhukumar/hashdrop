package files

import (
	"github.com/anxhukumar/hashdrop/cli/internal/api"
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"github.com/google/uuid"
)

type FileIDConflictMatches struct {
	FileName string    `json:"file_name"`
	FileID   uuid.UUID `json:"file_id"`
}

// Checks if the short fileId provided by the client, matches multiple files in the database
func CheckMultipleShortFileIDMatch(fileID string, queryParam map[string]string, token string) ([]FileIDConflictMatches, error) {

	resp := []FileIDConflictMatches{}

	err := api.GetJSON(config.ResolveFileIDConflictEndpoint, &resp, token, queryParam)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
