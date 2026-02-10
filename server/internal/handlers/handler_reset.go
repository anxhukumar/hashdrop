package handlers

import (
	"errors"
	"net/http"
)

func (s *Server) HandlerReset(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_reset")

	// Check current platform to ensure data can't be reset in production
	if s.Cfg.Platform != "dev" {
		err := errors.New("functionality accessed from wrong dev environment")
		msgToDev := "reset endpoint accessed outside dev environment"
		msgToClient := "can't access this functionality without a local development environment"
		RespondWithWarn(
			w,
			logger,
			msgToDev,
			msgToClient,
			err,
			http.StatusForbidden,
		)
		return
	}

	// Delete all users
	err := s.Store.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		msgToDev := "error while deleting all users from database"
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
	logger.Info("all users deleted successfully in dev reset")
}
