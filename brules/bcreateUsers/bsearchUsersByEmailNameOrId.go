package bcreateusers

import (
	"errors"
	"net/mail"
	"strings"

	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

// ErrInvalidUserSearch: the search is not one complete email address.
var ErrInvalidUserSearch = errors.New("type the complete email address of the user")

// BsearchUsersByEmailNameOrId finds the user to invite (2026-10-25, security
// S10): ONLY an exact, complete email address is accepted (no partial text,
// no names, no ids) and at most one user is returned, so the directory of
// users cannot be listed or harvested. (The name is kept so the callers do
// not change.)
func BsearchUsersByEmailNameOrId(text string) ([]databaseusers.FirebaseUserSearchResult, error) {
	email := strings.TrimSpace(text)
	if len(email) < 6 || len(email) > 254 || strings.ContainsAny(email, " ,;<>\r\n") {
		return nil, ErrInvalidUserSearch
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || !strings.Contains(email[strings.LastIndex(email, "@")+1:], ".") {
		return nil, ErrInvalidUserSearch
	}
	return databaseusers.SearchUserByExactEmail(email)
}
