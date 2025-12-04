package handlers

import (
	"log"

	"github.com/anxhukumar/hashdrop/internal/store"
)

// Server struct to give access of store and logger to each handler as a method
type Server struct {
	store  *store.Store
	logger *log.Logger
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store:  store,
		logger: log.Default(),
	}
}
