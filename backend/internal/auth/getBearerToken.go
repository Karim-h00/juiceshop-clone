package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization Header Missin")
	}
	splitHeader := strings.SplitN(authHeader, " ", 2)
	if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
		return "", fmt.Errorf("Invalid authorization format")
	}

	return strings.TrimSpace(splitHeader[1]), nil
}
