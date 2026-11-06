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

type UpdateRequestInput struct {
	RequestID string `json:"requestId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type UpdateRequestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HUpdateCompanyUserRequestStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendUpdateResponseCBM(w, http.StatusMethodNotAllowed, "Use POST or PUT")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendUpdateResponseCBM(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	currentUserID := strings.TrimSpace(claims.UID)

	var input UpdateRequestInput

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&input); err != nil {
		sendUpdateResponseCBM(w, http.StatusBadRequest, common.InvalidJSONMessage(err))
		return
	}

	input.RequestID = strings.TrimSpace(input.RequestID)
	input.Status = strings.TrimSpace(strings.ToUpper(input.Status))
	input.Message = strings.TrimSpace(input.Message)

	if input.RequestID == "" {
		sendUpdateResponseCBM(w, http.StatusBadRequest, "requestId is required")
		return
	}

	if input.Status == "" {
		sendUpdateResponseCBM(w, http.StatusBadRequest, "status is required")
		return
	}

	if err := busersrequests.BUpdateCompanyUserRequestStatus(
		input.RequestID,
		input.Status,
		input.Message,
		currentUserID,
	); err != nil {
		sendUpdateResponseCBM(w, http.StatusInternalServerError, common.ClientMessage(err, "The request could not be updated"))
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"requestId": input.RequestID,
		"status":    input.Status,
		"message":   input.Message,
		"userId":    currentUserID,
	})

	lt := tracking.LoginTracking{
		UserID:         currentUserID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPCBM(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "UpdateCompanyUserRequestStatus",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}

	sendUpdateResponseCBM(w, http.StatusOK, "request updated successfully")
}

func readClientIPCBM(r *http.Request) string {
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

func sendUpdateResponseCBM(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(UpdateRequestResponse{
		Status:  http.StatusText(status),
		Message: message,
	})
}
