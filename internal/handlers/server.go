package handlers

import "github.com/anxhukumar/hashdrop/internal/store"

// Server struct to give access of store to each handler as a method
type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store: store,
	}
}
