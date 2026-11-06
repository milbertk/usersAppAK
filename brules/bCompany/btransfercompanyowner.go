package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

type TransferCompanyOwnerRequest struct {
	CompanyID   string `json:"companyId"`
	NewOwnerUID string `json:"newOwnerUid"`
	LoggedUID   string `json:"-"`
}

type TransferCompanyOwnerResponse struct {
	CompanyID        string `json:"companyId"`
	PreviousOwnerUID string `json:"previousOwnerUid"`
	NewOwnerUID      string `json:"newOwnerUid"`
}

func BTransferCompanyOwner(req TransferCompanyOwnerRequest) (*TransferCompanyOwnerResponse, error) {
	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.NewOwnerUID = strings.TrimSpace(req.NewOwnerUID)
	req.LoggedUID = strings.TrimSpace(req.LoggedUID)

	if req.CompanyID == "" {
		return nil, fmt.Errorf("companyId is required")
	}

	if req.NewOwnerUID == "" {
		return nil, fmt.Errorf("newOwnerUid is required")
	}

	if req.LoggedUID == "" {
		return nil, fmt.Errorf("logged user is required")
	}

	if req.NewOwnerUID == req.LoggedUID {
		return nil, fmt.Errorf("the new owner must be different from the logged user")
	}

	company, err := databasecompany.SelectByCompanyID(req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("error loading company: %w", err)
	}

	currentOwnerUID := strings.TrimSpace(company.OwnerFirebaseUID)

	if currentOwnerUID == "" {
		return nil, fmt.Errorf("the company does not have a current owner")
	}

	if req.NewOwnerUID == currentOwnerUID {
		return nil, fmt.Errorf("the new owner must be different from the current owner")
	}

	if req.LoggedUID != currentOwnerUID {
		return nil, fmt.Errorf("only the current owner can transfer this company")
	}

	newOwner, err := databaseusers.SelectByUID(req.NewOwnerUID)
	if err != nil {
		return nil, fmt.Errorf("error loading new owner: %w", err)
	}

	if strings.TrimSpace(newOwner.UID) == "" {
		return nil, fmt.Errorf("new owner does not exist")
	}

	if strings.ToLower(strings.TrimSpace(newOwner.Active)) != "true" {
		return nil, fmt.Errorf("new owner is not active")
	}

	if strings.ToLower(strings.TrimSpace(newOwner.Disabled)) == "true" {
		return nil, fmt.Errorf("new owner is disabled")
	}

	company.CompanyID = req.CompanyID
	company.OwnerFirebaseUID = req.NewOwnerUID

	if err := company.Update(); err != nil {
		return nil, fmt.Errorf("error updating company owner: %w", err)
	}

	return &TransferCompanyOwnerResponse{
		CompanyID:        req.CompanyID,
		PreviousOwnerUID: currentOwnerUID,
		NewOwnerUID:      req.NewOwnerUID,
	}, nil
}
