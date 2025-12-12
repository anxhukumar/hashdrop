package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anxhukumar/hashdrop/internal/auth"
	"github.com/anxhukumar/hashdrop/internal/database"
	"github.com/google/uuid"
)

func (s *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// Get decoded incoming user json data
	var userIncoming UserIncoming
	if err := DecodeJson(r, &userIncoming); err != nil {
		RespondWithError(w, s.logger, "Invalid JSON payload", err, http.StatusBadRequest)
		return
	}

	// Get hashed password
	hashedPassword, err := auth.HashedPassword(userIncoming.Password)
	if err != nil {
		RespondWithError(w, s.logger, "Invalid password or password too short", err, http.StatusBadRequest)
		return
	}

	// Send data to db
	userDb := database.CreateNewUserParams{
		ID:             uuid.New(),
		Email:          userIncoming.Email,
		HashedPassword: hashedPassword,
	}
	userDbResponse, err := s.store.Queries.CreateNewUser(r.Context(), userDb)
	if err != nil {
		RespondWithError(w, s.logger, "Error while creating new user", err, http.StatusInternalServerError)
		return
	}

	// Decode output and send json to user
	UserOutgoing := UserOutgoing{
		ID:        userDbResponse.ID,
		CreatedAt: userDbResponse.CreatedAt,
		UpdatedAt: userDbResponse.UpdatedAt,
		Email:     userDbResponse.Email,
	}

	res, err := json.Marshal(UserOutgoing)
	if err != nil {
		RespondWithError(w, s.logger, "Error while sending created user data", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
