package securitygen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword,omitempty"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ChangePasswordResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var firebaseAuthClient *firebaseauth.Client

func HChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendChangePasswordResponse(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendChangePasswordResponse(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	userID := strings.TrimSpace(claims.UID)

	var req ChangePasswordRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendChangePasswordResponse(w, http.StatusBadRequest, common.InvalidJSONMessage(err))
		return
	}

	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

	if req.NewPassword == "" {
		sendChangePasswordResponse(w, http.StatusBadRequest, "New password is required")
		return
	}

	if req.ConfirmPassword == "" {
		sendChangePasswordResponse(w, http.StatusBadRequest, "Confirm password is required")
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		sendChangePasswordResponse(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if err := validatePassword(req.NewPassword); err != nil {
		sendChangePasswordResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userDb, err := databaseusers.SelectByUID(userID)
	if err != nil {
		sendChangePasswordResponse(
			w,
			http.StatusBadRequest,
			common.InternalMessage(err, "Error loading user"),
		)
		return
	}

	providerID := normalizeChangePasswordProviderID(userDb.ProviderID)

	if !isPasswordProviderID(providerID) {
		sendChangePasswordResponse(
			w,
			http.StatusBadRequest,
			"Password change is only available for email/password accounts. Current provider: "+providerID,
		)
		return
	}

	authClient, err := getFirebaseAuthClient(r.Context())
	if err != nil {
		sendChangePasswordResponse(
			w,
			http.StatusInternalServerError,
			common.InternalMessage(err, "The password could not be changed right now. Please try again later."),
		)
		return
	}

	params := (&firebaseauth.UserToUpdate{}).Password(req.NewPassword)

	_, err = authClient.UpdateUser(r.Context(), userID, params)
	if err != nil {
		sendChangePasswordResponse(
			w,
			http.StatusBadRequest,
			common.ClientMessage(err, "The password could not be changed. Please try again."),
		)
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"userId":   userID,
		"action":   "ChangePassword",
		"provider": providerID,
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPChangePassword(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "ChangePassword",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}

	sendChangePasswordResponse(w, http.StatusOK, "Password updated successfully.")
}

func normalizeChangePasswordProviderID(providerID string) string {
	providerID = strings.ToLower(strings.TrimSpace(providerID))

	if providerID == "" {
		return "unknown"
	}

	return providerID
}

func isPasswordProviderID(providerID string) bool {
	providerID = normalizeChangePasswordProviderID(providerID)

	return providerID == "password"
}

func getFirebaseAuthClient(ctx context.Context) (*firebaseauth.Client, error) {
	if firebaseAuthClient != nil {
		return firebaseAuthClient, nil
	}

	credentialsFile := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_FILE"))

	if credentialsFile == "" {
		credentialsFile = strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	}

	if credentialsFile == "" {
		credentialsFile = "./kdauth-9893b-firebase-adminsdk-82p6v-f263a3940b.json"
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf(
			"could not initialize Firebase app using credentials file %q: %w",
			credentialsFile,
			err,
		)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	firebaseAuthClient = client

	return firebaseAuthClient, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("Password must have at least 8 characters")
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("Password must include at least one uppercase letter")
	}

	if !hasLower {
		return fmt.Errorf("Password must include at least one lowercase letter")
	}

	if !hasNumber {
		return fmt.Errorf("Password must include at least one number")
	}

	if !hasSpecial {
		return fmt.Errorf("Password must include at least one special character")
	}

	return nil
}

func readClientIPChangePassword(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func sendChangePasswordResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ChangePasswordResponse{
		Status:  http.StatusText(status),
		Message: message,
	})
}
