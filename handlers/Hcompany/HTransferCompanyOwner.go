package hcompany

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/milbertk/security"
	"github.com/milbertk/tracking"

	bCreateCompany "github.com/milbertk/usersAppAK/brules/bCompany"

	"github.com/milbertk/usersAppAK/handlers/common"
)

type TransferCompanyOwnerRequest struct {
	CompanyID   string `json:"companyId"`
	NewOwnerUID string `json:"newOwnerUid"`
}

type TransferCompanyOwnerHTTPResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HTransferCompanyOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		sendTransferCompanyOwnerResponse(
			w,
			http.StatusMethodNotAllowed,
			"Use POST or PUT",
			nil,
		)
		return
	}

	claims, ok := security.GetClaims(r)
	if !ok || claims == nil || strings.TrimSpace(claims.UID) == "" {
		sendTransferCompanyOwnerResponse(
			w,
			http.StatusUnauthorized,
			"Missing user claims",
			nil,
		)
		return
	}

	loggedUID := strings.TrimSpace(claims.UID)

	var req TransferCompanyOwnerRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendTransferCompanyOwnerResponse(
			w,
			http.StatusBadRequest,
			common.InvalidJSONMessage(err),
			nil,
		)
		return
	}

	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.NewOwnerUID = strings.TrimSpace(req.NewOwnerUID)

	result, err := bCreateCompany.BTransferCompanyOwner(
		bCreateCompany.TransferCompanyOwnerRequest{
			CompanyID:   req.CompanyID,
			NewOwnerUID: req.NewOwnerUID,
			LoggedUID:   loggedUID,
		},
	)

	if err != nil {
		insertTransferCompanyOwnerTracking(
			r,
			loggedUID,
			claims.Email,
			req.CompanyID,
			req.NewOwnerUID,
			"",
			false,
			err.Error(),
		)

		sendTransferCompanyOwnerResponse(
			w,
			http.StatusBadRequest,
			common.ClientMessage(err, "The company could not be transferred"),
			nil,
		)
		return
	}

	insertTransferCompanyOwnerTracking(
		r,
		loggedUID,
		claims.Email,
		req.CompanyID,
		req.NewOwnerUID,
		result.PreviousOwnerUID,
		true,
		"Company owner updated successfully.",
	)

	sendTransferCompanyOwnerResponse(
		w,
		http.StatusOK,
		"Company owner updated successfully.",
		result,
	)
}

func insertTransferCompanyOwnerTracking(
	r *http.Request,
	loggedUID string,
	userEmail string,
	companyID string,
	newOwnerUID string,
	previousOwnerUID string,
	success bool,
	message string,
) {
	jsonTransfer, _ := json.Marshal(map[string]interface{}{
		"companyId":        companyID,
		"loggedUid":        loggedUID,
		"previousOwnerUid": previousOwnerUID,
		"newOwnerUid":      newOwnerUID,
		"success":          success,
		"message":          message,
	})

	lt := tracking.LoginTracking{
		UserID:         loggedUID,
		Email:          userEmail,
		DateLocalAcces: time.Now().Format("2006-01-02 15:04:05"),
		IP:             readClientIP(r),
		Platform:       "UserManagement",
		MacAddress:     "0",
		Browser:        r.UserAgent(),
		CountryCode:    "",
		GMTTime:        time.Now().UTC().Format("2006-01-02 15:04:05"),
		Lang:           r.Header.Get("Accept-Language"),
		Action:         "TransferCompanyOwner",
		Jsonstring:     string(jsonTransfer),
	}

	if err := lt.Insert(); err != nil {
		fmt.Println("Error inserting tracking:", err)
	} else {
		fmt.Println("Tracking inserted successfully")
	}
}

func sendTransferCompanyOwnerResponse(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(TransferCompanyOwnerHTTPResponse{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	})
}
