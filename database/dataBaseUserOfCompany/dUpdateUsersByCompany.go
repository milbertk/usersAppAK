package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Update updates role and active status of a company user.
// Uses composite WHERE for extra protection:
// company_user_id + company_id + user_id
func (r *CompanyUser) Update() error {

	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	r.CompanyUserID = strings.TrimSpace(r.CompanyUserID)
	r.CompanyID = strings.TrimSpace(r.CompanyID)
	r.UserID = strings.TrimSpace(r.UserID)
	r.Role = strings.ToUpper(strings.TrimSpace(r.Role))

	if r.CompanyUserID == "" {
		return fmt.Errorf("companyUserId is required")
	}

	if r.CompanyID == "" {
		return fmt.Errorf("companyId is required")
	}

	if r.UserID == "" {
		return fmt.Errorf("userId is required")
	}

	if r.Role == "" {
		r.Role = "USER"
	}

	query := `
		UPDATE public.company_users
		SET
			role = $4,
			active = $5,
			updated_at = NOW(),

			removed_at =
				CASE
					WHEN $5 = false THEN NOW()
					ELSE NULL
				END

		WHERE
			company_user_id = $1
			AND company_id = $2
			AND user_id = $3

		RETURNING
			updated_at,
			removed_at
	`

	err = db.QueryRow(
		query,
		r.CompanyUserID, // $1
		r.CompanyID,     // $2
		r.UserID,        // $3
		r.Role,          // $4
		r.Active,        // $5
	).Scan(
		&r.UpdatedAt,
		&r.RemovedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to update company user (company_user_id=%s company_id=%s user_id=%s): %v",
			r.CompanyUserID,
			r.CompanyID,
			r.UserID,
			err,
		)
	}

	return nil
}
