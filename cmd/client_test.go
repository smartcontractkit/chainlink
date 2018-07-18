package cmd_test

import (
	"path"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
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
