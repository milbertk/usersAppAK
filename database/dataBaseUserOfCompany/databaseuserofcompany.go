package databaseuserofcompany

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectActiveByCompanyAndUser returns an active company user record
// for the given company_id and user_id.
// If no record exists, it returns (nil, nil).
func SelectActiveByCompanyAndUser(companyID string, userID string) (*CompanyUser, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	companyID = strings.TrimSpace(companyID)
	userID = strings.TrimSpace(userID)

	if companyID == "" {
		return nil, fmt.Errorf("companyID is required")
	}

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	query := `
		SELECT
			company_user_id,
			company_id,
			user_id,
			role,
			active,
			joined_at,
			removed_at,
			created_at,
			updated_at
		FROM public.company_users
		WHERE company_id = $1
		  AND user_id = $2
		  AND active = true
		LIMIT 1
	`

	var item CompanyUser
	err = db.QueryRow(query, companyID, userID).Scan(
		&item.CompanyUserID,
		&item.CompanyID,
		&item.UserID,
		&item.Role,
		&item.Active,
		&item.JoinedAt,
		&item.RemovedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to validate company user duplication: %v", err)
	}

	return &item, nil
}
