package handlers

import (
	"log"
	"net/http"
)

// Helper function to log error and also give appropriate response to the client
func RespondWithError(w http.ResponseWriter, logger *log.Logger, msg string, err error, code int) {
	if err != nil {
		// Logs error message for the developer
		logger.Println(err)
	}
	// Sends error message to the client
	http.Error(w, msg, code)
}
