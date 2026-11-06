package bcreateusers

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/smtp"
	"os"
	"strings"

	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	FromName string
}

// InvitationPagePath is the p3 page where a user accepts or declines invitations.
const InvitationPagePath = "/manageinvitatinsbyuser"

// InvitationLink builds the link for the invitation email (2026-10-29):
// APP_PUBLIC_URL (base address of p3, e.g. https://app.yourdomain.com) +
// InvitationPagePath. Returns "" when APP_PUBLIC_URL is not set.
// Change .env and restart usersAppAK, no rebuild needed.
func InvitationLink() string {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("APP_PUBLIC_URL")), "/")
	if base == "" {
		return ""
	}
	return base + InvitationPagePath
}

// ErrEmailNotConfigured is returned when the general email settings are missing.
var ErrEmailNotConfigured = errors.New("email is not configured (set SMTP_HOST, SMTP_USER and SMTP_PASSWORD in .env)")

// LoadGeneralSMTPConfig reads the general (platform) email settings from the
// environment / .env (2026-10-28). Nothing is hard-coded: change .env and
// restart usersAppAK, no rebuild needed.
//
//	SMTP_HOST, SMTP_PORT (default 587), SMTP_USER, SMTP_PASSWORD,
//	SMTP_FROM (default SMTP_USER), SMTP_FROM_NAME (default "AlephK")
func LoadGeneralSMTPConfig() (SMTPConfig, error) {
	get := func(name, fallback string) string {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
		return fallback
	}

	cfg := SMTPConfig{
		Host:     get("SMTP_HOST", ""),
		Port:     get("SMTP_PORT", "587"),
		Username: get("SMTP_USER", ""),
		Password: os.Getenv("SMTP_PASSWORD"),
		FromName: get("SMTP_FROM_NAME", "AlephK"),
	}
	cfg.From = get("SMTP_FROM", cfg.Username)

	if cfg.Host == "" || cfg.Username == "" || cfg.Password == "" {
		return cfg, ErrEmailNotConfigured
	}
	return cfg, nil
}

func BSendInvitationEmail(
	currentUser databaseusers.FirebaseUser,
	targetEmail string,
	messageText string,
) error {

	if strings.TrimSpace(targetEmail) == "" {
		return fmt.Errorf("target email is required")
	}

	if strings.TrimSpace(messageText) == "" {
		messageText = "You have been invited to join the platform."
	}

	// General email settings from .env (2026-10-28) — no hard-coded server or password.
	smtpConfig, err := LoadGeneralSMTPConfig()
	if err != nil {
		return err
	}

	inviterName := currentUser.DisplayName
	if strings.TrimSpace(inviterName) == "" {
		inviterName = currentUser.Email
	}

	subject := smtpConfig.FromName + " invitation"

	linkText := ""
	if link := InvitationLink(); link != "" {
		linkText = fmt.Sprintf(`
To accept or decline the invitation, open:
%s

If you don't have an account yet, create one with this same email address, then open the link again.
`, link)
	} else {
		log.Println("[invitation] APP_PUBLIC_URL is not set: invitation email sent without a link")
	}

	body := fmt.Sprintf(
		`Hello,

%s invited you to join the platform.

Message:
%s

Please check your invitation and continue the registration process.
%s
Regards,
%s`,
		inviterName,
		messageText,
		linkText,
		smtpConfig.FromName,
	)

	return sendSMTPEmail(
		smtpConfig,
		targetEmail,
		subject,
		body,
	)
}

func sendSMTPEmail(
	config SMTPConfig,
	to string,
	subject string,
	body string,
) error {

	auth := smtp.PlainAuth(
		"",
		config.Username,
		config.Password,
		config.Host,
	)

	address := config.Host + ":" + config.Port

	fromHeader := config.From
	if strings.TrimSpace(config.FromName) != "" {
		fromHeader = mime.QEncoding.Encode("UTF-8", config.FromName) + " <" + config.From + ">"
	}

	message := []byte(
		"From: " + fromHeader + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + mime.QEncoding.Encode("UTF-8", subject) + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body,
	)

	err := smtp.SendMail(
		address,
		auth,
		config.From,
		[]string{to},
		message,
	)

	if err != nil {
		return fmt.Errorf("error sending invitation email: %v", err)
	}

	return nil
}
