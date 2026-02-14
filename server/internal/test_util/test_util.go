package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const sqlSchemaDir = "../sql/schema"

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create temp file for sqlite db
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	runMigrations(t, db)

	return db
}

func runMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	files, err := os.ReadDir(sqlSchemaDir)
	if err != nil {
		t.Fatalf("failed to read migrations dir: %v", err)
	}

	var names []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}

	sort.Strings(names)

	for _, name := range names {
		path := filepath.Join(sqlSchemaDir, name)

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read migration %s: %v", name, err)
		}

		sqlText := string(sqlBytes)

		// Extract only the -- +goose Up section
		upSQL := extractGooseUp(t, name, sqlText)

		statements := strings.Split(upSQL, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := db.Exec(stmt); err != nil {
				t.Fatalf("failed to run migration %s statement [%s]: %v", name, stmt, err)
			}
		}
	}
}

func extractGooseUp(t *testing.T, filename, sqlText string) string {
	t.Helper()

	upMarker := "-- +goose Up"
	downMarker := "-- +goose Down"

	upIdx := strings.Index(sqlText, upMarker)
	if upIdx == -1 {
		t.Fatalf("migration %s has no -- +goose Up section", filename)
	}

	start := upIdx + len(upMarker)

	downIdx := strings.Index(sqlText[start:], downMarker)
	if downIdx == -1 {
		// No Down section, take everything after Up
		return sqlText[start:]
	}

	return sqlText[start : start+downIdx]
}
