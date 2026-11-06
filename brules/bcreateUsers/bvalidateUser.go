package bcreateusers

import (
	"fmt"

	"github.com/milbertk/security"
	databaseusers "github.com/milbertk/usersAppAK/database/databaseUsers"
)

// BCreateUsers inserts a full FirebaseUser object
func BvalidateUser(user databaseusers.FirebaseUser) error {

	except, exist, erru := security.DValidateUser(user.UID)

	if except != true {

		return fmt.Errorf("❌ the user does not exist, can not be created again, revoer your account")
	}

	if erru != nil {

		return fmt.Errorf("❌ Error validating the user: %v", erru)
	}

	if exist != "true" {

		return fmt.Errorf("❌ the user exist, can not be created again, recover your account")
	}

	err := user.Update()
	if err != nil {
		return fmt.Errorf("❌ Insert failed: %v", err)
	}

	return nil
}
