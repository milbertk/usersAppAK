package databasecompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert inserts a new company with safe defaults
func (c *Company) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	if strings.TrimSpace(c.OwnerFirebaseUID) == "" {
		return fmt.Errorf("owner UID is required")
	}

	if strings.TrimSpace(c.CompanyName) == "" {
		return fmt.Errorf("company name is required")
	}

	if c.Status == "" {
		c.Status = "ACTIVE"
	}

	if c.Plan == "" {
		c.Plan = "FREE"
	}

	if c.Employees < 0 {
		c.Employees = 0
	}

	query := `
		INSERT INTO public.companies (
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
			status
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,
			$9,$10,
			$11,$12,
			$13,$14,$15,
			$16,$17,$18,
			$19,$20,$21,$22,$23,$24,$25
		)
		RETURNING company_id, created_at, updated_at
	`

	err = db.QueryRow(
		query,
		c.OwnerFirebaseUID,
		strings.TrimSpace(c.CompanyName),
		strings.TrimSpace(c.LegalName),
		c.IdentificationType,
		c.IdentificationNumber,
		strings.TrimSpace(c.Email),
		strings.TrimSpace(c.Phone),
		strings.TrimSpace(c.Website),
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
	).Scan(
		&c.CompanyID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		fmt.Println("❌ failed to insert company:", err)
		return fmt.Errorf("❌ failed to insert company: %v", err)
	}

	fmt.Println("✅ Company inserted successfully:", c.CompanyID)
	return nil
}
