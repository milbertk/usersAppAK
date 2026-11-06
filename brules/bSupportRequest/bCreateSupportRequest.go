package bSupportRequest

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/milbertk/usersAppAK/database/databaseSupportRequest"
)

// Support requests (2026-10-24): the "Support" form of the apps. The request
// is saved in tbl_support_request (usersAK) and emailed in the background to
// SUPPORT_EMAIL_TO with the SMTP of these environment variables:
//   SMTP_HOST, SMTP_PORT (587), SMTP_USER, SMTP_PASSWORD, SMTP_FROM (= SMTP_USER)
// When they are not set the request is still saved (email_sent = false).

var (
	// ErrInvalidSupportRequest wraps validation problems (HTTP 400).
	ErrInvalidSupportRequest = errors.New("invalid support request")
	// ErrTooManySupportRequests: too many requests from the user today (HTTP 429).
	ErrTooManySupportRequests = errors.New("too many support requests today, we will contact you soon")
)

const maxRequestsPerDay = 10

var (
	validCategories = map[string]bool{"technical": true, "account": true, "billing": true, "howTo": true, "other": true}
	validPriorities = map[string]bool{"low": true, "normal": true, "high": true}
	validMethods    = map[string]bool{"email": true, "phone": true, "whatsapp": true, "videoCall": true}
	validSlots      = map[string]bool{"morning": true, "afternoon": true, "evening": true}
	datePattern     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// CreateSupportRequestInput is what the form sends.
type CreateSupportRequestInput struct {
	CompanyID     string   `json:"companyId"`
	CompanyName   string   `json:"companyName"`
	Platform      string   `json:"platform"`
	Category      string   `json:"category"`
	Priority      string   `json:"priority"`
	Subject       string   `json:"subject"`
	Description   string   `json:"description"`
	ContactMethod string   `json:"contactMethod"`
	ContactValue  string   `json:"contactValue"`
	PreferredDate string   `json:"preferredDate"`
	TimeSlots     []string `json:"timeSlots"`
	ImageURLs     []string `json:"imageUrls"`
	CurrentPath   string   `json:"currentPath"`
	Language      string   `json:"language"`
}

func invalid(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrInvalidSupportRequest, fmt.Sprintf(format, args...))
}

func tooLong(value string, max int) bool { return utf8.RuneCountInString(value) > max }

// normalizeSupportRequest trims and validates the form (pure logic).
func normalizeSupportRequest(in CreateSupportRequestInput) (databaseSupportRequest.SupportRequest, error) {
	out := databaseSupportRequest.SupportRequest{
		CompanyID:     strings.TrimSpace(in.CompanyID),
		CompanyName:   strings.TrimSpace(in.CompanyName),
		Platform:      strings.ToLower(strings.TrimSpace(in.Platform)),
		Category:      strings.TrimSpace(in.Category),
		Priority:      strings.TrimSpace(in.Priority),
		Subject:       strings.TrimSpace(in.Subject),
		Description:   strings.TrimSpace(in.Description),
		ContactMethod: strings.TrimSpace(in.ContactMethod),
		ContactValue:  strings.TrimSpace(in.ContactValue),
		PreferredDate: strings.TrimSpace(in.PreferredDate),
		CurrentPath:   strings.TrimSpace(in.CurrentPath),
		Language:      strings.ToLower(strings.TrimSpace(in.Language)),
	}
	if out.Platform == "" {
		out.Platform = "shao"
	}
	if out.Priority == "" {
		out.Priority = "normal"
	}
	if strings.HasPrefix(out.Language, "es") {
		out.Language = "es"
	} else {
		out.Language = "en"
	}

	var slots []string
	for _, s := range in.TimeSlots {
		s = strings.TrimSpace(s)
		if validSlots[s] {
			slots = append(slots, s)
		}
	}
	out.TimeSlots = strings.Join(slots, ",")

	var urls []string
	for _, u := range in.ImageURLs {
		u = strings.TrimSpace(u)
		if strings.HasPrefix(u, "https://") && len(u) <= 1000 {
			urls = append(urls, u)
		}
	}
	if len(urls) > 3 {
		urls = urls[:3]
	}
	out.ImageURLs = strings.Join(urls, "\n")

	switch {
	case !validCategories[out.Category]:
		return out, invalid("category is not valid")
	case !validPriorities[out.Priority]:
		return out, invalid("priority is not valid")
	case !validMethods[out.ContactMethod]:
		return out, invalid("contactMethod is not valid")
	case out.ContactValue == "" || tooLong(out.ContactValue, 200):
		return out, invalid("contactValue is required (max 200 characters)")
	case out.Subject == "" || tooLong(out.Subject, 200):
		return out, invalid("subject is required (max 200 characters)")
	case out.Description == "" || tooLong(out.Description, 5000):
		return out, invalid("description is required (max 5000 characters)")
	case out.PreferredDate != "" && !datePattern.MatchString(out.PreferredDate):
		return out, invalid("preferredDate must be YYYY-MM-DD")
	case tooLong(out.CompanyID, 100) || tooLong(out.CompanyName, 200) || tooLong(out.CurrentPath, 300) || tooLong(out.Platform, 40):
		return out, invalid("a field is too long")
	}
	return out, nil
}

