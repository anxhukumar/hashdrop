package handlers

import (
	"log"

	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/store"
)

// Server struct to give access of store and logger to each handler as a method
type Server struct {
	store  *store.Store
	logger *log.Logger
	cfg    *config.Config
}

func NewServer(store *store.Store, cfg *config.Config) *Server {
	return &Server{
		store:  store,
		logger: log.Default(),
		cfg:    cfg,
	}
}
