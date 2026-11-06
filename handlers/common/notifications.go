package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	busersrequests "github.com/milbertk/usersAppAK/brules/bUsersRequests"
)

type NotificationItem struct {
	Importance     string `json:"importance"`     // high, medium, low
	DescriptionKey string `json:"descriptionKey"` // frontend i18n key
	URL            string `json:"url"`            // where frontend should redirect
	Status         string `json:"status"`         // active, unread, resolved, etc.
}

type NotificationsResponse struct {
	Status        string             `json:"status"`
	Message       string             `json:"message"`
	Notifications []NotificationItem `json:"notifications"`
}

func HGetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendNotificationsResponse(
			w,
			http.StatusMethodNotAllowed,
			"Use GET",
			[]NotificationItem{},
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendNotificationsResponse(
			w,
			http.StatusUnauthorized,
			"Missing user claims",
			[]NotificationItem{},
		)
		return
	}

	userID := strings.TrimSpace(claims.UID)

	notifications, err := buildNotifications(userID)
	if err != nil {
		sendNotificationsResponse(
			w,
			http.StatusBadRequest,
			InternalMessage(err, "Error getting notifications"),
			[]NotificationItem{},
		)
		return
	}

	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"userId":        userID,
		"notifications": notifications,
		"count":         len(notifications),
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIP(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "GetNotifications",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}

	sendNotificationsResponse(
		w,
		http.StatusOK,
		"Notifications obtained successfully.",
		notifications,
	)
}

func buildNotifications(userID string) ([]NotificationItem, error) {
	notifications := []NotificationItem{}

	hasPendingInvitedRequests, err := busersrequests.BPendingRequestUserrequestByCompany(
		userID,
		"invited",
	)
	if err != nil {
		return nil, fmt.Errorf("error checking pending invited requests: %v", err)
	}

	if hasPendingInvitedRequests {
		notifications = append(notifications, NotificationItem{
			Importance:     "high",
			DescriptionKey: "notifications.pendingInvitedRequests",
			URL:            "/manageinvitatinsbyuser",
			Status:         "active",
		})
	}

	return notifications, nil
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

func sendNotificationsResponse(
	w http.ResponseWriter,
	status int,
	message string,
	notifications []NotificationItem,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if notifications == nil {
		notifications = []NotificationItem{}
	}

	_ = json.NewEncoder(w).Encode(NotificationsResponse{
		Status:        http.StatusText(status),
		Message:       message,
		Notifications: notifications,
	})
}
