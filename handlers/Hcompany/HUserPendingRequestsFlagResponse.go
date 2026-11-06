package hcompany

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/tracking"

	busersrequests "github.com/milbertk/usersAppAK/brules/bUsersRequests"
	"github.com/milbertk/usersAppAK/helpers"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type PendingRequestsFlagResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	HasPending bool   `json:"hasPending"`
}

func HUserHasPendingInvitedRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendPendingFlagResponse(w, http.StatusMethodNotAllowed, "Method not allowed", false)
		return
	}

	userID, err := helpers.RequireLocalJWT(r)
	if err != nil {
		sendPendingFlagResponse(w, http.StatusUnauthorized, err.Error(), false)
		return
	}

	hasPending, err := busersrequests.BPendingRequestUserrequestByCompany(userID, "invited")
	if err != nil {
		sendPendingFlagResponse(w, http.StatusBadRequest, common.InternalMessage(err, "Error checking pending requests"), false)
		return
	}

	// Tracking
	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"userId":     userID,
		"hasPending": hasPending,
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          "",
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIP(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "UserHasPendingInvitedRequests",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}

	sendPendingFlagResponse(w, http.StatusOK, "Pending requests flag obtained successfully.", hasPending)
}

func readClientIP(r *http.Request) string {
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

func sendPendingFlagResponse(w http.ResponseWriter, status int, message string, hasPending bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(PendingRequestsFlagResponse{
		Status:     http.StatusText(status),
		Message:    message,
		HasPending: hasPending,
	})
}
