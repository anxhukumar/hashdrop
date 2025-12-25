package handlers

import (
	"log"

	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Server struct to give access of store and logger to each handler as a method
type Server struct {
	store    *store.Store
	logger   *log.Logger
	cfg      *config.Config
	s3Config aws.Config
	s3Client *s3.Client
}

func NewServer(store *store.Store, cfg *config.Config, s3Config aws.Config, s3Client *s3.Client) *Server {
	return &Server{
		store:    store,
		logger:   log.Default(),
		cfg:      cfg,
		s3Config: s3Config,
		s3Client: s3Client,
	}
}
