package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
)

func (s *Server) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user login data
	var userLoginIncoming UserLoginIncoming
	if err := DecodeJson(r, &userLoginIncoming); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Check if user is registered and get account details
	userData, err := s.store.Queries.GetUserByEmail(r.Context(), userLoginIncoming.Email)
	if err != nil {
		RespondWithError(w, s.logger, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

	// Check if password is correct
	isMatch, err := auth.CheckPasswordHash(userLoginIncoming.Password, userData.HashedPassword)
	if err != nil {
		RespondWithError(w, s.logger, "Error verifying password", err, http.StatusInternalServerError)
		return
	}

	if !isMatch {
		RespondWithError(w, s.logger, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	// Get s3 key of user
	s3Key, err := s.store.Queries.GetAnyS3KeyOfUser(r.Context(), userData.ID)
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
	err = s.store.Queries.DeleteUserById(r.Context(), userData.ID)
	if err != nil {
		RespondWithError(w, s.logger, "Error deleting users data from database", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
