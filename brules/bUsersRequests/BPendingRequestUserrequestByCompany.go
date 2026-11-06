package busersrequests

import (
	"fmt"
	"strings"

	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

// BPendingRequestUserrequestByCompany returns a flag indicating
// whether the current user has pending request notifications.
func BPendingRequestUserrequestByCompany(currentUserID, filterBy string) (bool, error) {
	currentUserID = strings.TrimSpace(currentUserID)
	filterBy = strings.TrimSpace(strings.ToLower(filterBy))

	fmt.Println("=== BPendingRequestUserrequestByCompany ===")
	fmt.Println("CurrentUserID:", currentUserID)
	fmt.Println("FilterBy:", filterBy)

	if currentUserID == "" {
		return false, fmt.Errorf("current user id is required")
	}

	switch filterBy {
	case "invited", "requested":
	default:
		return false, fmt.Errorf("filterBy must be 'invited' or 'requested'")
	}

	exists, err := databaserequestusercompany.UserHasPendingRequests(currentUserID, filterBy)
	if err != nil {
		return false, fmt.Errorf("error validating pending requests: %v", err)
	}

	return exists, nil
}
