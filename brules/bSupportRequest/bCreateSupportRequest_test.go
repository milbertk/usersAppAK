package bSupportRequest

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNormalizeSupportRequest(t *testing.T) {
	ok := CreateSupportRequestInput{
		Category: "technical", Subject: " Cannot save ", Description: "It fails",
		ContactMethod: "whatsapp", ContactValue: "+506 8888 7777", TimeSlots: []string{"morning", "night"},
		ImageURLs: []string{"https://x/a.png", "http://bad", "https://x/b.png", "https://x/c.png", "https://x/d.png"},
		Language:  "es-CR", PreferredDate: "2026-10-30",
	}
	got, err := normalizeSupportRequest(ok)
	if err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}
	if got.Subject != "Cannot save" || got.Priority != "normal" || got.Platform != "shao" || got.Language != "es" {
		t.Fatalf("defaults: %+v", got)
	}
	if got.TimeSlots != "morning" || strings.Count(got.ImageURLs, "\n") != 2 {
		t.Fatalf("slots/images: %q %q", got.TimeSlots, got.ImageURLs)
	}
	bad := []CreateSupportRequestInput{
		{Category: "x", Subject: "a", Description: "b", ContactMethod: "email", ContactValue: "c"},
		{Category: "other", Subject: "", Description: "b", ContactMethod: "email", ContactValue: "c"},
		{Category: "other", Subject: "a", Description: "b", ContactMethod: "fax", ContactValue: "c"},
		{Category: "other", Subject: "a", Description: "b", ContactMethod: "email", ContactValue: "c", PreferredDate: "30/10/2026"},
	}
	for _, in := range bad {
		if _, err := normalizeSupportRequest(in); !errors.Is(err, ErrInvalidSupportRequest) {
			t.Errorf("expected invalid for %+v", in)
		}
	}
	if code := newTicketCode(time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)); !strings.HasPrefix(code, "SUP-261024-") || len(code) != 15 {
		t.Fatalf("ticket code %q", code)
	}
}
