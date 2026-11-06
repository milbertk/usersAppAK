package haiagent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/tracking"
)

const aiTrackingPlatform = "AIAgent"

type AIAgentHTTPResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func sendAIAgentResponse(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(AIAgentHTTPResponse{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	})
}

// insertAIAgentTracking writes the tracking row (lt) for every AI endpoint.
func insertAIAgentTracking(
	r *http.Request,
	userID string,
	userEmail string,
	action string,
	payload map[string]interface{},
) {
	jsonPayload, _ := json.Marshal(payload)

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          userEmail,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPai(r),
		Platform:       aiTrackingPlatform,
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         action,
		Jsonstring:     string(jsonPayload),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}
}

func readClientIPai(r *http.Request) string {
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

func statusFromAIError(msg string) int {
	lower := strings.ToLower(msg)

	switch {
	case strings.Contains(lower, "does not exist"):
		return http.StatusNotFound
	case strings.Contains(lower, "is required"),
		strings.Contains(lower, "too long"):
		return http.StatusBadRequest
	case strings.Contains(lower, "error calling ai agent"):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
