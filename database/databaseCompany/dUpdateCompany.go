package databasecompany

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

func (c *Company) Update() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		UPDATE public.companies
		SET
			company_name = $1,
			legal_name = $2,
			identification_type = $3,
			identification_number = $4,
			email = $5,
			phone = $6,
			website = $7,
			country_code = $8,
			country_name = $9,
			state_code = $10,
			state_name = $11,
			city = $12,
			district = $13,
			postal_code = $14,
			address1 = $15,
			address2 = $16,
			directions = $17,
			category = $18,
			internal_category = $19,
			employees = $20,
			is_virtual = $21,
			active = $22,
			plan = $23,
			status = $24,
			updated_at = now()
		WHERE company_id = $25
		  AND deleted_at IS NULL
	`

	result, err := db.Exec(query,
		c.CompanyName,
		c.LegalName,
		c.IdentificationType,
		c.IdentificationNumber,
		c.Email,
		c.Phone,
		c.Website,
		c.CountryCode,
		c.CountryName,
		c.StateCode,
		c.StateName,
		c.City,
		c.District,
		c.PostalCode,
		c.Address1,
		c.Address2,
		c.Directions,
		c.Category,
		c.InternalCategory,
		c.Employees,
		c.IsVirtual,
		c.Active,
		c.Plan,
		c.Status,
		c.CompanyID,
	)

	if err != nil {
		return fmt.Errorf("❌ failed to update company: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("⚠️ 2no active company found with ID: %s", c.CompanyID)
	}

	fmt.Println("✅ Company updated successfully:", c.CompanyID)
	return nil
}
