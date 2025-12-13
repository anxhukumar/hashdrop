package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/anxhukumar/hashdrop/internal/config"
	"github.com/anxhukumar/hashdrop/internal/handlers"
	"github.com/anxhukumar/hashdrop/internal/store"
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

	mux.HandleFunc("POST /admin/reset", server.HandlerReset)
	mux.HandleFunc("GET /api/healthz", server.HandlerReadiness)
	mux.HandleFunc("POST /api/register", server.HandlerCreateUser)

	port := cfg.Port
	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running on port: %s\n", port)

	// Log error if server fails.
	log.Fatal(serv.ListenAndServe())
}
