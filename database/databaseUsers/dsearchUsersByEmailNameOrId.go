package databaseusers

import (
	"fmt"
	"strings"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SearchUsersByEmailNameOrID returns up to 100 nearest matches by uid, display name, or email.
func SearchUsersByEmailNameOrID(searchText string) ([]FirebaseUserSearchResult, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	searchText = strings.TrimSpace(searchText)
	if searchText == "" {
		return []FirebaseUserSearchResult{}, nil
	}

	likeValue := "%" + searchText + "%"

	query := `
		SELECT
			uid,
			COALESCE(displayname, '') AS displayname,
			COALESCE(email, '') AS email
		FROM usfirebasedata
		WHERE
			uid ILIKE $1
			OR displayname ILIKE $1
			OR email ILIKE $1
		ORDER BY
			CASE
				WHEN uid ILIKE $2 THEN 1
				WHEN email ILIKE $2 THEN 2
				WHEN displayname ILIKE $2 THEN 3
				ELSE 4
			END,
			displayname ASC,
			email ASC
		LIMIT 100
	`

	// $1 = contains search
	// $2 = starts with search (higher priority)
	startsWithValue := searchText + "%"

	rows, err := db.Query(query, likeValue, startsWithValue)
	if err != nil {
		return nil, fmt.Errorf("❌ query error: %v", err)
	}
	defer rows.Close()

	results := make([]FirebaseUserSearchResult, 0, 100)

	for rows.Next() {
		var u FirebaseUserSearchResult
		if err := rows.Scan(
			&u.UID,
			&u.DisplayName,
			&u.Email,
		); err != nil {
			return nil, fmt.Errorf("❌ scan error: %v", err)
		}
		results = append(results, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("❌ rows error: %v", err)
	}

	return results, nil
}
