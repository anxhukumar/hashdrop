package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerGetPassphraseSalt(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_get_passphrase_salt")

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

	dbFileData, err := s.Store.Queries.GetPassphraseSalt(
		r.Context(),
		database.GetPassphraseSaltParams{
			UserID: userID,
			ID:     fileID,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			msgToDev := "passphrase salt not available for this file"
			msgToClient := "passphrase salt not found"
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

		msgToDev := "error fetching passphrase salt from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	resp := PassphraseSaltRes{Salt: dbFileData.String}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}

	logger.Info("fetched passphrase salt successfully")
}
