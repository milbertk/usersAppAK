package hcompany

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/milbertk/security"

	buserofcompany "github.com/milbertk/usersAppAK/brules/busersOfCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type SearchCompanyUsersResponse struct {
	Status  string                                   `json:"status"`
	Message string                                   `json:"message"`
	Data    *buserofcompany.SearchCompanyUsersResult `json:"data,omitempty"`
}

func HSearchCompanyUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendSearchCompanyUsersResponse(
			w,
			http.StatusMethodNotAllowed,
			"ERROR",
			"Use GET",
			nil,
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendSearchCompanyUsersResponse(
			w,
			http.StatusUnauthorized,
			"ERROR",
			"Missing user claims",
			nil,
		)
		return
	}

	currentUserID := strings.TrimSpace(claims.UID)

	companyID := strings.TrimSpace(r.URL.Query().Get("companyId"))
	includeInactive := false

	includeInactiveParam := strings.TrimSpace(r.URL.Query().Get("includeInactive"))
	if includeInactiveParam != "" {
		parsed, err := strconv.ParseBool(includeInactiveParam)
		if err != nil {
			sendSearchCompanyUsersResponse(
				w,
				http.StatusBadRequest,
				"ERROR",
				"invalid includeInactive value",
				nil,
			)
			return
		}

		includeInactive = parsed
	}

	if companyID == "" {
		sendSearchCompanyUsersResponse(
			w,
			http.StatusBadRequest,
			"ERROR",
			"companyId is required",
			nil,
		)
		return
	}

	log.Printf(
		"[HSearchCompanyUsers] companyID=%q userID=%q includeInactive=%v",
		companyID,
		currentUserID,
		includeInactive,
	)

	result, err := buserofcompany.BSearchCompanyUsers(
		companyID,
		currentUserID,
		includeInactive,
	)

	if err != nil {
		msg := err.Error()
		lowerMsg := strings.ToLower(msg)

		switch {
		case strings.Contains(lowerMsg, "unauthorized") ||
			strings.Contains(lowerMsg, "not allowed") ||
			strings.Contains(lowerMsg, "forbidden"):
			sendSearchCompanyUsersResponse(w, http.StatusForbidden, "ERROR", msg, nil)

		case strings.Contains(lowerMsg, "required"):
			sendSearchCompanyUsersResponse(w, http.StatusBadRequest, "ERROR", msg, nil)

		case strings.Contains(lowerMsg, "not found") ||
			strings.Contains(lowerMsg, "no rows"):
			sendSearchCompanyUsersResponse(w, http.StatusNotFound, "ERROR", msg, nil)

		default:
			sendSearchCompanyUsersResponse(w, http.StatusInternalServerError, "ERROR", common.InternalMessage(err, "Error loading company users"), nil)
		}

		return
	}

	sendSearchCompanyUsersResponse(
		w,
		http.StatusOK,
		"OK",
		"company users loaded successfully",
		result,
	)
}

func sendSearchCompanyUsersResponse(
	w http.ResponseWriter,
	httpStatus int,
	status string,
	message string,
	data *buserofcompany.SearchCompanyUsersResult,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	_ = json.NewEncoder(w).Encode(SearchCompanyUsersResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
