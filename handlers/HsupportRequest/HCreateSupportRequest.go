package hsupportrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	bcaptcha "github.com/milbertk/usersAppAK/brules/bCaptcha"
	bsupportrequest "github.com/milbertk/usersAppAK/brules/bSupportRequest"
)

const supportTrackingPlatform = "UserManagement"

type supportResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func sendSupportResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(supportResponse{Status: http.StatusText(status), Message: message, Data: data})
}

func readClientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// insertSupportTracking writes the tracking row (lt).
func insertSupportTracking(r *http.Request, userID, email, action string, payload map[string]interface{}) {
	jsonPayload, _ := json.Marshal(payload)
	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIP(r),
		Platform:       supportTrackingPlatform,
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

// HCreateSupportRequest saves a support request and emails it (2026-10-24).
//
//	POST /supportrequests   (signed-in users)
//	body: { companyId, companyName, platform, category, priority, subject,
//	        description, contactMethod, contactValue, preferredDate,
//	        timeSlots[], imageUrls[], currentPath, language }
func HCreateSupportRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendSupportResponse(w, http.StatusMethodNotAllowed, "Use POST", nil)
		return
	}
	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendSupportResponse(w, http.StatusUnauthorized, "Missing user claims", nil)
		return
	}
	userID := strings.TrimSpace(claims.UID)

	// CAPTCHA (Cloudflare Turnstile, 2026-10-25): off when
	// TURNSTILE_SECRET_KEY is empty.
	if err := bcaptcha.BVerifyCaptcha(r.Header.Get(bcaptcha.HeaderName)); err != nil {
		code := bcaptcha.Code(err)
		insertSupportTracking(r, userID, claims.Email, "SupportRequestCaptchaFailed", map[string]interface{}{"code": code})
		status := http.StatusBadRequest
		if code == "CAPTCHA_UNAVAILABLE" {
			status = http.StatusServiceUnavailable
		}
		sendSupportResponse(w, status, "Please complete the security check and try again.", map[string]string{"code": code})
		return
	}

	var input bsupportrequest.CreateSupportRequestInput
	decoder := json.NewDecoder(io.LimitReader(r.Body, 16*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		sendSupportResponse(w, http.StatusBadRequest, "invalid JSON body", nil)
		return
	}

	request, err := bsupportrequest.BCreateSupportRequest(userID, claims.Email, input)
	if err != nil {
		status := http.StatusInternalServerError
		message := "The request could not be saved. Please try again."
		switch {
		case errors.Is(err, bsupportrequest.ErrInvalidSupportRequest):
			status, message = http.StatusBadRequest, err.Error()
		case errors.Is(err, bsupportrequest.ErrTooManySupportRequests):
			status, message = http.StatusTooManyRequests, err.Error()
		default:
			fmt.Println("[HCreateSupportRequest] error:", err)
		}
		insertSupportTracking(r, userID, claims.Email, "CreateSupportRequestFailed", map[string]interface{}{
			"status": status, "category": input.Category, "companyId": input.CompanyID,
		})
		sendSupportResponse(w, status, message, nil)
		return
	}

	insertSupportTracking(r, userID, claims.Email, "CreateSupportRequest", map[string]interface{}{
		"ticketCode": request.TicketCode, "category": request.Category, "priority": request.Priority,
		"companyId": request.CompanyID, "platform": request.Platform,
	})
	sendSupportResponse(w, http.StatusOK, "Support request saved.", map[string]interface{}{
		"ticketCode": request.TicketCode,
		"createdAt":  request.CreatedAt,
	})
}
