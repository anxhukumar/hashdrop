package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure database connection
	dbConn, err := sql.Open("sqlite3", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}
	defer dbConn.Close()

	// Create store struct instance
	store := store.NewStore(dbConn)
	// Provide the store and config to the server
	server := handlers.NewServer(store, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("DELETE /admin/reset", server.HandlerReset)
	mux.HandleFunc("GET /api/healthz", server.HandlerReadiness)
	mux.HandleFunc("POST /api/register", server.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", server.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", server.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", server.HandlerRevokeToken)

	port := cfg.Port
	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running on port: %s\n", port)

	// Log error if server fails.
	log.Fatal(serv.ListenAndServe())
}
