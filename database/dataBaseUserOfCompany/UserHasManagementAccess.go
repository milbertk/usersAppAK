package databaseuserofcompany

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// UserHasManagementAccess validates if a user belongs to the company and has
// permissions to manage/search company users.
//
// Access is granted when ANY of the following is true:
//  1. The user has an active row in public.company_users with role OWNER or MANAGER.
//  2. The user is the registered owner of the company in public.companies
//     (companies.owner_firebase_uid = userID), even if no row exists in
//     public.company_users.
func UserHasManagementAccess(companyID string, userID string) (bool, string, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return false, "", fmt.Errorf("DB connection error: %v", err)
	}

	println("Valida user")
	println(companyID, userID)

	companyID = strings.TrimSpace(companyID)
	userID = strings.TrimSpace(userID)

	if companyID == "" {
		return false, "", fmt.Errorf("companyId is required")
	}
	if userID == "" {
		return false, "", fmt.Errorf("userId is required")
	}

	// Single query that gathers all possible role sources for this (company, user):
	//   - role from company_users (only active rows)
	//   - synthetic "OWNER" role if the user is the company's registered owner
	//
	// `priority` lets us deterministically pick the strongest match:
	//   1 = OWNER from companies.owner_firebase_uid
	//   2 = role from company_users (whatever it is)
	query := `
		SELECT role, priority FROM (
			SELECT 'OWNER'::text AS role, 1 AS priority
			FROM public.companies
			WHERE company_id = $1
			  AND owner_firebase_uid = $2

			UNION ALL

			SELECT role, 2 AS priority
			FROM public.company_users
			WHERE company_id = $1
			  AND user_id = $2
			  AND active = true
		) AS access
		ORDER BY priority ASC
		LIMIT 1
	`

	var role string
	var priority int
	err = db.QueryRow(query, companyID, userID).Scan(&role, &priority)
	if err != nil {
		if err == sql.ErrNoRows {
			// No match in either source — user has no relationship with this company.
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to validate company access: %v", err)
	}

	role = strings.ToUpper(strings.TrimSpace(role))

	switch role {
	case "OWNER", "MANAGER":
		return true, role, nil
	default:
		// User exists in company_users but with a non-management role
		// (e.g. USER, SUPPORT). Return the role for context, but deny access.
		return false, role, nil
	}
}
