package hcompany

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	bCreateCompany "github.com/milbertk/usersAppAK/brules/bCompany"
	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CompanyCreateRequest struct {
	CompanyID            string `json:"companyId"`
	CompanyName          string `json:"companyName"`
	LegalName            string `json:"legalName"`
	IdentificationType   string `json:"identificationType"`
	IdentificationNumber string `json:"identificationNumber"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	Website              string `json:"website"`
	CountryCode          string `json:"countryCode"`
	CountryName          string `json:"countryName"`
	StateCode            string `json:"stateCode"`
	StateName            string `json:"stateName"`
	City                 string `json:"city"`
	District             string `json:"district"`
	PostalCode           string `json:"postalCode"`
	Address1             string `json:"address1"`
	Address2             string `json:"address2"`
	Directions           string `json:"directions"`
	Category             string `json:"category"`
	InternalCategory     string `json:"internalCategory"`
	Employees            int    `json:"employees"`
	IsVirtual            bool   `json:"isVirtual"`
	Plan                 string `json:"plan"`
	Status               string `json:"status"`
	Active               *bool  `json:"active"`
}

func HcompanyCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, "Use POST")
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendJSONResponse(w, http.StatusUnauthorized, "Missing user claims")
		return
	}

	ownerUID := strings.TrimSpace(claims.UID)

	var req CompanyCreateRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, common.InvalidJSONMessage(err))
		return
	}

	normalizeCompanyRequest(&req)

	if err := validateCompanyRequest(req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now().UTC()

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	if req.CompanyID == "" {
		req.CompanyID = uuid.NewString()
	}

	company := databasecompany.Company{
		CompanyID:            req.CompanyID,
		OwnerFirebaseUID:     ownerUID,
		CompanyName:          req.CompanyName,
		LegalName:            req.LegalName,
		IdentificationType:   req.IdentificationType,
		IdentificationNumber: req.IdentificationNumber,
		Email:                req.Email,
		Phone:                req.Phone,
		Website:              req.Website,
		CountryCode:          req.CountryCode,
		CountryName:          req.CountryName,
		StateCode:            req.StateCode,
		StateName:            req.StateName,
		City:                 req.City,
		District:             req.District,
		PostalCode:           req.PostalCode,
		Address1:             req.Address1,
		Address2:             req.Address2,
		Directions:           req.Directions,
		Category:             req.Category,
		InternalCategory:     req.InternalCategory,
		Employees:            req.Employees,
		IsVirtual:            req.IsVirtual,
		Active:               active,
		// plan and status are NEVER taken from the request (2026-10-25,
		// security S3): new companies start FREE / ACTIVE and an update keeps
		// the saved values (see BCreateOrUpdateCompany). Only AlephK changes them.
		Plan:      "",
		Status:    "",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	if company.Plan == "" {
		company.Plan = "FREE"
	}

	if company.Status == "" {
		company.Status = "ACTIVE"
	}

	action, err := bCreateCompany.BCreateOrUpdateCompany(company)
	if err != nil {
		msg := err.Error()
		msgLower := strings.ToLower(msg)

		if strings.Contains(msgLower, "exists") ||
			strings.Contains(msgLower, "already exists") ||
			strings.Contains(msgLower, "duplicate") {
			sendJSONResponse(w, http.StatusConflict, common.ClientMessage(err, "The company name or identification number already exists"))
			return
		}

		sendJSONResponse(w, http.StatusBadRequest, common.ClientMessage(err, "The company could not be saved"))
		return
	}

	insertCompanyTracking(r, ownerUID, claims.Email, req, company, action)

	if action == "update" {
		sendJSONResponse(w, http.StatusOK, "Company updated successfully.")
		return
	}

	sendJSONResponse(w, http.StatusOK, "Company created successfully.")
}

func normalizeCompanyRequest(req *CompanyCreateRequest) {
	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.CompanyName = strings.TrimSpace(req.CompanyName)
	req.LegalName = strings.TrimSpace(req.LegalName)
	req.IdentificationType = strings.TrimSpace(req.IdentificationType)
	req.IdentificationNumber = strings.TrimSpace(req.IdentificationNumber)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Website = strings.TrimSpace(req.Website)
	req.CountryCode = strings.TrimSpace(req.CountryCode)
	req.CountryName = strings.TrimSpace(req.CountryName)
	req.StateCode = strings.TrimSpace(req.StateCode)
	req.StateName = strings.TrimSpace(req.StateName)
	req.City = strings.TrimSpace(req.City)
	req.District = strings.TrimSpace(req.District)
	req.PostalCode = strings.TrimSpace(req.PostalCode)
	req.Address1 = strings.TrimSpace(req.Address1)
	req.Address2 = strings.TrimSpace(req.Address2)
	req.Directions = strings.TrimSpace(req.Directions)
	req.Category = strings.TrimSpace(req.Category)
	req.InternalCategory = strings.TrimSpace(req.InternalCategory)
	req.Plan = strings.TrimSpace(req.Plan)
	req.Status = strings.TrimSpace(req.Status)
}

func validateCompanyRequest(req CompanyCreateRequest) error {
	if req.CompanyID != "" {
		if _, err := uuid.Parse(req.CompanyID); err != nil {
			return fmt.Errorf("companyId must be a valid UUID")
		}
	}

	if req.CompanyName == "" {
		return fmt.Errorf("companyName is required")
	}

	if len(req.CompanyName) > 200 {
		return fmt.Errorf("companyName cannot exceed 200 characters")
	}

	if len(req.LegalName) > 200 {
		return fmt.Errorf("legalName cannot exceed 200 characters")
	}

	if req.IdentificationType == "" {
		return fmt.Errorf("identificationType is required")
	}

	if len(req.IdentificationType) > 50 {
		return fmt.Errorf("identificationType cannot exceed 50 characters")
	}

	if req.IdentificationNumber == "" {
		return fmt.Errorf("identificationNumber is required")
	}

	if len(req.IdentificationNumber) > 50 {
		return fmt.Errorf("identificationNumber cannot exceed 50 characters")
	}

	if req.Email != "" {
		if len(req.Email) > 150 {
			return fmt.Errorf("email cannot exceed 150 characters")
		}

		if _, err := mail.ParseAddress(req.Email); err != nil {
			return fmt.Errorf("email is not valid")
		}
	}

	if len(req.Phone) > 30 {
		return fmt.Errorf("phone cannot exceed 30 characters")
	}

	if req.Website != "" {
		if len(req.Website) > 200 {
			return fmt.Errorf("website cannot exceed 200 characters")
		}

		parsedURL, err := url.ParseRequestURI(req.Website)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("website is not valid")
		}
	}

	if len(req.CountryCode) > 10 {
		return fmt.Errorf("countryCode cannot exceed 10 characters")
	}

	if len(req.CountryName) > 100 {
		return fmt.Errorf("countryName cannot exceed 100 characters")
	}

	if len(req.StateCode) > 20 {
		return fmt.Errorf("stateCode cannot exceed 20 characters")
	}

	if len(req.StateName) > 100 {
		return fmt.Errorf("stateName cannot exceed 100 characters")
	}

	if len(req.City) > 100 {
		return fmt.Errorf("city cannot exceed 100 characters")
	}

	if len(req.District) > 100 {
		return fmt.Errorf("district cannot exceed 100 characters")
	}

	if len(req.PostalCode) > 20 {
		return fmt.Errorf("postalCode cannot exceed 20 characters")
	}

	if len(req.Address1) > 200 {
		return fmt.Errorf("address1 cannot exceed 200 characters")
	}

	if len(req.Address2) > 200 {
		return fmt.Errorf("address2 cannot exceed 200 characters")
	}

	if len(req.Category) > 100 {
		return fmt.Errorf("category cannot exceed 100 characters")
	}

	if len(req.InternalCategory) > 100 {
		return fmt.Errorf("internalCategory cannot exceed 100 characters")
	}

	if req.Employees < 0 {
		return fmt.Errorf("employees cannot be negative")
	}

	if len(req.Plan) > 50 {
		return fmt.Errorf("plan cannot exceed 50 characters")
	}

	if len(req.Status) > 50 {
		return fmt.Errorf("status cannot exceed 50 characters")
	}

	return nil
}

func insertCompanyTracking(
	r *http.Request,
	ownerUID string,
	claimEmail string,
	req CompanyCreateRequest,
	company databasecompany.Company,
	action string,
) {
	actionName := "CreateCompany"
	if action == "update" {
		actionName = "UpdateCompany"
	}

	jsonCompany, _ := json.Marshal(company)

	trackingEmail := req.Email
	if strings.TrimSpace(trackingEmail) == "" {
		trackingEmail = claimEmail
	}

	lt := tracking.LoginTracking{
		UserID:         ownerUID,
		Email:          trackingEmail,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIPr(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    req.CountryCode,
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         actionName,
		Jsonstring:     string(jsonCompany),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("❌ Error inserting tracking:", err)
	} else {
		fmt.Println("✅ Tracking inserted successfully")
	}
}

func readClientIPr(r *http.Request) string {
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

func sendJSONResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ResponseMessage{
		Status:  http.StatusText(status),
		Message: message,
	})
}
