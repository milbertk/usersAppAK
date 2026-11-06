package buserofcompany

import (
	"fmt"
	"strings"

	databaseuserofcompany "github.com/milbertk/usersAppAK/database/dataBaseUserOfCompany"
)

type UpdateCompanyUserInput struct {
	CompanyUserID string `json:"companyUserId"`
	CompanyID     string `json:"companyId"`
	UserID        string `json:"userId"`
	Role          string `json:"role"`
	Active        bool   `json:"active"`
}

type UpdateCompanyUserResult struct {
	Status          string                             `json:"status"`
	Message         string                             `json:"message"`
	RequesterUserID string                             `json:"requesterUserId"`
	RequesterRole   string                             `json:"requesterRole"`
	CompanyUser     *databaseuserofcompany.CompanyUser `json:"companyUser"`
}

func BUpdateCompanyUser(input UpdateCompanyUserInput, requesterUserID string) (*UpdateCompanyUserResult, error) {
	input.CompanyUserID = strings.TrimSpace(input.CompanyUserID)
	input.CompanyID = strings.TrimSpace(input.CompanyID)
	input.UserID = strings.TrimSpace(input.UserID)
	input.Role = strings.ToUpper(strings.TrimSpace(input.Role))
	requesterUserID = strings.TrimSpace(requesterUserID)

	if input.CompanyUserID == "" {
		return nil, fmt.Errorf("companyUserId is required")
	}

	if input.CompanyID == "" {
		return nil, fmt.Errorf("companyId is required")
	}

	if input.UserID == "" {
		return nil, fmt.Errorf("userId is required")
	}

	if requesterUserID == "" {
		return nil, fmt.Errorf("requester user id is required")
	}

	if input.Role == "" {
		input.Role = "USER"
	}

	if !isValidCompanyRole(input.Role) {
		return nil, fmt.Errorf("invalid role: %s", input.Role)
	}

	hasAccess, requesterRole, err := databaseuserofcompany.UserHasManagementAccess(
		input.CompanyID,
		requesterUserID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to validate requester access: %v", err)
	}

	if !hasAccess {
		return nil, fmt.Errorf("unauthorized: only OWNER or MANAGER can update company users")
	}

	item := &databaseuserofcompany.CompanyUser{
		CompanyUserID: input.CompanyUserID,
		CompanyID:     input.CompanyID,
		UserID:        input.UserID,
		Role:          input.Role,
		Active:        input.Active,
	}

	err = item.Update()
	if err != nil {
		return nil, fmt.Errorf("failed to update company user: %v", err)
	}

	return &UpdateCompanyUserResult{
		Status:          "OK",
		Message:         "company user updated successfully",
		RequesterUserID: requesterUserID,
		RequesterRole:   requesterRole,
		CompanyUser:     item,
	}, nil
}

func isValidCompanyRole(role string) bool {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "OWNER", "ADMIN", "MANAGER", "USER", "SUPPORT":
		return true
	default:
		return false
	}
}
