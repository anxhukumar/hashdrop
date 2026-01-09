package handlers

import (
	"errors"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/database"
)

func (s *Server) HandlerResolveFileMatches(w http.ResponseWriter, r *http.Request) {
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

	dbFileData, err := s.store.Queries.CheckShortFileIDConflict(
		r.Context(),
		database.CheckShortFileIDConflictParams{
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

	resp := []FileIDConflictMatches{}
	for _, data := range dbFileData {
		resp = append(resp,
			FileIDConflictMatches{
				FileName: data.FileName,
				FileID:   data.ID,
			},
		)
	}

	RespondWithJSON(w, http.StatusOK, resp)
}
