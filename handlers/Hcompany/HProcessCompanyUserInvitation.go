package hcompany

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	busersrequests "github.com/milbertk/usersAppAK/brules/bUsersRequests"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type ProcessInvitationInput struct {
	RequestID string `json:"requestId"`
	Action    string `json:"action"`
	Message   string `json:"message"`
}

type ProcessInvitationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HProcessCompanyUserInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendProcessInvitationResponse(w, http.StatusMethodNotAllowed, "Use POST or PUT")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendProcessInvitationResponse(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	currentUserID := strings.TrimSpace(claims.UID)

	var input ProcessInvitationInput

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&input); err != nil {
		sendProcessInvitationResponse(w, http.StatusBadRequest, common.InvalidJSONMessage(err))
		return
	}

	input.RequestID = strings.TrimSpace(input.RequestID)
	input.Action = strings.TrimSpace(strings.ToUpper(input.Action))
	input.Message = strings.TrimSpace(input.Message)

	if input.RequestID == "" {
		sendProcessInvitationResponse(w, http.StatusBadRequest, "requestId is required")
		return
	}

	if input.Action == "" {
		sendProcessInvitationResponse(w, http.StatusBadRequest, "action is required")
		return
	}

	switch input.Action {
	case "REJECT", "ACCEPT":
		// allowed values
	default:
		sendProcessInvitationResponse(w, http.StatusBadRequest, "action must be REJECT or ACCEPT")
		return
	}

	if dec.More() {
		sendProcessInvitationResponse(w, http.StatusBadRequest, "only one JSON object is allowed")
		return
	}

	if err := busersrequests.BProcessInvitedUserRequest(
		input.RequestID,
		input.Action,
		input.Message,
		currentUserID,
	); err != nil {
		sendProcessInvitationResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"requestId": input.RequestID,
		"action":    input.Action,
		"message":   input.Message,
		"userId":    currentUserID,
	})

	lt := tracking.LoginTracking{
		UserID:         currentUserID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPInvitation(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "ProcessCompanyUserInvitation",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}

	successMessage := "invitation processed successfully"

	if input.Action == "ACCEPT" {
		successMessage = "invitation accepted successfully"
	} else if input.Action == "REJECT" {
		successMessage = "invitation rejected successfully"
	}

	sendProcessInvitationResponse(w, http.StatusOK, successMessage)
}

func readClientIPInvitation(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func sendProcessInvitationResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ProcessInvitationResponse{
		Status:  http.StatusText(status),
		Message: message,
	})
}
