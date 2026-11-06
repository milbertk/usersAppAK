package securitygen

import (
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/usersAppAK/helpers"
)

// HuserLogout clears the "access_token" HttpOnly cookie set by HuserLogin.
//
// Method: POST
// Auth: none required
func HuserLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"status":"MethodNotAllowed","message":"Use POST"}`))
		return
	}

	// Same Domain / Secure rules as the login (2026-10-31), or the browser
	// would keep the cookie.
	cookieDomain := helpers.LoginCookieDomain(r)

	clearCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   helpers.LoginCookieSecure(r),
		SameSite: http.SameSiteLaxMode,
	}

	if cookieDomain != "" {
		clearCookie.Domain = cookieDomain
	}

	http.SetCookie(w, clearCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"OK","message":"Logged out"}`))
}

func shouldUseSecureCookie(domain string) bool {
	domain = strings.TrimSpace(strings.ToLower(domain))

	if domain == "" {
		return false
	}

	if domain == "localhost" || strings.HasSuffix(domain, ".local") {
		return false
	}

	return true
}
