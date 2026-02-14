package cloudfrontguard

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	validPath := writeTempPrivateKey(t)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid private key",
			path:    validPath,
			wantErr: false,
		},
		{
			name:    "File does not exist",
			path:    "/no/such/file.pem",
			wantErr: true,
		},
		{
			name: "Invalid PEM file",
			path: func() string {
				dir := t.TempDir()
				path := filepath.Join(dir, "bad.pem")
				if err := os.WriteFile(path, []byte("not a pem"), 0600); err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}

				return path
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := loadPrivateKey(tt.path)

			if (err != nil) != tt.wantErr {
				t.Fatalf("loadPrivateKey() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && key == nil {
				t.Fatalf("expected non-nil key")
			}
		})
	}
}
