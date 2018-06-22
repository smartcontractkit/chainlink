package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestTerminalAuthenticatorWithNoAcctNoPwdCreatesAccount(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	prompt := &cltest.MockCountingPrompt{EnteredStrings: []string{
		cltest.Password, "wrongconfirmation", cltest.Password, cltest.Password,
	}}

	auth := cmd.TerminalAuthenticator{Prompter: prompt}
	assert.False(t, app.Store.KeyStore.HasAccounts())
	assert.NoError(t, auth.Authenticate(app.Store, ""))
	assert.Equal(t, 4, prompt.Count)
	assert.Equal(t, 1, len(app.Store.KeyStore.Accounts()))
}

func TestTerminalAuthenticatorWithNoAcctWithInitialPwdCreatesAcct(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	auth := cmd.TerminalAuthenticator{Prompter: &cltest.MockCountingPrompt{}}

	assert.Equal(t, 0, len(app.Store.KeyStore.Accounts()))
	assert.NoError(t, auth.Authenticate(app.Store, "somepassword"))
	assert.True(t, app.Store.KeyStore.HasAccounts())
	assert.Equal(t, 1, len(app.Store.KeyStore.Accounts()))
}

func TestTerminalAuthenticatorWithAcctNoInitialPwdPromptLoop(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	// prompt loop tries all in array
	prompt := &cltest.MockCountingPrompt{
		EnteredStrings: []string{"wrongpassword", cltest.Password},
	}

	auth := cmd.TerminalAuthenticator{Prompter: prompt}
	assert.NoError(t, auth.Authenticate(app.Store, ""))
	assert.Equal(t, 2, prompt.Count)
}

func TestTerminalAuthenticatorWithAcctAndPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	tests := []struct {
		password  string
		wantError bool
	}{
		{cltest.Password, false},
		{"wrongpassword", true},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			auth := cmd.TerminalAuthenticator{Prompter: &cltest.MockCountingPrompt{}}
			err := auth.Authenticate(app.Store, test.password)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}
