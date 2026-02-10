package handlers

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/anxhukumar/hashdrop/server/internal/otp"
	"github.com/google/uuid"
)

func (s *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_create_user")

	// Get decoded incoming user json data
	var userIncoming UserIncoming
	if err := DecodeJson(r, &userIncoming); err != nil {
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

	// Get hashed password
	hashedPassword, err := auth.HashedPassword(userIncoming.Password)
	if err != nil {
		msgToDev := "user inserted password is invalid or password is too short"
		msgToClient := "invalid password"
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

	// Send data to db
	userDb := database.CreateNewUserParams{
		ID:             uuid.New(),
		Email:          userIncoming.Email,
		HashedPassword: hashedPassword,
	}
	userDbResponse, err := s.Store.Queries.CreateNewUser(r.Context(), userDb)
	if err != nil {
		msgToDev := "error creating new user in database"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Attach user_id in logger context to enhance logs
	logger = logger.With("user_id", userDbResponse.ID.String())

	// Generate and save otp in database and email it to the users email address
	err = otp.GenerateAndEmailOtp(
		r.Context(),
		userDb.ID,
		userDb.Email,
		s.Cfg.OtpHashingSecret,
		s.Store.Queries,
		s.SESClient,
	)
	if err != nil {
		msgToDev := "error generating or sending otp to users email"
		RespondWithError(
			w,
			logger,
			msgToDev,
			err,
			http.StatusInternalServerError,
		)
		return
	}

	// Marshal output and send json to user
	UserOutgoing := UserOutgoing{
		ID:        userDbResponse.ID,
		CreatedAt: userDbResponse.CreatedAt,
		UpdatedAt: userDbResponse.UpdatedAt,
		Email:     userDbResponse.Email,
	}

	if err := RespondWithJSON(w, http.StatusCreated, UserOutgoing); err != nil {
		logger.Error("failed to send response", "err", err)
		return
	}
}
