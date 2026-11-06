package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/milbertk/security"
)

type wrappedLoginResponse struct {
	Error  bool `json:"error"`
	Status int  `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

// ExtractBearerToken has returned the real JWT from Authorization header.
// It has supported:
//   - Bearer <jwt>
//   - Bearer "<jwt>"
//   - Bearer { "data": { "token": "<jwt>" } }   (wrapped response)
func ExtractBearerToken(r *http.Request) (string, error) {
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

	// Case 1: already JWT (3 segments => 2 dots)
	if strings.Count(decoded, ".") == 2 {
		return decoded, nil
	}

	// Case 2: wrapped JSON that contains data.token
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

// RequireLocalJWT has validated your local JWT and has returned UID.
func RequireLocalJWT(r *http.Request) (string, error) {
	// The request and the token are never printed (2026-10-25, security S11).
	token, err := ExtractBearerToken(r)
	if err != nil {
		return "", err
	}
	parsed, err := security.ValidateLocalJWT(token)
	if err != nil {
		return "", fmt.Errorf("invalid token or error decoding")
	}

	if parsed.UID == "" {
		return "", fmt.Errorf("invalid token (missing uid)")
	}

	return parsed.UID, nil
}

// OPTIONAL: If later you have needed Firebase ID tokens too, you can use this:
// func RequireFirebaseIDToken(r *http.Request) (string, error) {
// 	token, err := ExtractBearerToken(r)
// 	if err != nil {
// 		return "", err
// 	}
// 	parsed, err := security.ValidateFirebaseToken(token)
// 	if err != nil {
// 		return "", fmt.Errorf("invalid Firebase ID token")
// 	}
// 	if parsed.UID == "" {
// 		return "", fmt.Errorf("invalid token (missing uid)")
// 	}
// 	return parsed.UID, nil
// }
