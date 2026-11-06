package databaseRequestUserCompany

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// UpdateStatusAndMessage updates the status and message of a request.
func UpdateStatusAndMessage(requestID, status, message string) error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	requestID = strings.TrimSpace(requestID)
	status = strings.TrimSpace(strings.ToUpper(status))
	message = strings.TrimSpace(message)

	query := `
		UPDATE public.company_user_requests
		SET
			status = $2::text,
			message = $3::text,
			updated_at = NOW(),
			resolved_at = CASE
				WHEN $2::text IN ('APPROVED', 'REJECTED', 'CANCELED', 'EXPIRED') THEN NOW()
				ELSE resolved_at
			END
		WHERE request_id = $1
	`

	fmt.Println("=== UpdateStatusAndMessage ===")
	fmt.Println("RequestID:", requestID)
	fmt.Println("Status:", status)

	result, err := db.Exec(query, requestID, status, message)
	if err != nil {
		return fmt.Errorf("failed to update company user request: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no request found with request_id: %s", requestID)
	}

	return nil
}
