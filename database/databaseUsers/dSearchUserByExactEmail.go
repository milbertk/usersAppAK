package databaseusers

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SearchUserByExactEmail returns the user with exactly this email (case
// insensitive), or an empty list. No partial matches and never a list of
// users (2026-10-25, security S10).
func SearchUserByExactEmail(email string) ([]FirebaseUserSearchResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return []FirebaseUserSearchResult{}, nil
	}

	rows, err := db.Query(`
		SELECT
			uid,
			COALESCE(displayname, '') AS displayname,
			COALESCE(email, '') AS email
		FROM usfirebasedata
		WHERE lower(email) = lower($1)
		LIMIT 1
	`, email)
	if err != nil {
		return nil, fmt.Errorf("❌ query error: %v", err)
	}
	defer rows.Close()

	results := make([]FirebaseUserSearchResult, 0, 1)
	for rows.Next() {
		var u FirebaseUserSearchResult
		if err := rows.Scan(&u.UID, &u.DisplayName, &u.Email); err != nil {
			return nil, fmt.Errorf("❌ scan error: %v", err)
		}
		results = append(results, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("❌ rows error: %v", err)
	}
	return results, nil
}
