package handlers

import (
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/auth"
	"github.com/anxhukumar/hashdrop/server/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user json data
	var userIncoming UserIncoming
	if err := DecodeJson(r, &userIncoming); err != nil {
		RespondWithError(w, s.Logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get hashed password
	hashedPassword, err := auth.HashedPassword(userIncoming.Password)
	if err != nil {
		RespondWithError(w, s.Logger, "Invalid password or password too short", err, http.StatusBadRequest)
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
		RespondWithError(w, s.Logger, "Error creating new user", err, http.StatusInternalServerError)
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
		s.Logger.Println("failed to send response:", err)
		return
	}
}
