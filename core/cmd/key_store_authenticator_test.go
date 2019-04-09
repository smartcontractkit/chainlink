package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
)

func TestTerminalKeyStoreAuthenticator_WithNoAcctNoPwdCreatesAccount(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	prompt := &cltest.MockCountingPrompter{EnteredStrings: []string{
		cltest.Password, "wrongconfirmation", cltest.Password, cltest.Password,
	}}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	assert.False(t, app.Store.KeyStore.HasAccounts())
	_, err := auth.Authenticate(app.Store, "")
	assert.NoError(t, err)
	assert.Equal(t, 4, prompt.Count)
	assert.Equal(t, 1, len(app.Store.KeyStore.Accounts()))
}

func TestTerminalKeyStoreAuthenticator_WithNoAcctWithInitialPwdCreatesAcct(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{}}

	assert.Equal(t, 0, len(app.Store.KeyStore.Accounts()))
	_, err := auth.Authenticate(app.Store, "somepassword")
	assert.NoError(t, err)
	assert.True(t, app.Store.KeyStore.HasAccounts())
	assert.Equal(t, 1, len(app.Store.KeyStore.Accounts()))
}

func TestTerminalKeyStoreAuthenticator_WithAcctNoInitialPwdPromptLoop(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
	defer cleanup()

	// prompt loop tries all in array
	prompt := &cltest.MockCountingPrompter{
		EnteredStrings: []string{"wrongpassword", cltest.Password},
	}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	_, err := auth.Authenticate(app.Store, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, prompt.Count)
}

func TestTerminalKeyStoreAuthenticator_WithAcctAndPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey()
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
			auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{}}
			_, err := auth.Authenticate(app.Store, test.password)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}
