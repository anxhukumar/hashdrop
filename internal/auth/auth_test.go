package auth

import (
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
