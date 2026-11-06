package busersrequests

import (
	"fmt"
	"strings"

	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

// BGetPendingInvitationsByUser returns all pending invitations
// where the given user was invited.
func BGetPendingInvitationsByUser(userID string) ([]databaserequestusercompany.CompanyUserRequestWithCompany, error) {
	userID = strings.TrimSpace(userID)

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	requests, err := databaserequestusercompany.SelectPendingByUser(userID, "invited")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending invitations: %v", err)
	}

	return requests, nil
}
