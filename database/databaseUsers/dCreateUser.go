package databaseusers

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert inserts the FirebaseUser into the usfirebasedata table
func (u *FirebaseUser) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		INSERT INTO usfirebasedata (
			uid, lastsignintime, localstatus, creationtime,
			displayname, phonenumber, photourl, active,
			providerid, email, emailverified, disabled
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = db.Exec(query,
		u.UID, u.LastSignIn, u.LocalStatus, u.CreationTime,
		u.DisplayName, u.PhoneNumber, u.PhotoURL, u.Active,
		u.ProviderID, u.Email, u.EmailVerified, u.Disabled,
	)

	if err != nil {
		return fmt.Errorf("❌ Failed to insert user: %v", err)
	}

	fmt.Println("✅ User inserted successfully:", u.UID)
	return nil
}
