package databaseSupportRequest

import (
	"fmt"
	"strings"
	"time"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// MarkEmailSent records that the support email reached the mail server.
func MarkEmailSent(supportRequestID string) error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}
	_, err = db.Exec(`
		UPDATE public.tbl_support_request
		SET email_sent = TRUE, updated_at = NOW()
		WHERE support_request_id::text = $1
	`, strings.TrimSpace(supportRequestID))
	if err != nil {
		return fmt.Errorf("failed to update support request: %v", err)
	}
	return nil
}

// CountUserRequestsSince counts the requests of a user since a time
// (protection against repeated sends).
func CountUserRequestsSince(userUID string, since time.Time) (int, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return 0, fmt.Errorf("DB connection error: %v", err)
	}
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM public.tbl_support_request
		WHERE user_uid = $1 AND created_at >= $2
	`, strings.TrimSpace(userUID), since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count support requests: %v", err)
	}
	return count, nil
}
