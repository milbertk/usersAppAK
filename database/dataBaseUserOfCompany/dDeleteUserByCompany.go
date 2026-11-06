package databaseuserofcompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Delete soft deletes a company user by setting active = false.
func Delete(companyUserID string) error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	companyUserID = strings.TrimSpace(companyUserID)
	if companyUserID == "" {
		return fmt.Errorf("companyUserId is required")
	}

	query := `
		UPDATE public.company_users
		SET
			active = false,
			removed_at = NOW(),
			updated_at = NOW()
		WHERE company_user_id = $1
	`

	result, err := db.Exec(query, companyUserID)
	if err != nil {
		return fmt.Errorf("failed to delete company user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify delete result: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no company user found with the given id")
	}

	return nil
}
