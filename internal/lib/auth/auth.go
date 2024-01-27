package auth

import (
	"net/http"
	"strings"

	"github.com/c-mierez/rss-aggregator/internal/lib/err"
)

// GetAPIKey extracts the API key from the headers of an HTTP Request
// Example:
// Authorization: ApiKey {API_KEY}
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", err.ErrNoAuthenticationFound
	}

	auth := strings.Split(authHeader, " ")

	// Check if the header is malformed
	if len(auth) != 2 || auth[0] != "ApiKey" {
		return "", err.ErrMalformedHeader
	}

	return auth[1], nil

}
