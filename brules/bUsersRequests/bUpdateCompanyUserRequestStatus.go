package busersrequests

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

func BUpdateCompanyUserRequestStatus(
	requestID string,
	status string,
	message string,
	currentUserID string,
) error {
	requestID = strings.TrimSpace(requestID)
	status = strings.TrimSpace(strings.ToUpper(status))
	message = strings.TrimSpace(message)
	currentUserID = strings.TrimSpace(currentUserID)

	fmt.Println("=== BUpdateCompanyUserRequestStatus ===")
	fmt.Println("RequestID:", requestID)
	fmt.Println("Status:", status)
	fmt.Println("CurrentUserID:", currentUserID)

	if requestID == "" {
		return fmt.Errorf("request id is required")
	}

	if currentUserID == "" {
		return fmt.Errorf("current user id is required")
	}

	if status == "" {
		return fmt.Errorf("status is required")
	}

	switch status {
	case "PENDING", "APPROVED", "REJECTED", "EXPIRED", "CANCELED":
	default:
		return fmt.Errorf("invalid status")
	}

	// 1. Get request by ID
	req, err := databaserequestusercompany.SelectByRequestID(requestID)
	if err != nil {
		return fmt.Errorf("error validating request: %v", err)
	}

	if strings.TrimSpace(req.RequestID) == "" {
		return fmt.Errorf("request not found")
	}

	// 2. Validate company exists
	company, err := databasecompany.SelectByCompanyID(req.CompanyID)
	if err != nil {
		return fmt.Errorf("error validating company: %v", err)
	}

	ownerUID := strings.TrimSpace(company.OwnerFirebaseUID)
	invitedUID := strings.TrimSpace(req.InvitedUserID)

	// 3. Validate permissions
	if currentUserID != ownerUID && currentUserID != invitedUID {
		return fmt.Errorf("only the company owner or the invited user can update this request" + "|" + ownerUID + "|" + invitedUID + "|" + currentUserID)
	}

	// 4. Update request
	err = databaserequestusercompany.UpdateStatusAndMessage(requestID, status, message)
	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
}
