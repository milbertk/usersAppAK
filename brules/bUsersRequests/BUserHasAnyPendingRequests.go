package busersrequests

import (
	"fmt"
	"strings"

	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

// BUserHasAnyPendingRequests checks both invited and requested
func BUserHasAnyPendingRequests(userID string) (bool, error) {
	userID = strings.TrimSpace(userID)

	if userID == "" {
		return false, fmt.Errorf("userID is required")
	}

	// Check invited
	invited, err := databaserequestusercompany.UserHasPendingRequests(userID, "invited")
	if err != nil {
		return false, fmt.Errorf("error checking invited requests: %v", err)
	}

	if invited {
		return true, nil
	}

	// Check requested
	requested, err := databaserequestusercompany.UserHasPendingRequests(userID, "requested")
	if err != nil {
		return false, fmt.Errorf("error checking requested requests: %v", err)
	}

	return requested, nil
}
