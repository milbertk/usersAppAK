package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
)

func BUpdateCompany(company databasecompany.Company) error {

	// 1) Validate company_id
	if strings.TrimSpace(company.CompanyID) == "" {
		return fmt.Errorf("❌ company_id is required")
	}

	// 2) Load current company (exists check)
	current, err := databasecompany.SelectByCompanyID(company.CompanyID)
	if err != nil {
		return fmt.Errorf("❌ error validating the company: %v", err)
	}

	// 3) Security: owner must match
	if current.OwnerFirebaseUID != company.OwnerFirebaseUID {
		return fmt.Errorf("❌ you don't have permission to update this company")
	}

	// The plan and the status are changed only by AlephK, never by the owner
	// (2026-10-22; status 2026-10-25 security S3).
	company.Plan = current.Plan
	company.Status = current.Status

	// 4) If name changed, validate uniqueness (same owner)
	if strings.TrimSpace(company.CompanyName) != "" &&
		!strings.EqualFold(strings.TrimSpace(company.CompanyName), strings.TrimSpace(current.CompanyName)) {

		companyArray, err := databasecompany.SearchByName(company.OwnerFirebaseUID, company.CompanyName)
		if err != nil {
			return fmt.Errorf("❌ error validating the company name: %v", err)
		}

		// If SearchByName returns slice:
		if len(companyArray) > 0 {
			return fmt.Errorf("❌ 3the company name exists, please choose another")
		}
	}

	// 5) Update in DB
	if err := company.Update(); err != nil {
		return fmt.Errorf("❌ update failed: %v", err)
	}

	return nil
}
