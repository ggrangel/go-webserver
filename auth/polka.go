package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetPolkaApiKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	polkaKey := strings.TrimPrefix(authHeader, "ApiKey ")
	if polkaKey == authHeader {
		return "", fmt.Errorf("invalid Authorization header")
	}

	return polkaKey, nil
}
