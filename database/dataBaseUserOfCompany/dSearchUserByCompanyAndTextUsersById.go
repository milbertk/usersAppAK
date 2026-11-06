package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SearchByCompanyAndUserText returns company users for a company filtered by user_id text.
func SearchByCompanyAndUserText(companyID, searchText string, includeInactive bool) ([]CompanyUserSearchResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	companyID = strings.TrimSpace(companyID)
	searchText = strings.TrimSpace(searchText)

	if companyID == "" {
		return nil, fmt.Errorf("companyId is required")
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
	`

	args := []interface{}{companyID}
	argPos := 2

	if !includeInactive {
		query += ` AND active = true`
	}

	if searchText != "" {
		query += fmt.Sprintf(` AND user_id ILIKE $%d`, argPos)
		args = append(args, "%"+searchText+"%")
		argPos++
	}

	query += ` ORDER BY joined_at DESC`

	rows, err := db.Query(query, args...)
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
