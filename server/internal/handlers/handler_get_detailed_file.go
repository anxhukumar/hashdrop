package handlers

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
)

func (s *Server) HandlerGetDetailedFile(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_get_detailed_file")

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

	q := r.URL.Query()
	fileID := q.Get("id")

	if len(fileID) == 0 {
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
	logger = logger.With("file_id", fileID)

	dbFileData, err := s.Store.Queries.GetDetailedFile(
		r.Context(),
		database.GetDetailedFileParams{
			UserID:  userID,
			Column2: fileID + "%",
		},
	)
	if err != nil {
		msgToDev := "error fetching detailed file data from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if len(dbFileData) == 0 {
		msgToDev := "no files found matching given file id"
		msgToClient := "no files found"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
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

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}

	logger.Info("fetched detailed file information successfully")
}
