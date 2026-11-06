package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func ExistsActiveByCompanyAndUser(companyID string, userID string) (bool, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return false, fmt.Errorf("DB connection error: %v", err)
	}

	companyID = strings.TrimSpace(companyID)
	userID = strings.TrimSpace(userID)

	if companyID == "" {
		return false, fmt.Errorf("companyID is required")
	}

	if userID == "" {
		return false, fmt.Errorf("userID is required")
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM public.company_users
			WHERE company_id = $1
			  AND user_id = $2
			  AND active = true
		)
	`

	var exists bool
	err = db.QueryRow(query, companyID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to validate company user duplication: %v", err)
	}

	return exists, nil
}
