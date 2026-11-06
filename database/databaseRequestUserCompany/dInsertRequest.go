package databaseRequestUserCompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert inserts a new company user request into company_user_requests.
func (r *CompanyUserRequest) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		INSERT INTO public.company_user_requests (
			request_id,
			company_id,
			invited_user_id,
			requested_by_user_id,
			role_requested,
			status,
			message,
			created_at,
			updated_at,
			resolved_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			NOW(),
			NULL,
			NULL
		)
		RETURNING request_id, created_at
	`

	err = db.QueryRow(
		query,
		r.CompanyID,
		r.InvitedUserID,
		r.RequestedByUserID,
		r.RoleRequested,
		r.Status,
		r.Message,
	).Scan(&r.RequestID, &r.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert company user request: %v", err)
	}

	return nil
}
