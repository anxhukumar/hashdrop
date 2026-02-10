package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
)

func (s *Server) HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_delete_user")

	// Get decoded incoming user login data
	var userLoginIncoming UserLoginIncoming
	if err := DecodeJson(r, &userLoginIncoming); err != nil {
		msgToDev := "user posted invalid json data"
		msgToClient := "invalid JSON payload"
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

	// Check if user is registered & verified to get account details
	userData, err := s.Store.Queries.GetVerifiedUserByEmail(r.Context(), userLoginIncoming.Email)
	if err != nil {
		msgToDev := "invalid username or password while deleting user"
		msgToClient := "invalid username or password"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			err,
			http.StatusUnauthorized,
		)
		return
	}

	// Attach user_id in logger context to enhance logs
	logger = logger.With("user_id", userData.ID.String())

	// Check if password is correct
	isMatch, err := auth.CheckPasswordHash(userLoginIncoming.Password, userData.HashedPassword)
	if err != nil {
		msgToDev := "error while verifying user password hash"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if !isMatch {
		msgToDev := "password mismatch while deleting user"
		msgToClient := "invalid username or password"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			nil,
			http.StatusUnauthorized,
		)
		return
	}

	// Get s3 key of user
	s3Key, err := s.Store.Queries.GetAnyS3KeyOfUser(r.Context(), userData.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			msgToDev := "error getting s3 key of user from database"
			RespondWithError(
				w,
				logger,
				msgToDev,
				err,
				http.StatusInternalServerError,
			)
			return
		}
	} else {
		// Fetch prefix
		parts := strings.SplitN(s3Key, "/", 2)
		if len(parts) < 2 {
			msgToDev := "invalid s3 key format stored for user"
			RespondWithError(
				w,
				logger,
				msgToDev,
				nil,
				http.StatusInternalServerError,
			)
			return
		}
		s3UserPrefix := parts[0] + "/"

		// Delete all objects of user
		err = DeleteAllUserS3Obj(r.Context(), s.S3Client, s.Cfg.S3Bucket, s3UserPrefix)
		if err != nil {
			msgToDev := "error deleting all objects of user from s3"
			RespondWithError(
				w,
				logger,
				msgToDev,
				err,
				http.StatusInternalServerError,
			)
			return
		}
	}

	// Delete all users data from database
	err = s.Store.Queries.DeleteUserById(r.Context(), userData.ID)
	if err != nil {
		msgToDev := "error deleting user data from database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("account deleted")
}
