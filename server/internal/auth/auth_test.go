package auth

import (
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashAndCompare(t *testing.T) {

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Empty password string",
			password: "",
			wantErr:  true,
		},
		{
			name:     "Password length less than 8 chars",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "Valid password",
			password: "test1234",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashedPassword(tt.password)

			// Test if we get errors when expected
			if (err != nil) != tt.wantErr {
				t.Fatalf("HashedPassword() | error = %v, wantErr = %v", err, tt.wantErr)
			}

			// If error was expected, skip further validations
			if tt.wantErr {
				return
			}

			// Test if hash is empty
			if hash == "" {
				t.Fatalf("HashedPassword() | Received empty hash")
			}

			// Test password comare works
			isMatch, err := CheckPasswordHash(tt.password, hash)
			if err != nil {
				t.Fatalf("CheckPasswordHash() | Received unexpected error")
			}
			if !isMatch {
				t.Fatalf("CheckPasswordHash() | Received -False- when expected -True-")
			}

		})
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)

			// Test if we get errors when expected
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT() and MakeJWT() | error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Test if we get the expected userID
			if gotUserID != tt.wantUserID {
				t.Fatalf("ValidateJWT() and MakeJWT() | gotUserID = %v, want = %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    bool
	}{
		{
			name:       "Valid bearer token",
			authHeader: "Bearer valid_token_123",
			wantToken:  "valid_token_123",
			wantErr:    false,
		},
		{
			name:       "Missing Authorization header",
			authHeader: "",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "Missing Bearer prefix",
			authHeader: "Basic credentials",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "Empty token after Bearer",
			authHeader: "Bearer ",
			wantToken:  "",
			wantErr:    true,
		},
		{
			name:       "Token with extra whitespace",
			authHeader: "Bearer   token_with_spaces   ",
			wantToken:  "token_with_spaces",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			gotToken, err := GetBearerToken(headers)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken() | error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() | gotToken= %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	t.Run("Generates valid hex string", func(t *testing.T) {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Fatalf("MakeRefreshToken() unexpected error: %v", err)
		}

		// Should be 64 characters (32 bytes encoded as hex)
		if len(token) != 64 {
			t.Errorf("MakeRefreshToken() length = %d, want 64", len(token))
		}

		// Should be valid hex
		_, err = hex.DecodeString(token)
		if err != nil {
			t.Errorf("MakeRefreshToken() produced invalid hex: %v", err)
		}
	})

	t.Run("Generates unique tokens", func(t *testing.T) {
		token1, _ := MakeRefreshToken()
		token2, _ := MakeRefreshToken()

		if token1 == token2 {
			t.Error("MakeRefreshToken() generated duplicate tokens")
		}
	})
}
