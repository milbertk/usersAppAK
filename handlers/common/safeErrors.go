package common

import (
	"log"
	"strings"
)

// Safe error messages for the browser (2026-10-25, security S15).
//
// Business messages ("company name exists", "passwords do not match") are
// useful for the user and are returned as they are. Technical details
// (database, SQL, Firebase, network) must never reach the browser: they are
// written to the server log and the user receives a generic message.

// internalErrorMarkers identify technical errors (lower case).
var internalErrorMarkers = []string{
	"pq:", "sql:", "sql ", "database", "db connection", "connection", "query", "scan",
	"rows", "relation ", "column ", "syntax", "violates", "constraint",
	"dial ", "tcp", "timeout", "x509", "tls", "eof",
	"firebase", "googleapi", "credential", "json:", "unexpected", "panic",
	"insert failed", "update failed", "failed to",
}

// ClientMessage returns err's text when it is a business message, or fallback
// (and logs the detail) when it looks technical.
func ClientMessage(err error, fallback string) string {
	if err == nil {
		return fallback
	}
	message := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(err.Error()), "❌"))
	low := strings.ToLower(message)
	for _, marker := range internalErrorMarkers {
		if strings.Contains(low, marker) {
			log.Printf("[safeErrors] detail hidden from the client: %v", err)
			return fallback
		}
	}
	if message == "" {
		return fallback
	}
	return message
}

// InternalMessage always hides the detail: logs it and returns fallback.
func InternalMessage(err error, fallback string) string {
	if err != nil {
		log.Printf("[safeErrors] detail hidden from the client: %v", err)
	}
	return fallback
}

// InvalidJSONMessage logs the JSON decoding detail and returns a generic text.
func InvalidJSONMessage(err error) string {
	return InternalMessage(err, "Invalid JSON body")
}
