package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerGetFileHash(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_get_file_hash")

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		msgToDev := "user id missing in request context"
		RespondWithError(
			w,
			logger,
			msgToDev,
			nil,
			http.StatusInternalServerError,
		)
		return
	}

	// Attach user_id in logger context
	logger = logger.With("user_id", userID.String())

	fileIDStr := r.URL.Query().Get("id")
	if len(fileIDStr) == 0 {
		msgToDev := "file id missing in query parameter"
		msgToClient := "missing file id in query parameter"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusBadRequest,
		)
		return
	}

	// Attach file_id in logger context
	logger = logger.With("file_id", fileIDStr)

	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		msgToDev := "invalid file id format in query parameter"
		msgToClient := "invalid file id"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			err,
			http.StatusBadRequest,
		)
		return
	}

	dbFileData, err := s.Store.Queries.GetFileHash(
		r.Context(),
		database.GetFileHashParams{
			UserID: userID,
			ID:     fileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "file hash not found for given file id"
			msgToClient := "file not found"
			RespondWithWarn(
				w,
				logger,
				msgToDev,
				msgToClient,
				err,
				http.StatusNotFound,
			)
			return
		}
		msgToDev := "error fetching file hash from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	resp := FileHash{Hash: dbFileData.String}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}
}
