package sessions_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		email, pwd string
		role       sessions.UserRole
		wantError  bool
	}{
		{"good@email.com", cltest.Password, sessions.UserRoleAdmin, false},
		{"notld@email", cltest.Password, sessions.UserRoleEdit, false},
		{"view@email", cltest.Password, sessions.UserRoleView, false},
		{"good@email.com", "badpd", sessions.UserRoleAdmin, true},
		{"bademail", cltest.Password, sessions.UserRoleAdmin, true},
		{"bad@", cltest.Password, sessions.UserRoleAdmin, true},
		{"@email", cltest.Password, sessions.UserRoleAdmin, true},
		{"good@email.com", cltest.Password, sessions.UserRoleRun, false},
		{"good@email-pass-too-long.com", cltest.Password + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", sessions.UserRoleAdmin, true},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			user, err := sessions.NewUser(test.email, test.pwd, test.role)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.email, user.Email)
				assert.Equal(t, test.role, user.Role)
				assert.NotEmpty(t, user.HashedPassword)
				newHash, _ := utils.HashPassword(test.pwd)
				assert.NotEqual(t, newHash, user.HashedPassword, "Salt should prevent equality")
			}
		})
	}
}

func TestUserGenerateAuthToken(t *testing.T) {
	var user sessions.User
	token, err := user.GenerateAuthToken()
	require.NoError(t, err)
	assert.Equal(t, null.StringFrom(token.AccessKey), user.TokenKey)
	assert.NotEqual(t, null.StringFrom(token.Secret), user.TokenHashedSecret)
}

func TestAuthenticateUserByToken(t *testing.T) {
	var user sessions.User

	token, err := user.GenerateAuthToken()
	assert.NoError(t, err, "failed when generate auth token")
	ok, err := sessions.AuthenticateUserByToken(token, &user)
	require.NoError(t, err)
	assert.True(t, ok, "authentication must be successful")

	_, err = user.GenerateAuthToken()
	assert.NoError(t, err, "failed to generate auth token")
	ok, err = sessions.AuthenticateUserByToken(token, &user)
	require.NoError(t, err)
	assert.False(t, ok, "authentication must fail with past token")
}
