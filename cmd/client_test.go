package cmd_test

import (
	"path"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerminalCookieAuthenticator_AuthenticateWithoutSession(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	tests := []struct {
		name, email, pwd string
	}{
		{"bad email", "notreal", cltest.Password},
		{"bad pwd", cltest.UserEmail, "mostcommonwrongpwdever"},
		{"bad both", "notreal", "mostcommonwrongpwdever"},
		{"correct", cltest.UserEmail, cltest.Password},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entering := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{EnteredStrings: entering}

			tca := cmd.NewTerminalCookieAuthenticator(app.Config, prompter)
			cookie, err := tca.Authenticate()

			assert.Error(t, err)
			assert.Nil(t, cookie)
			assert.False(t, utils.FileExists(path.Join(app.Config.RootDir, "cookie")))
		})
	}
}

func TestTerminalCookieAuthenticator_AuthenticateWithSession(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	app.MustSeedUserSession()

	tests := []struct {
		name, email, pwd string
		wantError        bool
	}{
		{"bad email", "notreal", cltest.Password, true},
		{"bad pwd", cltest.UserEmail, "mostcommonwrongpwdever", true},
		{"bad both", "notreal", "mostcommonwrongpwdever", true},
		{"success", cltest.UserEmail, cltest.Password, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entering := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{EnteredStrings: entering}

			tca := cmd.NewTerminalCookieAuthenticator(app.Config, prompter)
			cookie, err := tca.Authenticate()

			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, cookie)
				assert.False(t, utils.FileExists(path.Join(app.Config.RootDir, "cookie")))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cookie)
				assert.True(t, utils.FileExists(path.Join(app.Config.RootDir, "cookie")))
			}
		})
	}
}

func TestTerminalCookieAuthenticator_Cookie(t *testing.T) {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	config := tc.Config
	tests := []struct {
		name      string
		rootDir   string
		wantError bool
	}{
		{"missing", config.RootDir, true},
		{"correct fixture", "../internal/fixtures", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.RootDir = test.rootDir
			tca := cmd.NewTerminalCookieAuthenticator(config, nil)
			cookie, err := tca.Cookie()
			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, cookie)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cookie)
			}
		})
	}
}

func TestTerminalAPIInitializer_InitializeWithoutAPIUser(t *testing.T) {
	tests := []struct {
		name           string
		enteredStrings []string
	}{
		{"correct", []string{"email", "password"}},
		{"incorrect pwd then correct", []string{"email", "", "email", "password"}},
		{"incorrect email then correct", []string{"", "password", "email", "password"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()

			mock := &cltest.MockCountingPrompter{EnteredStrings: test.enteredStrings}
			tai := cmd.NewPromptingAPIInitializer(mock)

			user, err := tai.Initialize(store)
			assert.NoError(t, err)
			assert.Equal(t, len(test.enteredStrings), mock.Count)
			assert.Equal(t, "", user.SessionID)

			persistedUser, err := store.FindUser()
			assert.NoError(t, err)

			assert.Equal(t, user.Email, persistedUser.Email)
			assert.Equal(t, user.HashedPassword, persistedUser.HashedPassword)
			assert.Equal(t, "", persistedUser.SessionID)
		})
	}
}

func TestTerminalAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	initialUser := cltest.MustUser(cltest.UserEmail, cltest.Password)
	require.NoError(t, store.Save(&initialUser))

	mock := &cltest.MockCountingPrompter{}
	tai := cmd.NewPromptingAPIInitializer(mock)

	user, err := tai.Initialize(store)
	assert.NoError(t, err)
	assert.Equal(t, 0, mock.Count)

	assert.Equal(t, initialUser.Email, user.Email)
	assert.Equal(t, initialUser.HashedPassword, user.HashedPassword)
	assert.Equal(t, "", user.SessionID)
}
