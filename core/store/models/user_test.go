package models_test

import (
	"testing"

	"chainlink/core/store/models"
	"chainlink/core/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		email, pwd string
		wantError  bool
	}{
		{"good@email.com", "goodpassword", false},
		{"notld@email", "goodpassword", false},
		{"good@email.com", "badpd", true},
		{"bademail", "goodpassword", true},
		{"bad@", "goodpassword", true},
		{"@email", "goodpassword", true},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			user, err := models.NewUser(test.email, test.pwd)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)
				assert.NotEmpty(t, user.HashedPassword)
				newHash, _ := utils.HashPassword(test.pwd)
				assert.NotEqual(t, newHash, user.HashedPassword, "Salt should prevent equality")
			}
		})
	}
}
