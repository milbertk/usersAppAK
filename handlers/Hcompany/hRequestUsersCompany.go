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

type ResponseMessageur struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CompanyUserRequestCreateRequest struct {
	CompanyID     string `json:"companyId"`
	InvitedUserID string `json:"invitedUserId"`
	RoleRequested string `json:"roleRequested"`
	Message       string `json:"message"`
	Status        string `json:"status"`
}

func HCreateCompanyUserRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponseur(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendJSONResponseur(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	requestedByUserID := strings.TrimSpace(claims.UID)

	var req CompanyUserRequestCreateRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		sendJSONResponseur(w, http.StatusBadRequest, common.InvalidJSONMessage(err))
		return
	}

	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.InvitedUserID = strings.TrimSpace(req.InvitedUserID)
	req.RoleRequested = strings.TrimSpace(req.RoleRequested)
	req.Message = strings.TrimSpace(req.Message)
	req.Status = strings.TrimSpace(req.Status)

	if req.CompanyID == "" {
		sendJSONResponseur(w, http.StatusBadRequest, "companyId is required")
		return
	}

	if req.InvitedUserID == "" {
		sendJSONResponseur(w, http.StatusBadRequest, "invitedUserId is required")
		return
	}

	if req.RoleRequested == "" {
		req.RoleRequested = "USER"
	}

	if req.Status == "" {
		req.Status = "PENDING"
	}

	request := databaseRequestUserCompany.CompanyUserRequest{
		CompanyID:         req.CompanyID,
		InvitedUserID:     req.InvitedUserID,
		RequestedByUserID: requestedByUserID,
		RoleRequested:     req.RoleRequested,
		Status:            req.Status,
		Message:           req.Message,
	}

	if err := busersrequests.BCreateCompanyUserRequest(request); err != nil {
		msg := err.Error()
		lowerMsg := strings.ToLower(msg)

		if strings.Contains(lowerMsg, "already a pending request") ||
			strings.Contains(lowerMsg, "already exists") ||
			strings.Contains(lowerMsg, "exists") {
			sendJSONResponseur(w, http.StatusConflict, msg)
			return
		}

		if strings.Contains(lowerMsg, "only the company owner can create this request") {
			sendJSONResponseur(w, http.StatusForbidden, msg)
			return
		}

		if strings.Contains(lowerMsg, "does not exist") {
			sendJSONResponseur(w, http.StatusNotFound, msg)
			return
		}

		sendJSONResponseur(w, http.StatusBadRequest, msg)
		return
	}

	jsonRequest, _ := json.Marshal(request)

	lt := tracking.LoginTracking{
		UserID:         requestedByUserID,
		Email:          claims.Email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPur(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "CreateCompanyUserRequest",
		Jsonstring:     string(jsonRequest),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}

	sendJSONResponseur(w, http.StatusOK, "Company user request created successfully.")
}

func readClientIPur(r *http.Request) string {
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

func sendJSONResponseur(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ResponseMessageur{
		Status:  http.StatusText(status),
		Message: message,
	})
}
