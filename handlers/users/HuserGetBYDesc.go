package Husers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	bcreateusers "github.com/milbertk/usersAppAK/brules/bcreateUsers"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
	"github.com/milbertk/usersAppAK/handlers/common"
)

type UserSearchResponse struct {
	Status  string                                   `json:"status"`
	Message string                                   `json:"message"`
	Users   []databaseusers.FirebaseUserSearchResult `json:"users"`
}

// HsearchUsersByEmailNameOrId finds ONE user by the exact email, to invite
// them to a company (2026-10-25, security S10).
//
//		GET /searchUser?search=<complete email>&companyId=<company>
//
//	  - Only the OWNER / MANAGER of companyId can search (the route uses
//	    security.RequireActiveCompanyUser; companyId may also come in the
//	    X-Company-ID header).
//	  - Only an exact, complete email; at most one user is returned. No partial
//	    text, names or ids, and never the list of users.
//	  - Every search is tracked (who searched, for which company, found or not).
func HsearchUsersByEmailNameOrId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendUserSearchResponse(w, http.StatusMethodNotAllowed, "Use GET", nil)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendUserSearchResponse(w, http.StatusUnauthorized, "Missing user claims", nil)
		return
	}
	requestedByUID := strings.TrimSpace(claims.UID)
	companyID, _ := security.GetActiveCompanyID(r)

	searchText := strings.TrimSpace(r.URL.Query().Get("search"))
	if strings.EqualFold(searchText, "null") {
		searchText = ""
	}

	users, err := bcreateusers.BsearchUsersByEmailNameOrId(searchText)
	if errors.Is(err, bcreateusers.ErrInvalidUserSearch) {
		insertUserSearchTracking(r, requestedByUID, claims.Email, companyID, false, "INVALID_SEARCH")
		sendUserSearchResponse(w, http.StatusBadRequest, "Type the complete email address of the user.", nil)
		return
	}
	if err != nil {
		insertUserSearchTracking(r, requestedByUID, claims.Email, companyID, false, "ERROR")
		sendUserSearchResponse(w, http.StatusInternalServerError, common.InternalMessage(err, "Error searching users"), nil)
		return
	}

	insertUserSearchTracking(r, requestedByUID, claims.Email, companyID, len(users) > 0, "OK")
	sendUserSearchResponse(w, http.StatusOK, "Users retrieved successfully.", users)
}

// insertUserSearchTracking records the search (never the searched email).
func insertUserSearchTracking(r *http.Request, userID, email, companyID string, found bool, result string) {
	jsonRequest, _ := json.Marshal(map[string]interface{}{
		"companyId": strings.TrimSpace(companyID),
		"found":     found,
		"result":    result,
	})

	lt := tracking.LoginTracking{
		UserID:         userID,
		Email:          email,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPUserSearch(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "SearchUserByEmail",
		Jsonstring:     string(jsonRequest),
	}
	if err := lt.Insert(); err != nil {
		log.Printf("[HsearchUsersByEmailNameOrId] Error inserting tracking: %v", err)
	}
}

func readClientIPUserSearch(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

func sendUserSearchResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	users []databaseusers.FirebaseUserSearchResult,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if users == nil {
		users = []databaseusers.FirebaseUserSearchResult{}
	}

	_ = json.NewEncoder(w).Encode(UserSearchResponse{
		Status:  http.StatusText(statusCode),
		Message: message,
		Users:   users,
	})
}
