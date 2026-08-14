package security

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct horse", "pepper")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	if !VerifyPassword(hash, "correct horse", "pepper") {
		t.Fatal("expected password to verify")
	}

	if VerifyPassword(hash, "wrong password", "pepper") {
		t.Fatal("expected wrong password to fail")
	}
}
