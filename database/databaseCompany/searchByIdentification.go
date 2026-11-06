package databasecompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SearchByIdentification(ownerUID string, identification string) ([]Company, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
	}

	query := `
		SELECT
			company_id,
			owner_firebase_uid,
			company_name,
			identification_number
		FROM public.companies
		WHERE owner_firebase_uid = $1
		  AND identification_number = $2
		  AND deleted_at IS NULL
	`

	rows, err := db.Query(query, ownerUID, identification)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.CompanyID,
			&c.OwnerFirebaseUID,
			&c.CompanyName,
			&c.IdentificationNumber,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		companies = append(companies, c)
	}

	return companies, nil
}
