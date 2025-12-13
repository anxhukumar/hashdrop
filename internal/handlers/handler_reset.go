package handlers

import (
	"errors"
	"net/http"
)

func (s *Server) HandlerReset(w http.ResponseWriter, r *http.Request) {

	// Check current platform to ensure data can't be reset in production
	if s.cfg.Platform != "dev" {
		err := errors.New("functionality accessed from wrong dev environment")
		RespondWithError(w, s.logger, "Can't access this functionality without a local development environment", err, http.StatusForbidden)
		return
	}

	// Delete all users
	err := s.store.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, s.logger, "Error while deleting all users", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
