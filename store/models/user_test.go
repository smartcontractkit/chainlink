package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	zeroTime := models.Time{}

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
				assert.Equal(t, zeroTime, user.CreatedAt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)
				newHash, _ := utils.HashPassword(test.pwd)
				assert.NotEmpty(t, newHash, user.HashedPassword)
				assert.NotEqual(t, zeroTime, user.CreatedAt)
			}
		})
	}
}
