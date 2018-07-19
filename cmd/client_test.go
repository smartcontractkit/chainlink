package cmd_test

import (
	"errors"
	"path"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
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
		{"bad pwd", cltest.APIEmail, "mostcommonwrongpwdever"},
		{"bad both", "notreal", "mostcommonwrongpwdever"},
		{"correct", cltest.APIEmail, cltest.Password},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := models.SessionRequest{Email: test.email, Password: test.pwd}
			tca := cmd.NewSessionCookieAuthenticator(app.Config)
			cookie, err := tca.Authenticate(sr)

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
		{"bad pwd", cltest.APIEmail, "mostcommonwrongpwdever", true},
		{"bad both", "notreal", "mostcommonwrongpwdever", true},
		{"success", cltest.APIEmail, cltest.Password, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := models.SessionRequest{Email: test.email, Password: test.pwd}
			tca := cmd.NewSessionCookieAuthenticator(app.Config)
			cookie, err := tca.Authenticate(sr)

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
			tca := cmd.NewSessionCookieAuthenticator(config)
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
			assert.Empty(t, user.SessionID)

			persistedUser, err := store.FindUser()
			assert.NoError(t, err)

			assert.Equal(t, user.Email, persistedUser.Email)
			assert.Equal(t, user.HashedPassword, persistedUser.HashedPassword)
			assert.Empty(t, persistedUser.SessionID)
		})
	}
}

func TestTerminalAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	initialUser := cltest.MustUser(cltest.APIEmail, cltest.Password)
	require.NoError(t, store.Save(&initialUser))

	mock := &cltest.MockCountingPrompter{}
	tai := cmd.NewPromptingAPIInitializer(mock)

	user, err := tai.Initialize(store)
	assert.NoError(t, err)
	assert.Equal(t, 0, mock.Count)

	assert.Equal(t, initialUser.Email, user.Email)
	assert.Equal(t, initialUser.HashedPassword, user.HashedPassword)
	assert.Empty(t, user.SessionID)
}

func TestFileAPIInitializer_InitializeWithoutAPIUser(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false},
		{"no file", "", true},
		{"incorrect file", "/tmp/doesnotexist", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()

			tfi := cmd.NewFileAPIInitializer(test.file)
			user, err := tfi.Initialize(store)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, cltest.APIEmail, user.Email)
				persistedUser, err := store.FindUser()
				assert.NoError(t, err)
				assert.Equal(t, persistedUser.Email, user.Email)
			}
		})
	}
}

func TestFileAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	initialUser := cltest.MustUser(cltest.APIEmail, cltest.Password)
	require.NoError(t, store.Save(&initialUser))

	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false},
		{"no file", "", true},
		{"incorrect file", "/tmp/doesnotexist", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tfi := cmd.NewFileAPIInitializer(test.file)
			user, err := tfi.Initialize(store)
			assert.NoError(t, err)
			assert.Equal(t, initialUser.Email, user.Email)
		})
	}
}

func TestPromptingSessionRequestBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		email, pwd string
	}{
		{"correct@input.com", "mypwd"},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			enteredStrings := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{EnteredStrings: enteredStrings}
			builder := cmd.NewPromptingSessionRequestBuilder(prompter)

			sr, err := builder.Build("")
			require.NoError(t, err)
			assert.Equal(t, test.email, sr.Email)
			assert.Equal(t, test.pwd, sr.Password)
		})
	}
}

func TestFileSessionRequestBuilder(t *testing.T) {
	t.Parallel()

	builder := cmd.NewFileSessionRequestBuilder()
	tests := []struct {
		name, file, wantEmail string
		wantError             error
	}{
		{"empty", "", "", errors.New("No API user credential file was passed")},
		{"correct file", "../internal/fixtures/apicredentials", "email@test.net", nil},
		{"incorrect file", "/tmp/dontexist", "", errors.New("open /tmp/dontexist: no such file or directory")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr, err := builder.Build(test.file)
			assert.Equal(t, test.wantEmail, sr.Email)
			if test.wantError != nil {
				assert.Equal(t, test.wantError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
