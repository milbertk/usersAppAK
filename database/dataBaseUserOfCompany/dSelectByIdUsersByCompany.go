package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectByID returns one company user by company_user_id.
func SelectByID(companyUserID string) (*CompanyUser, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	companyUserID = strings.TrimSpace(companyUserID)
	if companyUserID == "" {
		return nil, fmt.Errorf("companyUserId is required")
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
		WHERE company_user_id = $1
	`

	var item CompanyUser
	err = db.QueryRow(query, companyUserID).Scan(
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
		return nil, fmt.Errorf("failed to select company user: %v", err)
	}

	return &item, nil
}
