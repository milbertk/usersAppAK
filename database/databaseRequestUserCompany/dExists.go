package databaseRequestUserCompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func ExistsPendingRequest(companyID, invitedUserID string) (bool, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return false, fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM public.company_user_requests
			WHERE company_id = $1
			  AND invited_user_id = $2
			  AND status = 'PENDING'
		)
	`

	var exists bool
	err = db.QueryRow(query, companyID, invitedUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to validate pending request: %v", err)
	}

	return exists, nil
}
