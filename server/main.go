package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Initialize S3
	s3Config, s3Client, err := aws.InitS3(context.Background(), cfg.S3BucketRegion)
	if err != nil {
		log.Fatalf("Failed to initialize s3: %s", err)
	}

	// Configure database connection
	dbConn, err := sql.Open("sqlite3", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}
	defer dbConn.Close()

	// Create store struct instance
	store := store.NewStore(dbConn)
	// Provide the store, config and s3 dependencies
	server := handlers.NewServer(store, cfg, s3Config, s3Client)

	mux := http.NewServeMux()

	mux.HandleFunc("DELETE /admin/reset", server.HandlerReset)
	mux.HandleFunc("GET /api/healthz", server.HandlerReadiness)
	mux.HandleFunc("POST /api/register", server.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", server.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", server.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", server.HandlerRevokeToken)
	mux.Handle(
		"POST /api/files/presign",
		server.Auth(http.HandlerFunc(server.HandlerGeneratePresignLink)),
	)
	mux.Handle(
		"POST /api/files/complete",
		server.Auth(http.HandlerFunc(server.HandlerCompleteFileUpload)),
	)
	mux.Handle(
		"GET /api/files/all",
		server.Auth(http.HandlerFunc(server.HandlerGetAllFiles)),
	)
	mux.Handle(
		"GET /api/files",
		server.Auth(http.HandlerFunc(server.HandlerGetDetailedFile)),
	)
	mux.Handle(
		"GET /api/files/salt",
		server.Auth(http.HandlerFunc(server.HandlerGetPassphraseSalt)),
	)

	port := cfg.Port
	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running on port: %s\n", port)

	// Log error if server fails.
	log.Fatal(serv.ListenAndServe())
}
