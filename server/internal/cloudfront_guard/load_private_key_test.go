package cloudfrontguard

import (
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	validKey := generateTempPrivateKey(t)

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "Valid private key",
			key:     validKey,
			wantErr: false,
		},
		{
			name:    "Invalid private key",
			key:     "invalid-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := loadPrivateKey(tt.key)

			if (err != nil) != tt.wantErr {
				t.Fatalf("loadPrivateKey() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && key == nil {
				t.Fatalf("expected non-nil key")
			}
		})
	}
}
