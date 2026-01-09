package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

func (s *Server) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {

	// Get userID from context
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, s.logger, "Internal server error", errors.New("user id missing in context"), http.StatusInternalServerError)
		return
	}

	// Get s3 key of user
	s3Key, err := s.store.Queries.GetAnyS3KeyOfUser(r.Context(), userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, s.logger, "Error getting s3 key of user", err, http.StatusInternalServerError)
			return
		}
	} else {
		// Fetch prefix
		parts := strings.SplitN(s3Key, "/", 2)
		if len(parts) < 2 {
			RespondWithError(
				w, s.logger,
				"Invalid S3 key format",
				errors.New("no '/' found in s3 key"),
				http.StatusInternalServerError,
			)
			return
		}
		s3UserPrefix := parts[0] + "/"

		// Delete all objects of user
		err = DeleteAllUserS3Obj(r.Context(), s.s3Client, s.cfg.S3Bucket, s3UserPrefix)
		if err != nil {
			RespondWithError(w, s.logger, "Error deleting all objects of user from s3", err, http.StatusInternalServerError)
			return
		}
	}

	// Delete all users data from database
	err = s.store.Queries.DeleteUserById(r.Context(), userID)
	if err != nil {
		RespondWithError(w, s.logger, "Error deleting users data from database", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
