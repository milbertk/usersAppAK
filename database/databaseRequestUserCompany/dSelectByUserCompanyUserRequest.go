package databaseRequestUserCompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectPendingByUser returns all active pending requests for a user.
// filterBy can be: "invited" or "requested"
func SelectPendingByUser(userID, filterBy string) ([]CompanyUserRequestWithCompany, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	userID = strings.TrimSpace(userID)
	filterBy = strings.TrimSpace(strings.ToLower(filterBy))

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	var userField string

	switch filterBy {
	case "invited":
		userField = "cur.invited_user_id"
	case "requested":
		userField = "cur.requested_by_user_id"
	default:
		return nil, fmt.Errorf("filterBy must be 'invited' or 'requested'")
	}

	query := fmt.Sprintf(`
		SELECT
			cur.request_id,
			cur.company_id,
			COALESCE(c.company_name, '') AS company_name,
			cur.invited_user_id,
			cur.requested_by_user_id,
			cur.role_requested,
			cur.status,
			COALESCE(cur.message, '') AS message,
			cur.created_at,
			cur.updated_at,
			cur.resolved_at
		FROM public.company_user_requests cur
		INNER JOIN public.companies c
			ON c.company_id = cur.company_id
		WHERE %s = $1
		  AND cur.status = 'PENDING'
		ORDER BY cur.created_at DESC
	`, userField)

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending requests: %v", err)
	}
	defer rows.Close()

	var results []CompanyUserRequestWithCompany

	for rows.Next() {
		var req CompanyUserRequestWithCompany

		err := rows.Scan(
			&req.RequestID,
			&req.CompanyID,
			&req.CompanyName,
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
			return nil, fmt.Errorf("failed to scan pending request row: %v", err)
		}

		results = append(results, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
