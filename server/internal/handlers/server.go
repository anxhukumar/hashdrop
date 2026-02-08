package handlers

import (
	"log"

	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

// Server struct to give access of store and logger to each handler as a method
type Server struct {
	Store     *store.Store
	Logger    *log.Logger
	Cfg       *config.Config
	AwsConfig aws.Config
	S3Client  *s3.Client
	SESClient *sesv2.Client
}

func NewServer(store *store.Store, cfg *config.Config, awsConfig aws.Config, s3Client *s3.Client, sesClient *sesv2.Client) *Server {
	return &Server{
		Store:     store,
		Logger:    log.Default(),
		Cfg:       cfg,
		AwsConfig: awsConfig,
		S3Client:  s3Client,
		SESClient: sesClient,
	}
}
