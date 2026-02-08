package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/anxhukumar/hashdrop/server/internal/ratelimit"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Initialize aws
	awsConfig, s3Client, sesClient, err := aws.InitAWS(context.Background(), cfg.S3BucketRegion)
	if err != nil {
		log.Fatalf("Failed to initialize aws: %s", err)
	}

	// Configure database connection
	dbConn, err := sql.Open("sqlite3", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}
	defer dbConn.Close()

	// WAL mode for SQLite
	_, err = dbConn.Exec(`
    PRAGMA journal_mode=WAL;
    PRAGMA synchronous=NORMAL;
`)
	if err != nil {
		log.Fatalf("failed to enable WAL mode: %s", err)
	}

	// Initialize dependencies
	store := store.NewStore(dbConn)
	server := handlers.NewServer(store, cfg, awsConfig, s3Client, sesClient)

	// Rate limiting
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiters := ratelimit.NewDefaultLimiters(ctx)
	rl := &ratelimit.Binder{
		Server:   server,
		Limiters: limiters,
	}

	// Routes
	mux := http.NewServeMux()

	mux.Handle("DELETE /admin/reset", rl.Reset(http.HandlerFunc(server.HandlerReset)))
	mux.Handle("GET /api/healthz", rl.Healthz(http.HandlerFunc(server.HandlerReadiness)))
	mux.Handle("POST /api/user/register", rl.Auth(http.HandlerFunc(server.HandlerCreateUser)))
	mux.Handle("PATCH /api/user/verify", rl.Auth(http.HandlerFunc(server.HandlerVerifyUser)))
	mux.Handle("POST /api/user/login", rl.Auth(http.HandlerFunc(server.HandlerLogin)))
	mux.Handle("POST /api/token/refresh", rl.Token(http.HandlerFunc(server.HandlerRefreshToken)))
	mux.Handle("POST /api/token/revoke", rl.Token(http.HandlerFunc(server.HandlerRevokeToken)))
	mux.Handle("DELETE /api/user", rl.Auth(http.HandlerFunc(server.HandlerDeleteUser)))
	mux.Handle("POST /api/files/presign", server.Auth(rl.Upload(http.HandlerFunc(server.HandlerGeneratePresignLink))))
	mux.Handle("POST /api/files/complete", server.Auth(rl.Upload(http.HandlerFunc(server.HandlerCompleteFileUpload))))
	mux.Handle("GET /api/files/all", server.Auth(rl.List(http.HandlerFunc(server.HandlerGetAllFiles))))
	mux.Handle("GET /api/files/resolve", server.Auth(rl.List(http.HandlerFunc(server.HandlerResolveFileMatches))))
	mux.Handle("GET /api/files", server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetDetailedFile))))
	mux.Handle("GET /api/files/salt", server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetPassphraseSalt))))
	mux.Handle("GET /api/files/hash", server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetFileHash))))
	mux.Handle("DELETE /api/files", server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerDeleteFile))))
	mux.Handle("GET /api/files/download/{userIDHash}/{fileID}", http.HandlerFunc(server.HandlerGenerateDownloadLink))

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + cfg.Port,
	}

	go func() {
		log.Printf("Server running on port: %s\n", cfg.Port)
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := serv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
