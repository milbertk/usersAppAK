package databaseRequestUserCompany

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SelectByRequestID(requestID string) (CompanyUserRequest, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return CompanyUserRequest{}, fmt.Errorf("DB connection error: %v", err)
	}

	requestID = strings.TrimSpace(requestID)

	query := `
		SELECT
			request_id,
			company_id,
			invited_user_id,
			requested_by_user_id,
			role_requested,
			status,
			COALESCE(message, ''),
			created_at,
			updated_at,
			resolved_at
		FROM public.company_user_requests
		WHERE request_id = $1
		LIMIT 1
	`

	var req CompanyUserRequest

	err = db.QueryRow(query, requestID).Scan(
		&req.RequestID,
		&req.CompanyID,
		&req.InvitedUserID,
		&req.RequestedByUserID,
		&req.RoleRequested,
		&req.Status,
		&req.Message,
		&req.CreatedAt,
		&req.UpdatedAt,
		&req.ResolvedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return CompanyUserRequest{}, fmt.Errorf("request not found")
		}
		return CompanyUserRequest{}, fmt.Errorf("failed to get request: %v", err)
	}

	return req, nil
}
