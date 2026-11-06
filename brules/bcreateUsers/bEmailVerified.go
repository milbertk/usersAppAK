package bcreateusers

import "strings"

// BEffectiveEmailVerified decides the email_verified value saved for a user (2026-10-27).
//
// Firebase marks Facebook emails as "not verified", but Facebook only shares an
// email that the user already confirmed with Facebook. So:
//   - Facebook ("facebook.com") WITH an email  -> "true"
//   - everything else (Google, Microsoft, Apple, email + password, Facebook
//     WITHOUT an email)                         -> the value Firebase sent
//
// Microsoft keeps Firebase's value on purpose: work / school emails are set by
// each organization and are not confirmed by Microsoft.
func BEffectiveEmailVerified(firebaseClaim, signInProvider, email string) string {
	if strings.EqualFold(strings.TrimSpace(firebaseClaim), "true") {
		return "true"
	}
	if strings.EqualFold(strings.TrimSpace(signInProvider), "facebook.com") && strings.TrimSpace(email) != "" {
		return "true"
	}
	return firebaseClaim
}
