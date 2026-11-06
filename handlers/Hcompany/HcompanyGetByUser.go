package hcompany

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/milbertk/security"

	bCreateCompany "github.com/milbertk/usersAppAK/brules/bCompany"
	databaseuserofcompany "github.com/milbertk/usersAppAK/database/dataBaseUserOfCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type CompanyByUserListResponse struct {
	Status    string                                      `json:"status"`
	Message   string                                      `json:"message"`
	Companies []databaseuserofcompany.CompanyByUserResult `json:"companies"`
}

func HcompanyGetByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCompanyByUserListResponse(
			w,
			http.StatusMethodNotAllowed,
			"Use GET",
			nil,
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendCompanyByUserListResponse(
			w,
			http.StatusUnauthorized,
			"Missing user claims",
			nil,
		)
		return
	}

	userID := strings.TrimSpace(claims.UID)

	companies, err := bCreateCompany.BSearchCompaniesByUserID(userID)
	if err != nil {
		sendCompanyByUserListResponse(
			w,
			http.StatusInternalServerError,
			common.InternalMessage(err, "Error retrieving user companies"),
			nil,
		)
		return
	}

	sendCompanyByUserListResponse(
		w,
		http.StatusOK,
		"User companies retrieved successfully.",
		companies,
	)
}

func sendCompanyByUserListResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	companies []databaseuserofcompany.CompanyByUserResult,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if companies == nil {
		companies = []databaseuserofcompany.CompanyByUserResult{}
	}

	_ = json.NewEncoder(w).Encode(CompanyByUserListResponse{
		Status:    http.StatusText(statusCode),
		Message:   message,
		Companies: companies,
	})
}
