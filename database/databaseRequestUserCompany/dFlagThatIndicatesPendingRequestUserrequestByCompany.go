package databaseRequestUserCompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// UserHasPendingRequests returns true if the user has at least one pending request.
// filterBy can be: "invited" or "requested"
func UserHasPendingRequests(userID, filterBy string) (bool, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return false, fmt.Errorf("DB connection error: %v", err)
	}

	userID = strings.TrimSpace(userID)
	filterBy = strings.TrimSpace(strings.ToLower(filterBy))

	if userID == "" {
		return false, fmt.Errorf("userID is required")
	}

	var userField string

	switch filterBy {
	case "invited":
		userField = "invited_user_id"
	case "requested":
		userField = "requested_by_user_id"
	default:
		return false, fmt.Errorf("filterBy must be 'invited' or 'requested'")
	}

	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.company_user_requests
			WHERE %s = $1
			  AND status = 'PENDING'
		)
	`, userField)

	var exists bool

	err = db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to verify pending requests: %v", err)
	}

	return exists, nil
}
