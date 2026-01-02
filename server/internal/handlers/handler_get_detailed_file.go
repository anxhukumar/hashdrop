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
			"no matching file found",
			errors.New("no matches found for the file id"),
			http.StatusNotFound,
		)
		return
	}

	if len(dbFileData) > 1 {
		RespondWithError(w,
			s.logger,
			"multiple files matched the given partial id",
			errors.New("multiple files matched the given partial id. please provide a longer / full id"),
			http.StatusConflict,
		)
		return
	}

	resp := FileDetailedData{
		FileName:           dbFileData[0].FileName,
		ID:                 dbFileData[0].ID,
		Status:             dbFileData[0].Status,
		PlaintextSizeBytes: dbFileData[0].PlaintextSizeBytes.Int64,
		EncryptedSizeBytes: dbFileData[0].EncryptedSizeBytes.Int64,
		S3Key:              dbFileData[0].S3Key,
		KeyManagementMode:  dbFileData[0].KeyManagementMode.String,
		PlaintextHash:      dbFileData[0].PlaintextHash.String,
	}

	RespondWithJSON(w, http.StatusOK, resp)

}
