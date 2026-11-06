package hcompany

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/milbertk/security"

	bCreateCompany "github.com/milbertk/usersAppAK/brules/bCompany"
	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type CompanyListResponse struct {
	Status    string                    `json:"status"`
	Message   string                    `json:"message"`
	Companies []databasecompany.Company `json:"companies"`
}

func HcompanyGetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendCompanyListResponse(w, http.StatusMethodNotAllowed, "Use GET", nil)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendCompanyListResponse(w, http.StatusUnauthorized, "Missing user claims", nil)
		return
	}

	ownerUID := strings.TrimSpace(claims.UID)

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	plan := strings.TrimSpace(r.URL.Query().Get("plan"))
	includeInactive := strings.TrimSpace(r.URL.Query().Get("includeInactive"))

	companies, err := bCreateCompany.BSelectAllCompanies(ownerUID)
	if err != nil {
		sendCompanyListResponse(
			w,
			http.StatusInternalServerError,
			common.InternalMessage(err, "Error retrieving companies"),
			nil,
		)
		return
	}

	filtered := make([]databasecompany.Company, 0, len(companies))

	for _, c := range companies {
		if status != "" && !strings.EqualFold(c.Status, status) {
			continue
		}

		if plan != "" && !strings.EqualFold(c.Plan, plan) {
			continue
		}

		if includeInactive != "true" && !c.Active {
			continue
		}

		if search != "" && !companyMatchesSearch(c, search) {
			continue
		}

		filtered = append(filtered, c)
	}

	sendCompanyListResponse(
		w,
		http.StatusOK,
		"Companies retrieved successfully.",
		filtered,
	)
}

func companyMatchesSearch(c databasecompany.Company, search string) bool {
	search = strings.ToLower(strings.TrimSpace(search))

	return strings.Contains(strings.ToLower(c.CompanyID), search) ||
		strings.Contains(strings.ToLower(c.CompanyName), search) ||
		strings.Contains(strings.ToLower(c.LegalName), search) ||
		strings.Contains(strings.ToLower(c.IdentificationType), search) ||
		strings.Contains(strings.ToLower(c.IdentificationNumber), search) ||
		strings.Contains(strings.ToLower(c.Email), search) ||
		strings.Contains(strings.ToLower(c.Phone), search) ||
		strings.Contains(strings.ToLower(c.Category), search) ||
		strings.Contains(strings.ToLower(c.InternalCategory), search)
}

func sendCompanyListResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	companies []databasecompany.Company,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(CompanyListResponse{
		Status:    http.StatusText(statusCode),
		Message:   message,
		Companies: companies,
	})
}
