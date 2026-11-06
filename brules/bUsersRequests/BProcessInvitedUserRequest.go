package busersrequests

import (
	"fmt"
	"strings"

	databaseuserofcompany "github.com/milbertk/usersAppAK/database/dataBaseUserOfCompany"
	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

func BProcessInvitedUserRequest(
	requestID string,
	actionCode string,
	message string,
	currentUserID string,
) error {
	requestID = strings.TrimSpace(requestID)
	actionCode = strings.TrimSpace(strings.ToUpper(actionCode))
	message = strings.TrimSpace(message)
	currentUserID = strings.TrimSpace(currentUserID)

	fmt.Println("=== BProcessInvitedUserRequest ===")
	fmt.Println("RequestID:", requestID)
	fmt.Println("ActionCode:", actionCode)
	fmt.Println("CurrentUserID:", currentUserID)

	if requestID == "" {
		return fmt.Errorf("request id is required")
	}

	if currentUserID == "" {
		return fmt.Errorf("current user id is required")
	}

	if actionCode == "" {
		return fmt.Errorf("action code is required")
	}

	switch actionCode {
	case "REJECT", "ACCEPT":
		// valid actions
	default:
		return fmt.Errorf("invalid action code")
	}

	// 1. Get request by ID
	req, err := databaserequestusercompany.SelectByRequestID(requestID)
	if err != nil {
		return fmt.Errorf("error validating request: %v", err)
	}

	if strings.TrimSpace(req.RequestID) == "" {
		return fmt.Errorf("request not found")
	}

	// 2. Only invited user can process the request
	invitedUID := strings.TrimSpace(req.InvitedUserID)
	if invitedUID == "" {
		return fmt.Errorf("invited user id is empty in request")
	}

	if currentUserID != invitedUID {
		return fmt.Errorf("only the invited user can process this request")
	}

	// 3. Only pending requests can be processed
	currentStatus := strings.TrimSpace(strings.ToUpper(req.Status))
	if currentStatus != "PENDING" {
		return fmt.Errorf("only pending requests can be processed")
	}

	// 4. Process action
	switch actionCode {
	case "REJECT":
		err = databaserequestusercompany.UpdateStatusAndMessage(
			requestID,
			"REJECTED",
			message,
		)
		if err != nil {
			return fmt.Errorf("reject process failed: %v", err)
		}

	case "ACCEPT":
		companyID := strings.TrimSpace(req.CompanyID)
		roleRequested := strings.TrimSpace(strings.ToUpper(req.RoleRequested))

		if companyID == "" {
			return fmt.Errorf("company id is empty in request")
		}

		if roleRequested == "" {
			roleRequested = "USER"
		}

		// 4.1 Validate duplicate membership
		exists, err := databaseuserofcompany.ExistsActiveByCompanyAndUser(companyID, invitedUID)
		if err != nil {
			return fmt.Errorf("error validating existing company membership: %v", err)
		}

		if exists {
			return fmt.Errorf("user already belongs to this company")
		}

		// 4.2 Insert user into company_users
		companyUser := &databaseuserofcompany.CompanyUser{
			CompanyID: companyID,
			UserID:    invitedUID,
			Role:      roleRequested,
		}

		err = companyUser.Insert()
		if err != nil {
			return fmt.Errorf("accept process failed inserting company user: %v", err)
		}

		// 4.3 Update request status
		err = databaserequestusercompany.UpdateStatusAndMessage(
			requestID,
			"APPROVED",
			message,
		)
		if err != nil {
			return fmt.Errorf("accept process failed updating request status: %v", err)
		}
	}

	return nil
}
