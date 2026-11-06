package databaseSupportRequest

import "time"

// SupportRequest is one row of public.tbl_support_request (usersAK): a
// support request sent from any platform by a signed-in user (2026-10-24).
type SupportRequest struct {
	SupportRequestID string    `json:"supportRequestId"`
	TicketCode       string    `json:"ticketCode"`
	UserUID          string    `json:"-"`
	UserEmail        string    `json:"userEmail"`
	CompanyID        string    `json:"companyId,omitempty"`
	CompanyName      string    `json:"companyName,omitempty"`
	Platform         string    `json:"platform"`
	Category         string    `json:"category"`
	Priority         string    `json:"priority"`
	Subject          string    `json:"subject"`
	Description      string    `json:"description"`
	ContactMethod    string    `json:"contactMethod"`
	ContactValue     string    `json:"contactValue"`
	PreferredDate    string    `json:"preferredDate,omitempty"` // YYYY-MM-DD or ""
	TimeSlots        string    `json:"timeSlots,omitempty"`     // comma separated
	ImageURLs        string    `json:"-"`                       // one per line
	CurrentPath      string    `json:"currentPath,omitempty"`
	Language         string    `json:"language,omitempty"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
}
