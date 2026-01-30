package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/anxhukumar/hashdrop/server/internal/aws"
	"github.com/anxhukumar/hashdrop/server/internal/config"
	"github.com/anxhukumar/hashdrop/server/internal/handlers"
	"github.com/anxhukumar/hashdrop/server/internal/ratelimit"
	"github.com/anxhukumar/hashdrop/server/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/time/rate"
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

	// TODO: REMOVE THIS IF WE MIGRATE FROM SQLITE TO ANY OTHER DB
	_, err = dbConn.Exec(`
    PRAGMA journal_mode=WAL;
    PRAGMA synchronous=NORMAL;
`)
	if err != nil {
		log.Fatalf("failed to enable WAL mode: %s", err)
	}

	// Create store struct instance
	store := store.NewStore(dbConn)
	// Provide the store, config and s3 dependencies
	server := handlers.NewServer(store, cfg, s3Config, s3Client)

	// Rate limiting
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	limiters := &ratelimit.Limiters{
		// ---------------------------------------------------------
		// PUBLIC / ADMIN
		// ---------------------------------------------------------

		// Reset: Extreme protection. Only required for dev.
		// 1 request every 30 seconds.
		ResetGlobalLimiter: rate.NewLimiter(rate.Every(30*time.Second), 1),

		// Healthz: Standard uptime monitoring.
		// 5 requests per second global.
		HealthzGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),

		// Auth (Register/Login/DeleteAccount):
		// Global: 10 per second.
		// IP: 1 request every 5 seconds.
		AuthGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),
		AuthIPLimiter:     ratelimit.NewKeyRateLimiter(ctx, rate.Limit(0.1), 2),

		// Token (Refresh/Revoke):
		// Frequent but lightweight.
		TokenGlobalLimiter: rate.NewLimiter(rate.Limit(20), 40),
		TokenIPLimiter:     ratelimit.NewKeyRateLimiter(ctx, rate.Limit(2), 10),

		// ---------------------------------------------------------
		// PRIVATE (S3 / DB INTENSIVE)
		// ---------------------------------------------------------

		// Upload (Presign/Complete):
		// Global: 10 uploads/sec.
		// User: 1 upload every 2 seconds.
		UploadGlobalLimiter: rate.NewLimiter(rate.Limit(5), 10),
		UploadUserLimiter:   ratelimit.NewKeyRateLimiter(ctx, rate.Limit(0.2), 3),

		// List (GetAllFiles, ResolveMatches):
		// Global: 50/sec (SQLite reads are very fast).
		// User: 5/sec.
		ListGlobalLimiter: rate.NewLimiter(rate.Limit(40), 80),
		ListUserLimiter:   ratelimit.NewKeyRateLimiter(ctx, rate.Limit(5), 10),

		// FileMeta (Detail, Salt, Hash, Delete):
		// Global: 30/sec.
		// User: 10/sec.
		FileMetaGlobalLimiter: rate.NewLimiter(rate.Limit(20), 40),
		FileMetaUserLimiter:   ratelimit.NewKeyRateLimiter(ctx, rate.Limit(2), 5),
	}

	rl := &ratelimit.Binder{
		Server:   server,
		Limiters: limiters,
	}

	mux := http.NewServeMux()

	mux.Handle(
		"DELETE /admin/reset",
		rl.Reset(http.HandlerFunc(server.HandlerReset)),
	)
	mux.Handle(
		"GET /api/healthz",
		rl.Healthz(http.HandlerFunc(server.HandlerReadiness)),
	)
	mux.Handle(
		"POST /api/user/register",
		rl.Auth(http.HandlerFunc(server.HandlerCreateUser)),
	)
	mux.Handle(
		"POST /api/user/login",
		rl.Auth(http.HandlerFunc(server.HandlerLogin)),
	)
	mux.Handle(
		"POST /api/token/refresh",
		rl.Token(http.HandlerFunc(server.HandlerRefreshToken)),
	)
	mux.Handle(
		"POST /api/token/revoke",
		rl.Token(http.HandlerFunc(server.HandlerRevokeToken)),
	)
	mux.Handle(
		"DELETE /api/user",
		rl.Auth(http.HandlerFunc(server.HandlerDeleteUser)),
	)
	mux.Handle(
		"POST /api/files/presign",
		server.Auth(rl.Upload(http.HandlerFunc(server.HandlerGeneratePresignLink))),
	)
	mux.Handle(
		"POST /api/files/complete",
		server.Auth(rl.Upload(http.HandlerFunc(server.HandlerCompleteFileUpload))),
	)
	mux.Handle(
		"GET /api/files/all",
		server.Auth(rl.List(http.HandlerFunc(server.HandlerGetAllFiles))),
	)
	mux.Handle(
		"GET /api/files/resolve",
		server.Auth(rl.List(http.HandlerFunc(server.HandlerResolveFileMatches))),
	)
	mux.Handle(
		"GET /api/files",
		server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetDetailedFile))),
	)
	mux.Handle(
		"GET /api/files/salt",
		server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetPassphraseSalt))),
	)
	mux.Handle(
		"GET /api/files/hash",
		server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerGetFileHash))),
	)
	mux.Handle(
		"DELETE /api/files",
		server.Auth(rl.FileMeta(http.HandlerFunc(server.HandlerDeleteFile))),
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
