package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
)

// BCreateOrUpdateCompany decides whether to create or update.
// Returns action = "create" or "update".
func BCreateOrUpdateCompany(company databasecompany.Company) (string, error) {
	company.CompanyID = strings.TrimSpace(company.CompanyID)
	company.CompanyName = strings.TrimSpace(company.CompanyName)
	company.OwnerFirebaseUID = strings.TrimSpace(company.OwnerFirebaseUID)
	company.IdentificationNumber = strings.TrimSpace(company.IdentificationNumber)

	// Only ids in the logs: no names or identification numbers (2026-10-25,
	// security S11).
	fmt.Println("=== BCreateOrUpdateCompany === companyId:", company.CompanyID)

	if company.OwnerFirebaseUID == "" {
		return "", fmt.Errorf("owner uid is required")
	}
	if company.CompanyName == "" {
		return "", fmt.Errorf("company name is required")
	}
	if company.IdentificationNumber == "" {
		return "", fmt.Errorf("identification number is required")
	}
	println("pass 1")
	println(company.CompanyID)

	var existing *databasecompany.Company
	println("pass 1.1id")
	if company.CompanyID != "" {
		println("pass 1.1id2")
		found, err := databasecompany.SelectByCompanyID(company.CompanyID)
		println("pass 1.1id3")
		if err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "no active company found") &&
				!strings.Contains(strings.ToLower(err.Error()), "no rows") {
				return "", fmt.Errorf("error validating company id: %v", err)
			}
		} else {

			println("pass 1.1id4")
			existing = found
		}
	}
	println("pass 1.1")
	byName, err := databasecompany.SearchByName(company.OwnerFirebaseUID, company.CompanyName)
	if err != nil {
		return "", fmt.Errorf("error validating company name: %v", err)
	}
	println("pass 1.1")
	fmt.Println("SearchByName count:", len(byName))

	byIdentification, err := databasecompany.SearchByIdentification(
		company.OwnerFirebaseUID,
		company.IdentificationNumber,
	)
	if err != nil {
		return "", fmt.Errorf("error validating identification number: %v", err)
	}
	println("pass 1.2")
	fmt.Println("SearchByIdentification count:", len(byIdentification))

	if existing != nil {
		if existing.OwnerFirebaseUID != company.OwnerFirebaseUID {
			return "", fmt.Errorf("you are not allowed to update this company")
		}

		// The plan is changed only by AlephK (after a plan request), never
		// from the company form (2026-10-22): keep the saved one.
		company.Plan = existing.Plan
		// Same for the status (2026-10-25, security S3).
		company.Status = existing.Status

		for _, c := range byName {
			if strings.TrimSpace(c.CompanyID) != "" && c.CompanyID != company.CompanyID {
				return "", fmt.Errorf("the company name exists, please choose another")
			}
		}

		for _, c := range byIdentification {
			if strings.TrimSpace(c.CompanyID) != "" && c.CompanyID != company.CompanyID {
				return "", fmt.Errorf("the identification number already exists for this owner")
			}
		}

		if err := company.Update(); err != nil {
			return "", fmt.Errorf("update failed: %v", err)
		}
		println("pass 1.3")
		return "update", nil
	}

	if len(byName) > 0 {
		return "", fmt.Errorf("the company name exists, please choose another")
	}

	if len(byIdentification) > 0 {
		return "", fmt.Errorf("the identification number already exists for this owner")
	}
	println("pass 1.4")
	// New companies always start on the FREE plan and ACTIVE (2026-10-22,
	// status 2026-10-25 security S3).
	company.Plan = "FREE"
	company.Status = "ACTIVE"
	if err := company.Insert(); err != nil {
		return "", fmt.Errorf("insert failed: %v", err)
	}
	println("pass 1.5")
	return "create", nil
}
