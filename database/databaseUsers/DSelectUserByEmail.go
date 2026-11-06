package databaseusers

import (
	"database/sql"
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// DSelectUserByEmail retrieves a single FirebaseUser by email
func DSelectUserByEmail(email string) (*FirebaseUser, error) {

	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		SELECT uid,
		       lastsignintime,
		       localstatus,
		       creationtime,
		       displayname,
		       phonenumber,
		       photourl,
		       active,
		       providerid,
		       email,
		       emailverified,
		       disabled
		FROM usfirebasedata
		WHERE LOWER(email) = LOWER($1)
	`

	var u FirebaseUser

	err = db.QueryRow(query, email).Scan(
		&u.UID,
		&u.LastSignIn,
		&u.LocalStatus,
		&u.CreationTime,
		&u.DisplayName,
		&u.PhoneNumber,
		&u.PhotoURL,
		&u.Active,
		&u.ProviderID,
		&u.Email,
		&u.EmailVerified,
		&u.Disabled,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("⚠️ no user found with email: %s", email)
		}

		return nil, fmt.Errorf("❌ query error: %v", err)
	}

	fmt.Println("✅ User loaded successfully:", u.UID) // id only, no email (security S11)

	return &u, nil
}
