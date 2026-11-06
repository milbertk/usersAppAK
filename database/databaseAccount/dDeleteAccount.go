package databaseaccount

import (
	"fmt"

	dbmanager "github.com/milbertk/usersAppAK/database/dbManager"
)

// Delete deletes an existing DatabaseAccount record from the account table
func (a *DatabaseAccountS) Delete() error {
	db, err := dbmanager.GetPostgresDB()
	if err != nil {
		return fmt.Errorf("❌ DB connection error: %v", err)
	}

	query := `
		DELETE FROM account
		WHERE id = $1
	`

	result, err := db.Exec(query, a.ID)
	if err != nil {
		return fmt.Errorf("❌ failed to delete account: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("⚠️ no account found with ID: %s", a.ID)
	}

	fmt.Println("✅ Account deleted successfully:", a.ID)
	return nil
}
