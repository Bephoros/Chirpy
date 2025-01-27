package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization header is not a bearer token")
	}

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		return "", errors.New("bearer token is empty")
	}

	return token, nil
}
