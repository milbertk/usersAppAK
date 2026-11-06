package databaseSupportRequest

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert saves a new support request (status OPEN) and fills its id and date.
func (s *SupportRequest) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("DB connection error: %v", err)
	}

	s.Status = "OPEN"
	query := `
		INSERT INTO public.tbl_support_request (
			ticket_code, user_uid, user_email, company_id, company_name, platform,
			category, priority, subject, description, contact_method, contact_value,
			preferred_date, time_slots, image_urls, current_path, language, status
		)
		VALUES (
			$1, $2, NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6,
			$7, $8, $9, $10, $11, $12,
			NULLIF($13, '')::date, NULLIF($14, ''), NULLIF($15, ''), NULLIF($16, ''), NULLIF($17, ''), $18
		)
		RETURNING support_request_id, created_at
	`

	err = db.QueryRow(
		query,
		s.TicketCode, s.UserUID, s.UserEmail, s.CompanyID, s.CompanyName, s.Platform,
		s.Category, s.Priority, s.Subject, s.Description, s.ContactMethod, s.ContactValue,
		s.PreferredDate, s.TimeSlots, s.ImageURLs, s.CurrentPath, s.Language, s.Status,
	).Scan(&s.SupportRequestID, &s.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert support request: %v", err)
	}
	return nil
}
