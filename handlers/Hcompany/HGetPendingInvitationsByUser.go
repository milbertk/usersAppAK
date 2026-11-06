package hcompany

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	busersrequests "github.com/milbertk/usersAppAK/brules/bUsersRequests"
	"github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type PendingInvitationsResponse struct {
	Status        string                                                     `json:"status"`
	Message       string                                                     `json:"message"`
	Notifications []databaseRequestUserCompany.CompanyUserRequestWithCompany `json:"notifications"`
}

func HGetPendingInvitationsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendPendingInvitationsResponse(
			w,
			http.StatusMethodNotAllowed,
			"Use GET",
			nil,
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendPendingInvitationsResponse(
			w,
			http.StatusUnauthorized,
			"Missing user claims",
			nil,
		)
		return
	}

	userID := strings.TrimSpace(claims.UID)

	requests, err := busersrequests.BGetPendingInvitationsByUser(userID)
	if err != nil {
		msg := err.Error()
		lowerMsg := strings.ToLower(msg)

		if strings.Contains(lowerMsg, "user id is required") {
			sendPendingInvitationsResponse(w, http.StatusBadRequest, msg, nil)
			return
		}

		sendPendingInvitationsResponse(
			w,
			http.StatusBadRequest,
			common.InternalMessage(err, "Error getting pending invitations"),
			nil,
		)
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"userID":        userID,
		"notifications": requests,
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPGNU(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "GetPendingInvitationsByUser",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}

	sendPendingInvitationsResponse(
		w,
		http.StatusOK,
		"Pending invitations retrieved successfully.",
		requests,
	)
}

func readClientIPGNU(r *http.Request) string {
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

func sendPendingInvitationsResponse(
	w http.ResponseWriter,
	status int,
	message string,
	notifications []databaseRequestUserCompany.CompanyUserRequestWithCompany,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if notifications == nil {
		notifications = []databaseRequestUserCompany.CompanyUserRequestWithCompany{}
	}

	_ = json.NewEncoder(w).Encode(PendingInvitationsResponse{
		Status:        http.StatusText(status),
		Message:       message,
		Notifications: notifications,
	})
}
