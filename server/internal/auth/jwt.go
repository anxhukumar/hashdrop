package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Create JWT
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "hashdrop",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	res, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error creating signed jwt string: %w", err)
	}

	return res, nil
}

// Validate JWT
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claimsStruct := jwt.RegisteredClaims{
		Issuer: "hashdrop",
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			// validate signing method
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error validating jwt: %w", err)
	}

	// Get the userID
	idString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error getting userID string while validating JWT: %w", err)
	}

	// Parse userID string to UUID
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing UUID while validating JWT: %w", err)
	}

	return id, nil
}
