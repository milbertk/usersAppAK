package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
)

// BSelectAllCompanies retrieves all companies for an owner
func BSelectAllCompanies(ownerUID string) ([]databasecompany.Company, error) {

	// 1) Validate owner
	if strings.TrimSpace(ownerUID) == "" {
		return nil, fmt.Errorf("❌ owner UID is required")
	}

	// 2) Load companies from database
	companies, err := databasecompany.SelectAllByOwner(ownerUID)
	if err != nil {
		return nil, fmt.Errorf("❌ error loading companies: %v", err)
	}

	return companies, nil
}
