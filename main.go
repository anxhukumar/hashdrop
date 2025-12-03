package main

import (
	"log"
	"net/http"

	"github.com/anxhukumar/hashdrop/internal/config"
	"github.com/anxhukumar/hashdrop/internal/handlers"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)

	port := cfg.Port
	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running on port: %s\n", port)

	// Logs error if server fails
	log.Fatal(serv.ListenAndServe())
}
