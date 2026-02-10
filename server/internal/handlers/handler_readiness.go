package handlers

import "net/http"

// HandlerReadiness reports whether the API is healthy and ready to serve requests
func (s *Server) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.With("handler", "handler_readiness")

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		logger.Error("failed to write readiness response", "err", err)
		return
	}
}
