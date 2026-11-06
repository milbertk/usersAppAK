package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type wrappedLoginResponseOnlyS struct {
	Error  bool `json:"error"`
	Status int  `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

// ExtractTokenOnly has returned ONLY the JWT string from Authorization header.
// It has supported:
//
//   - Bearer <jwt>
//   - Bearer "<jwt>"
//   - Bearer { "data": { "token": "<jwt>" } }
func ExtractTokenOnly(r *http.Request) (string, error) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("missing or invalid Authorization header")
	}

	raw := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, _ := url.QueryUnescape(raw)
	decoded = strings.TrimSpace(strings.Trim(decoded, `"`))

	if decoded == "" {
		return "", fmt.Errorf("empty token")
	}

	// Case 1: already JWT
	if strings.Count(decoded, ".") == 2 {
		return decoded, nil
	}

	// Case 2: wrapped JSON
	if strings.HasPrefix(decoded, "{") {
		var w wrappedLoginResponse
		if err := json.Unmarshal([]byte(decoded), &w); err != nil {
			return "", fmt.Errorf("invalid wrapped token JSON")
		}

		if w.Data.Token == "" {
			return "", fmt.Errorf("token not found inside wrapper")
		}

		return strings.TrimSpace(w.Data.Token), nil
	}

	return "", fmt.Errorf("unknown token format")
}
