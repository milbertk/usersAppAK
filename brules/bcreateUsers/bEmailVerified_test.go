package bcreateusers

import "testing"

func TestBEffectiveEmailVerified(t *testing.T) {
	cases := []struct {
		claim, provider, email, want string
	}{
		{"true", "google.com", "a@gmail.com", "true"},
		{"false", "facebook.com", "a@hotmail.com", "true"},
		{"false", "FACEBOOK.COM", " a@hotmail.com ", "true"},
		{"false", "facebook.com", "", "false"},
		{"false", "facebook.com", "   ", "false"},
		{"false", "microsoft.com", "a@company.com", "false"},
		{"false", "password", "a@b.com", "false"},
		{"true", "password", "a@b.com", "true"},
		{"true", "apple.com", "x@privaterelay.appleid.com", "true"},
	}
	for _, c := range cases {
		if got := BEffectiveEmailVerified(c.claim, c.provider, c.email); got != c.want {
			t.Errorf("claim=%q provider=%q email=%q: got %q, want %q", c.claim, c.provider, c.email, got, c.want)
		}
	}
}
