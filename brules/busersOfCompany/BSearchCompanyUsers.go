package buserofcompany

import (
	"fmt"
	"strings"

	databaseuserofcompany "github.com/milbertk/usersAppAK/database/dataBaseUserOfCompany"
)

type SearchCompanyUsersResult struct {
	CompanyID       string                                          `json:"companyId"`
	RequesterUserID string                                          `json:"requesterUserId"`
	RequesterRole   string                                          `json:"requesterRole"`
	IncludeInactive bool                                            `json:"includeInactive"`
	Users           []databaseuserofcompany.CompanyUserSearchResult `json:"users"`
	Total           int                                             `json:"total"`
}

func BSearchCompanyUsers(companyID string, requesterUserID string, includeInactive bool) (*SearchCompanyUsersResult, error) {
	companyID = strings.TrimSpace(companyID)
	requesterUserID = strings.TrimSpace(requesterUserID)

	if companyID == "" {
		return nil, fmt.Errorf("companyId is required")
	}
	if requesterUserID == "" {
		return nil, fmt.Errorf("requester user id is required")
	}

	println("validation2")
	println(companyID, requesterUserID)
	hasAccess, role, err := databaseuserofcompany.UserHasManagementAccess(companyID, requesterUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate requester access: %v", err)
	}

	if !hasAccess {
		return nil, fmt.Errorf("unauthorized: only OWNER or MANAGER can search company users")
	}

	users, err := databaseuserofcompany.SearchByCompany(companyID, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("failed to search company users: %v", err)
	}

	result := &SearchCompanyUsersResult{
		CompanyID:       companyID,
		RequesterUserID: requesterUserID,
		RequesterRole:   role,
		IncludeInactive: includeInactive,
		Users:           users,
		Total:           len(users),
	}

	return result, nil
}
