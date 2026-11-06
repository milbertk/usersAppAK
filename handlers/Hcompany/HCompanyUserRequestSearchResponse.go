package hcompany

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	busersrequests "github.com/milbertk/usersAppAK/brules/bUsersRequests"
	databaseRequestUserCompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type CompanyUserRequestSearchResponse struct {
	Status   string                                                      `json:"status"`
	Message  string                                                      `json:"message"`
	Requests []databaseRequestUserCompany.CompanyUserRequestSearchResult `json:"requests"`
	Page     int                                                         `json:"page"`
	PageSize int                                                         `json:"pageSize"`
	Total    int                                                         `json:"total"`
}

func HSearchCompanyUserRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSONResponse(w, http.StatusMethodNotAllowed, "Use GET")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendJSONResponse(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	requestedByUserID := strings.TrimSpace(claims.UID)

	companyID := strings.TrimSpace(r.URL.Query().Get("companyId"))
	searchText := strings.TrimSpace(r.URL.Query().Get("search"))
	pageStr := strings.TrimSpace(r.URL.Query().Get("page"))
	pageSizeStr := strings.TrimSpace(r.URL.Query().Get("pageSize"))

	if companyID == "" {
		sendJSONResponse(w, http.StatusBadRequest, "companyId is required")
		return
	}

	page := 1
	pageSize := 10

	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage <= 0 {
			sendJSONResponse(w, http.StatusBadRequest, "page must be a positive number")
			return
		}

		page = parsedPage
	}

	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize <= 0 {
			sendJSONResponse(w, http.StatusBadRequest, "pageSize must be a positive number")
			return
		}

		pageSize = parsedPageSize
	}

	if pageSize > 100 {
		pageSize = 100
	}

	results, err := busersrequests.BSearchRequest(
		companyID,
		searchText,
		requestedByUserID,
	)

	if err != nil {
		msg := err.Error()
		lowerMsg := strings.ToLower(msg)

		if strings.Contains(lowerMsg, "not allowed") ||
			strings.Contains(lowerMsg, "you are not allowed") ||
			strings.Contains(lowerMsg, "only the company owner") {
			sendJSONResponse(w, http.StatusForbidden, msg)
			return
		}

		if strings.Contains(lowerMsg, "company not found") ||
			strings.Contains(lowerMsg, "no active company found") ||
			strings.Contains(lowerMsg, "no rows") {
			sendJSONResponse(w, http.StatusNotFound, msg)
			return
		}

		sendJSONResponse(w, http.StatusBadRequest, common.ClientMessage(err, "Error searching requests"))
		return
	}

	total := len(results)

	start := (page - 1) * pageSize
	if start >= total {
		start = total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	pagedResults := results[start:end]

	insertSearchCompanyUserRequestsTracking(
		r,
		requestedByUserID,
		claims.Email,
		companyID,
		searchText,
		page,
		pageSize,
		total,
	)

	sendJSONResponseRequests(
		w,
		http.StatusOK,
		"Company user requests retrieved successfully.",
		pagedResults,
		page,
		pageSize,
		total,
	)
}

func insertSearchCompanyUserRequestsTracking(
	r *http.Request,
	requestedByUserID string,
	userEmail string,
	companyID string,
	searchText string,
	page int,
	pageSize int,
	total int,
) {
	jsonSearch, _ := json.Marshal(map[string]interface{}{
		"companyId": companyID,
		"search":    searchText,
		"page":      page,
		"pageSize":  pageSize,
		"total":     total,
	})

	lt := tracking.LoginTracking{
		UserID:         requestedByUserID,
		Email:          userEmail,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIP(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "SearchCompanyUserRequests",
		Jsonstring:     string(jsonSearch),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}
}

func sendJSONResponseRequests(
	w http.ResponseWriter,
	status int,
	message string,
	requests []databaseRequestUserCompany.CompanyUserRequestSearchResult,
	page int,
	pageSize int,
	total int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(CompanyUserRequestSearchResponse{
		Status:   http.StatusText(status),
		Message:  message,
		Requests: requests,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}
