package huserlogin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/milbertk/class"
	"github.com/milbertk/security"
	"github.com/milbertk/tracking"
	bcreateusers "github.com/milbertk/usersAppAK/brules/bcreateUsers"
	databaseaccount "github.com/milbertk/usersAppAK/database/databaseAccount"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
	"github.com/milbertk/usersAppAK/helpers"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Expires string `json:"expires,omitempty"`
	UID     string `json:"uid,omitempty"`
	Email   string `json:"email,omitempty"`
}

func HuserLogin(w http.ResponseWriter, r *http.Request) {

	// El CORS ya lo maneja el router en el outer layer, no lo dupliquemos aquí.

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, LoginResponse{
			Status:  "MethodNotAllowed",
			Message: "Use POST",
		})
		return
	}

	// --- Parse Bearer
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		writeJSON(w, http.StatusUnauthorized, LoginResponse{
			Status:  "Unauthorized",
			Message: "Missing or invalid Authorization header",
		})
		return
	}
	raw := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, err := url.QueryUnescape(raw)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, LoginResponse{
			Status:  "BadRequest",
			Message: "Invalid token encoding",
		})
		return
	}
	firebaseIDToken := strings.Trim(decoded, `"`)

	// --- Validate Firebase; get claims
	tokenNew, err := security.ValidateFirebaseToken(firebaseIDToken)
	if err != nil {
		fmt.Println("ValidateFirebaseToken error:", err)
		writeJSON(w, http.StatusUnauthorized, LoginResponse{
			Status:  "Unauthorized",
			Message: "Invalid Firebase token",
		})
		return
	}
	claims := tokenNew.Claims

	uid := class.GetClaims(claims, "user_id")
	email := class.GetClaims(claims, "email")

	Ex, exist, err := security.DValidateUser(uid)
	if err != nil {
		fmt.Println("DValidateUser error:", err)
		writeJSON(w, http.StatusBadRequest, LoginResponse{
			Status:  "BadRequest",
			Message: "Error validating user: " + exist,
		})
		return
	}

	if !Ex {
		uid := tokenNew.UID
		email := class.GetClaims(claims, "email")
		displayName := class.GetClaims(claims, "name")
		photoURL := class.GetClaims(claims, "picture")
		phoneNumber := class.GetClaims(claims, "phone_number")
		providerID := tokenNew.Firebase.SignInProvider
		// Facebook with an email counts as verified (2026-10-27), see bEmailVerified.go.
		emailVerified := bcreateusers.BEffectiveEmailVerified(class.GetClaims(claims, "email_verified"), providerID, email)

		user := databaseusers.FirebaseUser{
			UID:           uid,
			LastSignIn:    time.Unix(tokenNew.AuthTime, 0).UTC().Format("2006-01-02T15:04:05Z"),
			LocalStatus:   "ACTIVE",
			CreationTime:  time.Unix(tokenNew.AuthTime, 0).UTC().Format("2006-01-02T15:04:05Z"),
			DisplayName:   displayName,
			PhoneNumber:   phoneNumber,
			PhotoURL:      photoURL,
			Active:        "true",
			ProviderID:    providerID,
			Email:         email,
			EmailVerified: emailVerified,
			Disabled:      "false",
		}

		if err := bcreateusers.BCreateUsers(user); err != nil {
			fmt.Println("BCreateUsers error:", err)
			writeJSON(w, http.StatusBadRequest, LoginResponse{
				Status:  "BadRequest",
				Message: common.InternalMessage(err, "Error creating user"),
			})
			return
		}

		account := databaseaccount.DatabaseAccountS{
			OwnerID:       uid,
			CreateDate:    time.Now(),
			StatusAccount: "ACTIVE",
		}
		if err := account.Insert(); err != nil {
			fmt.Println("account.Insert error:", err)
			writeJSON(w, http.StatusBadRequest, LoginResponse{
				Status:  "BadRequest",
				Message: common.InternalMessage(err, "Error creating the account"),
			})
			return
		}

	} else {
		userDb, err := databaseusers.SelectByUID(uid)
		if err != nil {
			fmt.Println("SelectByUID error:", err)
			writeJSON(w, http.StatusBadRequest, LoginResponse{
				Status:  "BadRequest",
				Message: common.InternalMessage(err, "Error loading user"),
			})
			return
		}

		// Facebook with an email counts as verified (2026-10-27), see bEmailVerified.go.
		effectiveVerified := bcreateusers.BEffectiveEmailVerified(class.GetClaims(claims, "email_verified"), tokenNew.Firebase.SignInProvider, email)

		// Provider and photo (2026-10-29): before, they were read from claim names that do not
		// exist ("ProviderID", "PhotoURL") and saved EMPTY, so email/password users could not
		// change their email or password. Now: the real sign-in provider and the "picture" claim,
		// keeping the saved values when the token does not bring them.
		currentProvider := strings.TrimSpace(tokenNew.Firebase.SignInProvider)
		if currentProvider == "" {
			currentProvider = userDb.ProviderID
		}
		currentPhoto := strings.TrimSpace(class.GetClaims(claims, "picture"))
		if currentPhoto == "" {
			currentPhoto = userDb.PhotoURL
		}

		if userDb.EmailVerified != effectiveVerified || userDb.ProviderID != currentProvider {
			userforUpdate := databaseusers.FirebaseUser{
				UID:           userDb.UID,
				LastSignIn:    time.Now().Format("2006-01-02 15:04:05"),
				CreationTime:  userDb.CreationTime,
				DisplayName:   userDb.DisplayName,
				PhoneNumber:   userDb.PhoneNumber,
				PhotoURL:      currentPhoto,
				Active:        userDb.Active,
				ProviderID:    currentProvider,
				Email:         userDb.Email,
				EmailVerified: effectiveVerified,
				Disabled:      userDb.Disabled,
			}
			if err := userforUpdate.Update(); err != nil {
				// Sólo logueamos, no abortamos el login por esto.
				fmt.Println("userforUpdate.Update error:", err)
			}
		}
	}

	// --- Generate local JWT
	jwt, err := security.GenerateLocalJWT(
		uid,
		email,
		class.GetClaims(claims, "name"),
		class.GetClaims(claims, "picture"),
		class.GetClaims(claims, "signInProvider"),
		class.GetClaims(claims, "role"),
		true,
	)
	if err != nil {
		fmt.Println("GenerateLocalJWT error:", err)
		writeJSON(w, http.StatusUnauthorized, LoginResponse{
			Status:  "Unauthorized",
			Message: "Could not generate token",
		})
		return
	}

	userClaims, err := security.ValidateLocalJWT(jwt)
	if err != nil {
		fmt.Println("ValidateLocalJWT error:", err)
		writeJSON(w, http.StatusBadRequest, LoginResponse{
			Status:  "BadRequest",
			Message: "Invalid User Token",
		})
		return
	}

	// --- GeoIP / tracking — TOLERANTE A FALLOS
	// IMPORTANTE: el código anterior usaba log.Fatal aquí, que mataba el
	// servidor entero si el archivo .mmdb no estaba presente.
	// Eso explicaba el "loading infinito": el handler nunca respondía
	// porque el proceso se moría a mitad de camino.
	// Ahora hacemos el tracking en best-effort: si falla, logueamos pero
	// el login sigue adelante.
	go func(token, userID, email string) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Printf("‼ tracking panic recovered: %v\n", rec)
			}
		}()

		collector, err := tracking.NewCollector("./GeoLite2-Country.mmdb")
		if err != nil {
			fmt.Println("GeoIP load error (tracking skipped):", err)
			return
		}
		defer collector.Close()

		info := collector.Extract(r)

		lt := tracking.LoginTracking{
			UserID:         userID,
			Email:          email,
			DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
			IP:             info.IP,
			Platform:       "UserManagment",
			MacAddress:     "0",
			Browser:        info.Browser,
			CountryCode:    info.CountryCode,
			GMTTime:        time.Now().Format("2006-01-02 15:04:05"),
			Lang:           info.Lang,
			Action:         "Login",
			Jsonstring:     token,
		}
		if err := lt.Insert(); err != nil {
			fmt.Println("❌ tracking insert error:", err)
		}
	}(firebaseIDToken, userClaims.UID, userClaims.Email)

	// --- Cookie settings
	/*domain := os.Getenv("AUTH_COOKIE_DOMAIN")
	if domain == "" {
		domain = "192.168.100.110"
	}*/
	expires := time.Now().Add(7 * 24 * time.Hour)
	maxAge := int((7 * 24 * time.Hour).Seconds())
	//secure := !isDevHost(domain)

	// Cookie domain (2026-10-31): the allowed domain of this request
	// (AUTH_COOKIE_DOMAIN list, else ALLOWED_DOMAINS), so the session is shared
	// by app.aleph-k.com and the APIs; localhost / LAN = host-only as before.
	// Secure when the browser used https (AUTH_COOKIE_SECURE can force it).
	cookieDomain := helpers.LoginCookieDomain(r)

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    jwt,
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   helpers.LoginCookieSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
	if cookieDomain != "" {
		cookie.Domain = cookieDomain
	}
	http.SetCookie(w, cookie)

	writeJSON(w, http.StatusOK, LoginResponse{
		Status:  "OK",
		Message: "Login successful",
		Token:   jwt,
		Expires: expires.Format(time.RFC3339),
		UID:     uid,
		Email:   email,
	})
}

func isDevHost(domain string) bool {
	return domain == "localhost" || strings.HasSuffix(domain, ".local")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func bearerToken(r *http.Request) (string, error) {
	ah := r.Header.Get("Authorization")
	if !strings.HasPrefix(ah, "Bearer ") {
		return "", errors.New("no bearer")
	}
	return strings.TrimPrefix(ah, "Bearer "), nil
}

// loginCookieSecure mirrors shouldUseSecureCookie in HuserLogout.go (2026-10-28):
// Secure only when a real cookie domain is set (not localhost / .local).
func loginCookieSecure(domain string) bool {
	domain = strings.TrimSpace(strings.ToLower(domain))
	if domain == "" || domain == "localhost" || strings.HasSuffix(domain, ".local") {
		return false
	}
	return true
}
