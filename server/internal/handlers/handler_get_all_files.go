package handlers

import (
	"net/http"
)

func (s *Server) HandlerGetAllFiles(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_get_all_files")

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

	// Attach user_id in logger context to enhance logs
	logger = logger.With("user_id", userID.String())

	dbFileData, err := s.Store.Queries.GetAllFilesOfUser(r.Context(), userID)
	if err != nil {
		msgToDev := "error while fetching all files of user from database"
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
		msgToDev := "no files found for user"
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

	resp := []FilesMetadata{}
	for _, data := range dbFileData {
		resp = append(resp,
			FilesMetadata{
				FileName:           data.FileName,
				EncryptedSizeBytes: data.EncryptedSizeBytes.Int64,
				Status:             data.Status,
				KeyManagementMode:  data.KeyManagementMode.String,
				CreatedAt:          data.CreatedAt,
				ID:                 data.ID,
			},
		)
	}

	if err := RespondWithJSON(w, http.StatusOK, resp); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}

	logger.Info("fetched all files of user successfully")
}
