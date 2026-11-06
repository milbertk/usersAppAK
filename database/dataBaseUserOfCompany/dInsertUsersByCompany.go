package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert inserts a new company user into company_users.
func (r *CompanyUser) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	r.CompanyID = strings.TrimSpace(r.CompanyID)
	r.UserID = strings.TrimSpace(r.UserID)
	r.Role = strings.ToUpper(strings.TrimSpace(r.Role))

	if r.Role == "" {
		r.Role = "USER"
	}

	query := `
		INSERT INTO public.company_users (
			company_user_id,
			company_id,
			user_id,
			role,
			active,
			joined_at,
			removed_at,
			created_at,
			updated_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			true,
			NOW(),
			NULL,
			NOW(),
			NULL
		)
		RETURNING 
			company_user_id,
			active,
			joined_at,
			created_at
	`

	err = db.QueryRow(
		query,
		r.CompanyID,
		r.UserID,
		r.Role,
	).Scan(
		&r.CompanyUserID,
		&r.Active,
		&r.JoinedAt,
		&r.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert company user: %v", err)
	}

	return nil
}
