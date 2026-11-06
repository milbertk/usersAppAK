package husercreate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/milbertk/class"
	"github.com/milbertk/security"
	"github.com/milbertk/tracking"
	bcreateusers "github.com/milbertk/usersAppAK/brules/bcreateUsers"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HuserCreate(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	// The Authorization header (Bearer token) is NEVER written to the logs
	// (2026-10-25, security S11).

	if !strings.HasPrefix(authHeader, "Bearer ") {
		sendJSONResponse(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, _ := url.QueryUnescape(rawToken)
	token := strings.Trim(decoded, `"`)

	tokenNew, err := security.ValidateFirebaseToken(token)
	if err != nil {
		fmt.Println("❌ Firebase token validation failed:", err)
		sendJSONResponse(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	claims := tokenNew.Claims
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

	trackingPayload := map[string]interface{}{"provider": providerID, "emailVerified": emailVerified}

	if err := bcreateusers.BCreateUsers(user); err != nil {
		fmt.Println("❌ Error inserting user:", err)
		insertUserCreateTracking(r, uid, email, "CreateUserFailed", trackingPayload)
		sendJSONResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	insertUserCreateTracking(r, uid, email, "CreateUser", trackingPayload)

	fmt.Println("✅ Firebase token is valid")
	fmt.Println("UID:", uid)
	sendJSONResponse(w, http.StatusOK, "Token validated and user created successfully.")
}

// Helper function for JSON responses
func sendJSONResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := ResponseMessage{
		Status:  http.StatusText(status),
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

// insertUserCreateTracking writes the tracking row (lt) for /createuser (2026-10-27).
// Same pattern as the other endpoints; the Firebase token is never stored.
func insertUserCreateTracking(r *http.Request, userID, email, action string, payload map[string]interface{}) {
	jsonPayload, _ := json.Marshal(payload)
	clientIP := r.RemoteAddr
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		clientIP = strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             clientIP,
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         action,
		Jsonstring:     string(jsonPayload),
	}
	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	}
}
