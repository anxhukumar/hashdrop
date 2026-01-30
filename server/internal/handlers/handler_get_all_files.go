package handlers

import (
	"errors"
	"net/http"
)

func (s *Server) HandlerGetAllFiles(w http.ResponseWriter, r *http.Request) {

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.Logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	dbFileData, err := s.Store.Queries.GetAllFilesOfUser(r.Context(), userID)
	if err != nil {
		RespondWithError(w, s.Logger, "Error while fetching user data", err, http.StatusInternalServerError)
		return
	}

	if len(dbFileData) == 0 {
		RespondWithError(
			w,
			s.Logger,
			"no matching file found",
			errors.New("no matches found for the file id"),
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

	RespondWithJSON(w, http.StatusOK, resp)
}
