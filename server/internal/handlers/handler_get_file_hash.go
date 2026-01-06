package handlers

import (
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerGetFileHash(w http.ResponseWriter, r *http.Request) {

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	fileIdStr := r.URL.Query().Get("id")
	if len(fileIdStr) == 0 {
		RespondWithError(w,
			s.logger,
			"Missing file id in query parameter",
			errors.New("file id missing in query"),
			http.StatusBadRequest)
		return
	}

	file_id, err := uuid.Parse(fileIdStr)
	if err != nil {
		RespondWithError(w, s.logger, "invalid file id", err, http.StatusBadRequest)
		return
	}

	dbFileData, err := s.store.Queries.GetFileHash(r.Context(), database.GetFileHashParams{UserID: userID, ID: file_id})
	if err != nil {
		RespondWithError(w, s.logger, "Error fetching file hash", err, http.StatusInternalServerError)
		return
	}

	resp := FileHash{Hash: dbFileData.String}

	RespondWithJSON(w, http.StatusOK, resp)
}
