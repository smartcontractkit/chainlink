package dashboard_lib

import (
	"encoding/base64"
	"testing"
)

func TestDecodeBasicAuth(t *testing.T) {
	// Valid Base64-encoded string
	validUser := "user"
	validPassword := "password"
	validAuthString := base64.StdEncoding.EncodeToString([]byte(validUser + ":" + validPassword))

	t.Run("Valid Base64 string", func(t *testing.T) {
		user, password, err := DecodeBasicAuth(validAuthString)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user != validUser {
			t.Errorf("expected user %s, got %s", validUser, user)
		}
		if password != validPassword {
			t.Errorf("expected password %s, got %s", validPassword, password)
		}
	})

	// Invalid Base64 string
	t.Run("Invalid Base64 string", func(t *testing.T) {
		invalidBase64String := "invalidbase64!!"
		_, _, err := DecodeBasicAuth(invalidBase64String)
		if err == nil {
			t.Fatal("expected error for invalid base64 string, got nil")
		}
	})

	// Plain text input (fallback)
	t.Run("Plain text input", func(t *testing.T) {
		plainAuthString := "user:password"
		_, _, err := DecodeBasicAuth(plainAuthString)
		if err == nil {
			t.Fatal("expected error for invalid base64 string, got nil")
		}
	})

	// Invalid format (no colon)
	t.Run("Missing colon in input", func(t *testing.T) {
		invalidFormat := base64.StdEncoding.EncodeToString([]byte("invalidformat"))
		_, _, err := DecodeBasicAuth(invalidFormat)
		if err == nil {
			t.Fatal("expected error for missing colon in the decoded string, got nil")
		}
	})

	// Edge case: empty string
	t.Run("Empty string", func(t *testing.T) {
		emptyString := ""
		_, _, err := DecodeBasicAuth(emptyString)
		if err == nil {
			t.Fatal("expected error for empty string, got nil")
		}
	})
}
