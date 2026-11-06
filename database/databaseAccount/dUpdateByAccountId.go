package databaseaccount

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Update updates an existing DatabaseAccount record in the account table
func (a *DatabaseAccountS) Update() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		UPDATE account
		SET
			ownerid = $1,
			createdate = $2,
			statusaccount = $3
		WHERE id = $4
	`

	result, err := db.Exec(query,
		a.OwnerID,
		a.CreateDate,
		a.StatusAccount,
		a.ID,
	)
	if err != nil {
		return fmt.Errorf("❌ failed to update account: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("⚠️ no account found with ID: %s", a.ID)
	}

	fmt.Println("✅ Account updated successfully:", a.ID)
	return nil
}
