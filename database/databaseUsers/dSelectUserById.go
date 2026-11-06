package databaseusers

import (
	"database/sql"
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// SelectByUID retrieves a single FirebaseUser by UID (primary key)
func SelectByUID(uid string) (*FirebaseUser, error) {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return nil, fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		SELECT uid, lastsignintime, localstatus, creationtime,
		       displayname, phonenumber, photourl, active,
		       providerid, email, emailverified, disabled
		FROM usfirebasedata
		WHERE uid = $1
	`

	var u FirebaseUser
	err = db.QueryRow(query, uid).Scan(
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
			return nil, fmt.Errorf("⚠️ no user found with UID: %s", uid)
		}
		return nil, fmt.Errorf("❌ query error: %v", err)
	}

	fmt.Println("✅ User loaded successfully:", u.UID)
	return &u, nil
}
