package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SearchByCompany returns company users for a given company_id.
func SearchByCompany(companyID string, includeInactive bool) ([]CompanyUserSearchResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	companyID = strings.TrimSpace(companyID)
	if companyID == "" {
		return nil, fmt.Errorf("companyId is required")
	}

	query := `
		SELECT
			cu.company_user_id,
			cu.company_id,
			cu.user_id,
			COALESCE(ufd.displayname, '') AS user_name,
			COALESCE(ufd.email, '') AS email,
			cu.role,
			cu.active,
			TO_CHAR(cu.joined_at, 'YYYY-MM-DD"T"HH24:MI:SS') AS joined_at,
			CASE
				WHEN cu.removed_at IS NOT NULL THEN TO_CHAR(cu.removed_at, 'YYYY-MM-DD"T"HH24:MI:SS')
				ELSE NULL
			END AS removed_at,
			TO_CHAR(cu.created_at, 'YYYY-MM-DD"T"HH24:MI:SS') AS created_at,
			TO_CHAR(cu.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS') AS updated_at
		FROM public.company_users cu
		LEFT JOIN public.usfirebasedata ufd
			ON ufd.uid = cu.user_id
		WHERE cu.company_id = $1
	`

	if !includeInactive {
		query += ` AND cu.active = true`
	}

	query += ` ORDER BY cu.joined_at DESC`

	rows, err := db.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to search company users: %v", err)
	}
	defer rows.Close()

	var results []CompanyUserSearchResult

	for rows.Next() {
		var item CompanyUserSearchResult

		err := rows.Scan(
			&item.CompanyUserID,
			&item.CompanyID,
			&item.UserID,
			&item.UserName,
			&item.Email,
			&item.Role,
			&item.Active,
			&item.JoinedAt,
			&item.RemovedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company user row: %v", err)
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
