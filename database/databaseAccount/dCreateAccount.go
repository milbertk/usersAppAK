package databaseaccount

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Insert inserts the DatabaseAccount into the account table
func (a *DatabaseAccountS) Insert() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		INSERT INTO account (
			 ownerid, createdate, statusaccount
		)
		VALUES ($1, $2, $3)
	`

	_, err = db.Exec(query,

		a.OwnerID,
		a.CreateDate,
		a.StatusAccount,
	)

	if err != nil {
		return fmt.Errorf("❌ Failed to insert account: %v", err)
	}

	fmt.Println("✅ Account inserted successfully:", a.ID)
	return nil
}
