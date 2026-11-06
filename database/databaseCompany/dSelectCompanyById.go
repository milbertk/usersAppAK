package databasecompany

import (
	"database/sql"
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func SelectByCompanyID(companyID string) (*Company, error) {
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
		WHERE company_id = $1
		  AND deleted_at IS NULL
	`

	var c Company
	err = db.QueryRow(query, companyID).Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("⚠️ 1no active company found with ID: %s", companyID)
		}
		return nil, fmt.Errorf("❌ query error: %v", err)
	}

	return &c, nil
}
