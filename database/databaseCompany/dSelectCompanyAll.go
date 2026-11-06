package databasecompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SelectAllByOwner(ownerUID string) ([]Company, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		SELECT
			company_id,
			owner_firebase_uid,
			company_name,
			legal_name,
			identification_type,
			identification_number,
			email,
			phone,
			website,
			country_code,
			country_name,
			state_code,
			state_name,
			city,
			district,
			postal_code,
			address1,
			address2,
			directions,
			category,
			internal_category,
			employees,
			is_virtual,
			active,
			plan,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM public.companies
		WHERE owner_firebase_uid = $1
		  AND deleted_at IS NULL
		ORDER BY company_name
	`

	rows, err := db.Query(query, ownerUID)
	if err != nil {
		return nil, fmt.Errorf("❌ query error: %v", err)
	}
	defer rows.Close()

	var companies []Company

	for rows.Next() {
		var c Company

		err := rows.Scan(
			&c.CompanyID,
			&c.OwnerFirebaseUID,
			&c.CompanyName,
			&c.LegalName,
			&c.IdentificationType,
			&c.IdentificationNumber,
			&c.Email,
			&c.Phone,
			&c.Website,
			&c.CountryCode,
			&c.CountryName,
			&c.StateCode,
			&c.StateName,
			&c.City,
			&c.District,
			&c.PostalCode,
			&c.Address1,
			&c.Address2,
			&c.Directions,
			&c.Category,
			&c.InternalCategory,
			&c.Employees,
			&c.IsVirtual,
			&c.Active,
			&c.Plan,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("❌ scan error: %v", err)
		}

		companies = append(companies, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("❌ rows error: %v", err)
	}

	return companies, nil
}
