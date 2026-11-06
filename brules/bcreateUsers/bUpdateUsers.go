package bcreateusers

import (
	"fmt"
	"strings"

	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

// BUpdateUsers updates an existing FirebaseUser
func BUpdateUsers(user databaseusers.FirebaseUser) error {

	// 1) Validate UID
	if strings.TrimSpace(user.UID) == "" {
		return fmt.Errorf("❌ UID is required to update the user")
	}

	// 2) Validate that the user exists
	existingUser, err := databaseusers.SelectByUID(user.UID)
	if err != nil {
		return fmt.Errorf("❌ user does not exist or cannot be loaded: %v", err)
	}

	// 3) Optional business rules (examples)
	// You can restrict changes here if needed
	// Example: do not allow changing provider
	user.ProviderID = existingUser.ProviderID

	// 4) Update user
	if err := user.Update(); err != nil {
		return fmt.Errorf("❌ Update failed: %v", err)
	}

	return nil
}
