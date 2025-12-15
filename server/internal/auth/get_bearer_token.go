package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("authorization header not found")
	}

	if !strings.HasPrefix(value, "Bearer ") {
		return "", errors.New("bearer token not found")
	}

	token := strings.TrimSpace(strings.TrimPrefix(value, "Bearer "))
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
