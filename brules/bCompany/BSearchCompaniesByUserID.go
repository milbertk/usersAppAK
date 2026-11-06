package bCreateCompany

import (
	"fmt"
	"strings"

	databaseuserofcompany "github.com/milbertk/usersAppAK/database/dataBaseUserOfCompany"
)

// BSearchCompaniesByUserID retrieves active companies where the user is owner or member.
func BSearchCompaniesByUserID(userID string) ([]databaseuserofcompany.CompanyByUserResult, error) {

	// 1) Validate user ID
	userID = strings.TrimSpace(userID)

	if userID == "" {
		return nil, fmt.Errorf("❌ user ID is required")
	}

	// 2) Load companies from database
	companies, err := databaseuserofcompany.SearchCompaniesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("❌ error loading companies by user ID: %v", err)
	}

	// 3) Return result
	return companies, nil
}
