package busersrequests

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
	databasecompanyuserrequest "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

func BCreateCompanyUserRequest(
	req databasecompanyuserrequest.CompanyUserRequest,
) error {

	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.InvitedUserID = strings.TrimSpace(req.InvitedUserID)
	req.RequestedByUserID = strings.TrimSpace(req.RequestedByUserID)
	req.RoleRequested = strings.TrimSpace(req.RoleRequested)
	req.Status = strings.TrimSpace(req.Status)
	req.Message = strings.TrimSpace(req.Message)

	fmt.Println("=== BCreateCompanyUserRequest ===")
	fmt.Println("CompanyID:", req.CompanyID)
	fmt.Println("InvitedUserID:", req.InvitedUserID)
	fmt.Println("RequestedByUserID:", req.RequestedByUserID)
	fmt.Println("RoleRequested:", req.RoleRequested)
	fmt.Println("Status:", req.Status)

	if req.CompanyID == "" {
		return fmt.Errorf("company id is required")
	}

	if req.InvitedUserID == "" {
		return fmt.Errorf("invited user id is required")
	}

	if req.RequestedByUserID == "" {
		return fmt.Errorf("requested by user id is required")
	}

	if req.RoleRequested == "" {
		req.RoleRequested = "USER"
	}

	if req.Status == "" {
		req.Status = "PENDING"
	}

	if req.InvitedUserID == req.RequestedByUserID {
		return fmt.Errorf("you cannot create a request for yourself")
	}

	// 1. Validate company exists
	company, err := databasecompany.SelectByCompanyID(req.CompanyID)
	if err != nil {
		return fmt.Errorf("error validating company: %v", err)
	}

	// 2. Validate creator is the owner of the company
	if strings.TrimSpace(company.OwnerFirebaseUID) != req.RequestedByUserID {
		return fmt.Errorf("only the company owner can create this request")
	}

	// 3. Validate there is no other pending request for the same user in the same company
	existsPending, err := databasecompanyuserrequest.ExistsPendingRequest(
		req.CompanyID,
		req.InvitedUserID,
	)
	if err != nil {
		return fmt.Errorf("error validating pending request: %v", err)
	}

	if existsPending {
		return fmt.Errorf("there is already a pending request for this user in this company")
	}

	// 4. Insert request
	if err := req.Insert(); err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	return nil
}
