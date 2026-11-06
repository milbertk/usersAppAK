package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
)

// BSearchCompanyByName searches companies by partial name
// Returns a list of matching companies
func BSearchCompanyByName(ownerUID, search string) ([]databasecompany.Company, error) {

	// 1) Validate owner
	if strings.TrimSpace(ownerUID) == "" {
		return nil, fmt.Errorf("❌ owner UID is required")
	}

	// 2) Validate search
	if strings.TrimSpace(search) == "" {
		return nil, fmt.Errorf("❌ search text is required")
	}

	// 3) Query database
	companies, err := databasecompany.SearchByName(ownerUID, search)
	if err != nil {
		return nil, fmt.Errorf("❌ error searching companies: %v", err)
	}

	return companies, nil
}
