package securitygen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type ChangeEmailRequest struct {
	CurrentEmail string `json:"currentEmail"`
	NewEmail     string `json:"newEmail"`
	Password     string `json:"password,omitempty"`
}

type ChangeEmailResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}

var changeEmailFirebaseAuthClient *firebaseauth.Client

func HChangeEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendChangeEmailResponse(w, http.StatusMethodNotAllowed, "Use POST or PUT", "")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendChangeEmailResponse(w, http.StatusUnauthorized, "Missing user claims", "")
		return
	}

	userID := strings.TrimSpace(claims.UID)

	var req ChangeEmailRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendChangeEmailResponse(w, http.StatusBadRequest, common.InvalidJSONMessage(err), "")
		return
	}

	req.CurrentEmail = strings.ToLower(strings.TrimSpace(req.CurrentEmail))
	req.NewEmail = strings.ToLower(strings.TrimSpace(req.NewEmail))
	req.Password = strings.TrimSpace(req.Password)

	if req.CurrentEmail == "" {
		sendChangeEmailResponse(w, http.StatusBadRequest, "Current email is required", "")
		return
	}

	if req.NewEmail == "" {
		sendChangeEmailResponse(w, http.StatusBadRequest, "New email is required", "")
		return
	}

	if !isValidEmail(req.CurrentEmail) {
		sendChangeEmailResponse(w, http.StatusBadRequest, "Current email is not valid", "")
		return
	}

	if !isValidEmail(req.NewEmail) {
		sendChangeEmailResponse(w, http.StatusBadRequest, "New email is not valid", "")
		return
	}

	if req.CurrentEmail == req.NewEmail {
		sendChangeEmailResponse(w, http.StatusBadRequest, "The new email must be different from the current email", "")
		return
	}

	userDb, err := databaseusers.SelectByUID(userID)
	if err != nil {
		sendChangeEmailResponse(w, http.StatusBadRequest, common.InternalMessage(err, "Error loading user"), "")
		return
	}

	providerID := normalizeChangeEmailProviderID(userDb.ProviderID)

	if !isEmailPasswordProvider(providerID) {
		sendChangeEmailResponse(
			w,
			http.StatusBadRequest,
			"Email change is only available for email/password accounts. Current provider: "+providerID,
			"",
		)
		return
	}

	dbEmail := strings.ToLower(strings.TrimSpace(userDb.Email))

	if dbEmail == "" {
		sendChangeEmailResponse(w, http.StatusBadRequest, "The current account does not have a registered email", "")
		return
	}

	if dbEmail != req.CurrentEmail {
		sendChangeEmailResponse(w, http.StatusBadRequest, "Current email does not match the account email", "")
		return
	}

	authClient, err := getChangeEmailFirebaseAuthClient(r.Context())
	if err != nil {
		sendChangeEmailResponse(
			w,
			http.StatusInternalServerError,
			common.InternalMessage(err, "The email could not be changed right now. Please try again later."),
			"",
		)
		return
	}

	params := (&firebaseauth.UserToUpdate{}).
		Email(req.NewEmail).
		EmailVerified(false)

	_, err = authClient.UpdateUser(r.Context(), userID, params)
	if err != nil {
		sendChangeEmailResponse(
			w,
			http.StatusBadRequest,
			common.ClientMessage(err, "The email could not be changed. Please check it and try again."),
			"",
		)
		return
	}

	userForUpdate := databaseusers.FirebaseUser{
		UID:           userDb.UID,
		LastSignIn:    userDb.LastSignIn,
		LocalStatus:   userDb.LocalStatus,
		CreationTime:  userDb.CreationTime,
		DisplayName:   userDb.DisplayName,
		PhoneNumber:   userDb.PhoneNumber,
		PhotoURL:      userDb.PhotoURL,
		Active:        userDb.Active,
		ProviderID:    userDb.ProviderID,
		Email:         req.NewEmail,
		EmailVerified: "false",
		Disabled:      userDb.Disabled,
	}

	if err := userForUpdate.Update(); err != nil {
		sendChangeEmailResponse(
			w,
			http.StatusBadRequest,
			common.InternalMessage(err, "The email was changed, but your profile could not be updated. Please contact support."),
			req.NewEmail,
		)
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"userId":       userID,
		"action":       "ChangeEmail",
		"provider":     providerID,
		"currentEmail": req.CurrentEmail,
		"newEmail":     req.NewEmail,
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          req.NewEmail,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPChangeEmail(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "ChangeEmail",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}

	sendChangeEmailResponse(
		w,
		http.StatusOK,
		"Email updated successfully. Please verify the new email address.",
		req.NewEmail,
	)
}

func normalizeChangeEmailProviderID(providerID string) string {
	providerID = strings.ToLower(strings.TrimSpace(providerID))

	if providerID == "" {
		return "unknown"
	}

	return providerID
}

func isEmailPasswordProvider(providerID string) bool {
	providerID = normalizeChangeEmailProviderID(providerID)

	return providerID == "password"
}

func isValidEmail(value string) bool {
	value = strings.TrimSpace(value)

	if value == "" {
		return false
	}

	_, err := mail.ParseAddress(value)
	return err == nil
}

func getChangeEmailFirebaseAuthClient(ctx context.Context) (*firebaseauth.Client, error) {
	if changeEmailFirebaseAuthClient != nil {
		return changeEmailFirebaseAuthClient, nil
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

	changeEmailFirebaseAuthClient = client

	return changeEmailFirebaseAuthClient, nil
}

func readClientIPChangeEmail(r *http.Request) string {
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

func sendChangeEmailResponse(w http.ResponseWriter, status int, message string, email string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ChangeEmailResponse{
		Status:  http.StatusText(status),
		Message: message,
		Email:   email,
	})
}