func newTicketCode(now time.Time) string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("SUP-%s-%02X%02X", now.Format("060102"), b[0], b[1])
}

// BCreateSupportRequest validates, saves and emails a support request.
func BCreateSupportRequest(userUID, userEmail string, in CreateSupportRequestInput) (*databaseSupportRequest.SupportRequest, error) {
	userUID = strings.TrimSpace(userUID)
	if userUID == "" {
		return nil, invalid("user is required")
	}

	request, err := normalizeSupportRequest(in)
	if err != nil {
		return nil, err
	}

	count, err := databaseSupportRequest.CountUserRequestsSince(userUID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}
	if count >= maxRequestsPerDay {
		return nil, ErrTooManySupportRequests
	}

	request.UserUID = userUID
	request.UserEmail = strings.TrimSpace(userEmail)

	// The ticket code is unique; retry once on a (very rare) collision.
	for attempt := 0; attempt < 2; attempt++ {
		request.TicketCode = newTicketCode(time.Now())
		if err = request.Insert(); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	go sendSupportEmail(request)
	return &request, nil
}

func env(name, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return fallback
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func sendSupportEmail(request databaseSupportRequest.SupportRequest) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[bSupportRequest] email panic: %v", r)
		}
	}()

	to := env("SUPPORT_EMAIL_TO", "")
	host := env("SMTP_HOST", "")
	user := env("SMTP_USER", "")
	pass := env("SMTP_PASSWORD", "")
	if to == "" || host == "" || user == "" {
		log.Printf("[bSupportRequest] %s saved; email not sent (set SUPPORT_EMAIL_TO, SMTP_HOST, SMTP_USER, SMTP_PASSWORD)", request.TicketCode)
		return
	}
	port := env("SMTP_PORT", "587")
	from := env("SMTP_FROM", user)

	body := strings.Join([]string{
		"New support request " + request.TicketCode,
		"",
		"Priority: " + request.Priority + " · Category: " + request.Category,
		"Platform: " + request.Platform + " · Screen: " + emptyDash(request.CurrentPath),
		"Company: " + emptyDash(request.CompanyName) + " (" + emptyDash(request.CompanyID) + ")",
		"User: " + emptyDash(request.UserEmail) + " (" + request.UserUID + ")",
		"Language: " + request.Language,
		"",
		"Contact by: " + request.ContactMethod + " — " + request.ContactValue,
		"Preferred date: " + emptyDash(request.PreferredDate) + " · Time: " + emptyDash(request.TimeSlots),
		"",
		"Subject: " + request.Subject,
		"",
		request.Description,
		"",
		"Images:",
		emptyDash(request.ImageURLs),
		"",
		"Close it with the SQL in usersAppAK/sql/supportRequests.sql.",
	}, "\n")

	subject := fmt.Sprintf("[Support %s] %s — %s", request.Priority, request.TicketCode, request.Subject)
	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + strings.ReplaceAll(subject, "\n", " ") + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" + body,
	)

	auth := smtp.PlainAuth("", user, pass, host)
	if err := smtp.SendMail(host+":"+port, auth, from, []string{to}, message); err != nil {
		log.Printf("[bSupportRequest] email failed: %s %v", request.TicketCode, err)
		return
	}
	if err := databaseSupportRequest.MarkEmailSent(request.SupportRequestID); err != nil {
		log.Printf("[bSupportRequest] email_sent not saved: %s %v", request.TicketCode, err)
	}
}
