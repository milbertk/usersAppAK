package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SearchCompaniesByUserID(userID string) ([]CompanyByUserResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	userID = strings.TrimSpace(userID)

	if userID == "" {
		return nil, fmt.Errorf("userId is required")
	}

	query := `
		SELECT
			c.company_id,
			c.company_name,
			c.owner_firebase_uid AS user_id,
			'OWNER' AS kind
		FROM public.companies AS c
		WHERE
			c.status = 'ACTIVE'
			AND c.owner_firebase_uid = $1

		UNION ALL

		SELECT
			c.company_id,
			c.company_name,
			cu.user_id,
			cu.role AS kind
		FROM public.companies AS c
		INNER JOIN public.company_users AS cu
			ON c.company_id = cu.company_id
		WHERE
			c.status = 'ACTIVE'
			AND cu.active = true
			AND cu.user_id = $1
			AND c.owner_firebase_uid <> $1

		ORDER BY company_name ASC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to search companies by user id: %v", err)
	}
	defer rows.Close()

	var results []CompanyByUserResult

	for rows.Next() {
		var item CompanyByUserResult

		err := rows.Scan(
			&item.CompanyID,
			&item.CompanyName,
			&item.UserID,
			&item.Kind,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company row: %v", err)
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
