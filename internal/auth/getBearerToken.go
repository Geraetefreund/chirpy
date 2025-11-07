package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("missing Authorization header")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", errors.New("invalid auth scheme")
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
