package handlers

import (
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
)

func (s *Server) HandlerGetDetailedFile(w http.ResponseWriter, r *http.Request) {

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	file_id := q.Get("id")

	if len(file_id) == 0 {
		RespondWithError(w,
			s.logger,
			"Missing file id in query parameter",
			errors.New("file id missing in query"),
			http.StatusBadRequest)
		return
	}

	dbFileData, err := s.store.Queries.GetDetailedFile(
		r.Context(),
		database.GetDetailedFileParams{
			UserID:  userID,
			Column2: file_id + "%",
		},
	)
	if err != nil {
		RespondWithError(w, s.logger, "Error fetching file data", err, http.StatusInternalServerError)
		return
	}

	if len(dbFileData) == 0 {
		RespondWithError(
			w,
			s.logger,
			"no files found",
			errors.New("no matches found for the file id"),
			http.StatusNotFound,
		)
		return
	}

	resp := []FileDetailedData{}
	for _, data := range dbFileData {
		resp = append(resp,
			FileDetailedData{
				FileName:           data.FileName,
				ID:                 data.ID,
				Status:             data.Status,
				PlaintextSizeBytes: data.PlaintextSizeBytes.Int64,
				EncryptedSizeBytes: data.EncryptedSizeBytes.Int64,
				S3Key:              data.S3Key,
				KeyManagementMode:  data.KeyManagementMode.String,
				PlaintextHash:      data.PlaintextHash.String,
			},
		)
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
