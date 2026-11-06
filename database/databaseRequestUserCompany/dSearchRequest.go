package databaseRequestUserCompany

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SearchByCompanyAndUserText returns requests for one company,
// filtering by invited user display name or email.
func SearchByCompanyAndUserText(companyID, searchText string) ([]CompanyUserRequestSearchResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	searchText = strings.TrimSpace(searchText)

	query := `
		SELECT
			cur.request_id,
			cur.company_id,
			cur.invited_user_id,
			cur.requested_by_user_id,
			cur.role_requested,
			cur.status,
			COALESCE(cur.message, '') AS message,
			cur.created_at,
			cur.updated_at,
			cur.resolved_at,
			COALESCE(u.displayname, '') AS invited_user_name,
			COALESCE(u.email, '') AS invited_user_email
		FROM public.company_user_requests cur
		LEFT JOIN public.usfirebasedata u
			ON u.uid = cur.invited_user_id
		WHERE cur.company_id = $1
		  AND (
				$2 = ''
				OR COALESCE(u.displayname, '') ILIKE '%' || $2 || '%'
				OR COALESCE(u.email, '') ILIKE '%' || $2 || '%'
		  )
		ORDER BY cur.created_at DESC
	`

	rows, err := db.Query(query, companyID, searchText)
	if err != nil {
		return nil, fmt.Errorf("failed to search company user requests: %w", err)
	}
	defer rows.Close()

	var results []CompanyUserRequestSearchResult

	for rows.Next() {
		var item CompanyUserRequestSearchResult

		var updatedAt sql.NullTime
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&item.RequestID,
			&item.CompanyID,
			&item.InvitedUserID,
			&item.RequestedByUserID,
			&item.RoleRequested,
			&item.Status,
			&item.Message,
			&item.CreatedAt,
			&updatedAt,
			&resolvedAt,
			&item.InvitedUserName,
			&item.InvitedUserEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company user request row: %w", err)
		}

		// Handle nullable time fields correctly (*time.Time)
		if updatedAt.Valid {
			t := updatedAt.Time
			item.UpdatedAt = &t
		} else {
			item.UpdatedAt = nil
		}

		if resolvedAt.Valid {
			t := resolvedAt.Time
			item.ResolvedAt = &t
		} else {
			item.ResolvedAt = nil
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
