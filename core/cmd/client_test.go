package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerminalCookieAuthenticator_AuthenticateWithoutSession(t *testing.T) {
	t.Parallel()

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
			config := orm.NewConfig()

			sr := models.SessionRequest{Email: test.email, Password: test.pwd}
			store := &cmd.MemoryCookieStore{}
			tca := cmd.NewSessionCookieAuthenticator(config, store)
			cookie, err := tca.Authenticate(sr)

			assert.Error(t, err)
			assert.Nil(t, cookie)
			cookie, err = store.Retrieve()
			assert.NoError(t, err)
			assert.Nil(t, cookie)
		})
	}
}

func TestTerminalCookieAuthenticator_AuthenticateWithSession(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBlockByNumber,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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
			store := &cmd.MemoryCookieStore{}
			tca := cmd.NewSessionCookieAuthenticator(app.Config.Config, store)
			cookie, err := tca.Authenticate(sr)

			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, cookie)

				cookie, err = store.Retrieve()
				assert.NoError(t, err)
				assert.Nil(t, cookie)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cookie)

				retrievedCookie, err := store.Retrieve()
				assert.NoError(t, err)
				assert.Equal(t, cookie, retrievedCookie)
			}
		})
	}
}

func TestDiskCookieStore_Retrieve(t *testing.T) {
	t.Parallel()

	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config := tc.Config

	tests := []struct {
		name      string
		rootDir   string
		wantError bool
	}{
		{"missing", config.RootDir(), true},
		{"correct fixture", "../internal/fixtures", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.Set("ROOT", test.rootDir)
			store := cmd.DiskCookieStore{Config: config}
			cookie, err := store.Retrieve()
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
	t.Parallel()

	tests := []struct {
		name           string
		enteredStrings []string
		isTerminal     bool
		isError        bool
	}{
		{"correct", []string{"good@email.com", "password"}, true, false},
		{"incorrect pwd then correct", []string{"good@email.com", "", "good@email.com", "password"}, true, false},
		{"incorrect email then correct", []string{"", "password", "good@email.com", "password"}, true, false},
		{"not a terminal", []string{}, false, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			mock := &cltest.MockCountingPrompter{EnteredStrings: test.enteredStrings, NotTerminal: !test.isTerminal}
			tai := cmd.NewPromptingAPIInitializer(mock)

			// Remove fixture user
			_, err := store.DeleteUser()
			require.NoError(t, err)

			user, err := tai.Initialize(store)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(test.enteredStrings), mock.Count)

				persistedUser, err := store.FindUser()
				assert.NoError(t, err)

				assert.Equal(t, user.Email, persistedUser.Email)
				assert.Equal(t, user.HashedPassword, persistedUser.HashedPassword)
			}
		})
	}
}

func TestTerminalAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	initialUser := cltest.MustRandomUser()
	require.NoError(t, store.SaveUser(&initialUser))

	mock := &cltest.MockCountingPrompter{}
	tai := cmd.NewPromptingAPIInitializer(mock)

	user, err := tai.Initialize(store)
	assert.NoError(t, err)
	assert.Equal(t, 0, mock.Count)

	assert.Equal(t, initialUser.Email, user.Email)
	assert.Equal(t, initialUser.HashedPassword, user.HashedPassword)
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
			store, cleanup := cltest.NewStore(t)
			// Clear out fixture user
			store.DeleteUser()
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
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{"correct", "internal/fixtures/apicredentials", false},
		{"no file", "", true},
		{"incorrect file", "/tmp/doesnotexist", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tfi := cmd.NewFileAPIInitializer(test.file)
			user, err := tfi.Initialize(store)
			assert.NoError(t, err)
			assert.Equal(t, cltest.APIEmail, user.Email)
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
		wantError             bool
	}{
		{"empty", "", "", true},
		{"correct file", "../internal/fixtures/apicredentials", cltest.APIEmail, false},
		{"incorrect file", "/tmp/dontexist", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr, err := builder.Build(test.file)
			assert.Equal(t, test.wantEmail, sr.Email)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
