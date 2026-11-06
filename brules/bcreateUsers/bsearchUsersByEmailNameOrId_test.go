package bcreateusers

import (
	"errors"
	"testing"
)

func TestUserSearchRejectsPartialText(t *testing.T) {
	for _, bad := range []string{"", "a", "ana", "@gmail.com", "ana@", "ana@gmail", "ana perez@gmail.com", "uid-123456", "a@b.c,d@e.f"} {
		if _, err := BsearchUsersByEmailNameOrId(bad); !errors.Is(err, ErrInvalidUserSearch) {
			t.Errorf("%q should be rejected, got %v", bad, err)
		}
	}
}
