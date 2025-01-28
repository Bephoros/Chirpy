package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidAuthorizationHeader = errors.New("invalid Authorization header format")

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidAuthorizationHeader
	}

	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidAuthorizationHeader
	}

	apiKey := strings.TrimPrefix(authHeader, prefix)
	if apiKey == "" {
		return "", ErrInvalidAuthorizationHeader
	}

	return apiKey, nil
}
