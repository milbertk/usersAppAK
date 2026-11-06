package databaseusers

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Update updates an existing FirebaseUser record in the usfirebasedata table
func (u *FirebaseUser) Update() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		UPDATE usfirebasedata
		SET 
			lastsignintime = $1,
			creationtime = $2,
			displayname = $3,
			phonenumber = $4,
			photourl = $5,
			active = $6,
			providerid = $7,
			email = $8,
			emailverified = $9,
			disabled = $10
		WHERE uid = $11
	`

	result, err := db.Exec(query,
		u.LastSignIn,
		u.CreationTime,
		u.DisplayName,
		u.PhoneNumber,
		u.PhotoURL,
		u.Active,
		u.ProviderID,
		u.Email,
		u.EmailVerified,
		u.Disabled,
		u.UID,
	)
	if err != nil {
		return fmt.Errorf("❌ failed to update user: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("⚠️ no user found with UID: %s", u.UID)
	}

	fmt.Println("✅ User updated successfully:", u.UID)
	return nil
}
