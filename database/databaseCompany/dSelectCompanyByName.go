package databasecompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SearchByName(ownerUID, search string) ([]Company, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		SELECT
			company_id,
			company_name,
			category,
			active,
			status,
			updated_at
		FROM public.companies
		WHERE owner_firebase_uid = $1
		  AND deleted_at IS NULL
		  AND company_name ILIKE '%' || $2 || '%'
		ORDER BY company_name
		LIMIT 50
	`

	rows, err := db.Query(query, ownerUID, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.CompanyID,
			&c.CompanyName,
			&c.Category,
			&c.Active,
			&c.Status,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
