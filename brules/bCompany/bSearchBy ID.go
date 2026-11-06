package bCreateCompany

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
)

// BSelectCompanyByID retrieves a company by company_id
func BSelectCompanyByID(companyID string) (*databasecompany.Company, error) {

	// 1) Validate input
	if strings.TrimSpace(companyID) == "" {
		return nil, fmt.Errorf("❌ company_id is required")
	}

	// 2) Load company from database
	company, err := databasecompany.SelectByCompanyID(companyID)
	if err != nil {
		return nil, fmt.Errorf("❌ error loading company: %v", err)
	}

	return company, nil
}
