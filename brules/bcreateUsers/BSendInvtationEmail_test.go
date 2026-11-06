package bcreateusers

import (
	"errors"
	"testing"
)

func TestLoadGeneralSMTPConfig(t *testing.T) {
	for _, k := range []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASSWORD", "SMTP_FROM", "SMTP_FROM_NAME"} {
		t.Setenv(k, "")
	}
	if _, err := LoadGeneralSMTPConfig(); !errors.Is(err, ErrEmailNotConfigured) {
		t.Fatalf("empty settings must return ErrEmailNotConfigured, got %v", err)
	}

	t.Setenv("SMTP_HOST", "mail.example.com")
	t.Setenv("SMTP_USER", "notify@example.com")
	t.Setenv("SMTP_PASSWORD", "x")
	cfg, err := LoadGeneralSMTPConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "587" || cfg.From != "notify@example.com" || cfg.FromName != "AlephK" {
		t.Errorf("defaults wrong: %+v", cfg)
	}

	t.Setenv("SMTP_FROM_NAME", "Kuantum K")
	t.Setenv("SMTP_FROM", "no-reply@example.com")
	cfg, _ = LoadGeneralSMTPConfig()
	if cfg.FromName != "Kuantum K" || cfg.From != "no-reply@example.com" {
		t.Errorf("overrides wrong: %+v", cfg)
	}
}

func TestInvitationLink(t *testing.T) {
	t.Setenv("APP_PUBLIC_URL", "")
	if got := InvitationLink(); got != "" {
		t.Errorf("empty base must give empty link, got %q", got)
	}
	t.Setenv("APP_PUBLIC_URL", " https://app.example.com/ ")
	if got := InvitationLink(); got != "https://app.example.com/manageinvitatinsbyuser" {
		t.Errorf("got %q", got)
	}
}
