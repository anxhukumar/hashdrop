package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerGetPassphraseSalt(w http.ResponseWriter, r *http.Request) {

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

	dbFileData, err := s.store.Queries.GetPassphraseSalt(
		r.Context(),
		database.GetPassphraseSaltParams{
			UserID: userID,
			ID:     file_id,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(
				w, s.logger,
				"Passphrase salt not available for this file",
				err,
				http.StatusNotFound,
			)
			return
		}

		RespondWithError(
			w, s.logger,
			"Error fetching file data",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	resp := passphraseSaltRes{Salt: dbFileData.String}

	RespondWithJSON(w, http.StatusOK, resp)
}
