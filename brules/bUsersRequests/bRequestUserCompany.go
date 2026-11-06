package busersrequests

import (
	"fmt"
	"strings"

	databasecompany "github.com/milbertk/usersAppAK/database/databaseCompany"
	databaserequestusercompany "github.com/milbertk/usersAppAK/database/databaseRequestUserCompany"
)

// BSearchRequest validates that the searching user is the owner of the company
// and then returns the matching requests.
func BSearchRequest(companyID, searchText, requesterUID string) ([]databaserequestusercompany.CompanyUserRequestSearchResult, error) {
	companyID = strings.TrimSpace(companyID)
	searchText = strings.TrimSpace(searchText)
	requesterUID = strings.TrimSpace(requesterUID)

	fmt.Println("=== BSearchRequest ===")
	fmt.Println("CompanyID:", companyID)
	fmt.Println("SearchText:", searchText)
	fmt.Println("RequesterUID:", requesterUID)

	if companyID == "" {
		return nil, fmt.Errorf("company id is required")
	}

	if requesterUID == "" {
		return nil, fmt.Errorf("requester uid is required")
	}

	// 1. Validate company exists
	company, err := databasecompany.SelectByCompanyID(companyID)
	if err != nil {
		return nil, fmt.Errorf("error validating company: %v", err)
	}

	// 2. Validate requester is the owner
	if strings.TrimSpace(company.OwnerFirebaseUID) != requesterUID {
		return nil, fmt.Errorf("you are not allowed to search requests for this company")
	}

	// 3. Execute search
	results, err := databaserequestusercompany.SearchByCompanyAndUserText(companyID, searchText)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	return results, nil
}
