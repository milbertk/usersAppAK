package securitygen

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type MeResponse struct {
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	Authenticated bool      `json:"authenticated"`
	Inactive      bool      `json:"inactive"`
	EmailVerified bool      `json:"email_verified"`
	User          *UserData `json:"user"`
}

type UserData struct {
	UID            string `json:"uid"`
	Email          string `json:"email"`
	Name           string `json:"name,omitempty"`
	Picture        string `json:"picture,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Role           string `json:"role,omitempty"`
	SignInProvider string `json:"sign_in_provider,omitempty"`
	ProviderID     string `json:"provider_id,omitempty"`

	LocalStatus   string `json:"local_status"`
	Active        string `json:"active"`
	Disabled      string `json:"disabled"`
	EmailVerified string `json:"email_verified_raw,omitempty"`

	LastSignIn   string `json:"last_sign_in,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`

	Issuer    string `json:"issuer,omitempty"`
	IssuedAt  string `json:"issued_at,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func HuserMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMeJSON(w, http.StatusMethodNotAllowed, MeResponse{
			Status:        "MethodNotAllowed",
			Message:       "Use GET",
			Authenticated: false,
			User:          nil,
		})
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		writeMeJSON(w, http.StatusUnauthorized, MeResponse{
			Status:        "Unauthorized",
			Message:       "Missing user claims",
			Authenticated: false,
			User:          nil,
		})
		return
	}

	userDb, err := databaseusers.SelectByUID(claims.UID)
	if err != nil {
		writeMeJSON(w, http.StatusUnauthorized, MeResponse{
			Status:        "Unauthorized",
			Message:       common.InternalMessage(err, "User not found"),
			Authenticated: false,
			User:          nil,
		})
		return
	}

	if strings.TrimSpace(userDb.UID) == "" {
		writeMeJSON(w, http.StatusUnauthorized, MeResponse{
			Status:        "Unauthorized",
			Message:       "User does not exist",
			Authenticated: false,
			User:          nil,
		})
		return
	}

	disabled := equalsCI(userDb.Disabled, "true")
	active := equalsCI(userDb.Active, "true")
	localActive := equalsCI(userDb.LocalStatus, "ACTIVE")
	emailVerified := equalsCI(userDb.EmailVerified, "true")

	inactive := disabled || !active || !localActive

	user := &UserData{
		UID:            userDb.UID,
		Email:          userDb.Email,
		Name:           userDb.DisplayName,
		Picture:        userDb.PhotoURL,
		PhoneNumber:    userDb.PhoneNumber,
		Role:           claims.Role,
		SignInProvider: claims.SignInProvider,
		ProviderID:     userDb.ProviderID,

		LocalStatus:   userDb.LocalStatus,
		Active:        userDb.Active,
		Disabled:      userDb.Disabled,
		EmailVerified: userDb.EmailVerified,

		LastSignIn:   userDb.LastSignIn,
		CreationTime: userDb.CreationTime,

		Issuer: claims.Issuer,
	}

	if claims.IssuedAt != nil {
		user.IssuedAt = claims.IssuedAt.UTC().Format(time.RFC3339)
	}

	if claims.ExpiresAt != nil {
		user.ExpiresAt = claims.ExpiresAt.UTC().Format(time.RFC3339)
	}

	if inactive {
		writeMeJSON(w, http.StatusOK, MeResponse{
			Status:        "Inactive",
			Message:       buildInactiveMessage(disabled, active, localActive),
			Authenticated: false,
			Inactive:      true,
			EmailVerified: emailVerified,
			User:          user,
		})
		return
	}

	writeMeJSON(w, http.StatusOK, MeResponse{
		Status:        "OK",
		Message:       "Session active",
		Authenticated: true,
		Inactive:      false,
		EmailVerified: emailVerified,
		User:          user,
	})
}

func equalsCI(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func buildInactiveMessage(disabled, active, localActive bool) string {
	switch {
	case disabled:
		return "Account is disabled"
	case !active:
		return "Account is not active"
	case !localActive:
		return "Account local status is not ACTIVE"
	default:
		return "Account is inactive"
	}
}

func writeMeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
