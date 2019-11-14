package cmd_test

import (
	"testing"

	"chainlink/core/cmd"
	"chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
)

func TestTerminalKeyStoreAuthenticator_WithNoAcctNoPwdCreatesAccount(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	prompt := &cltest.MockCountingPrompter{
		T: t,
		EnteredStrings: []string{
			cltest.Password,
			"wrongconfirmation",
			cltest.Password,
			cltest.Password,
		},
	}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	assert.False(t, store.KeyStore.HasAccounts())
	_, err := auth.Authenticate(store, "")
	assert.NoError(t, err)
	assert.Equal(t, 4, prompt.Count)
	assert.Len(t, store.KeyStore.Accounts(), 1)
}

func TestTerminalKeyStoreAuthenticator_WithNoAcctWithInitialPwdCreatesAcct(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}

	assert.Len(t, store.KeyStore.Accounts(), 0)
	_, err := auth.Authenticate(store, "somepassword")
	assert.NoError(t, err)
	assert.True(t, store.KeyStore.HasAccounts())
	assert.Len(t, store.KeyStore.Accounts(), 1)
}

func TestTerminalKeyStoreAuthenticator_WithAcctNoInitialPwdPromptLoop(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// prompt loop tries all in array
	prompt := &cltest.MockCountingPrompter{
		T:              t,
		EnteredStrings: []string{"wrongpassword", cltest.Password, cltest.Password, cltest.Password},
	}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	_, err := auth.Authenticate(store, "")
	assert.NoError(t, err)
	assert.Equal(t, 4, prompt.Count)
}

func TestTerminalKeyStoreAuthenticator_WithAcctAndPwd(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
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
			auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
			_, err := auth.Authenticate(store, test.password)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}
