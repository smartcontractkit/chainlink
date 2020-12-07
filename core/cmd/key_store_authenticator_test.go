package cmd_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestTerminalKeyStoreAuthenticator_WithNoAcctNoPwdCreatesAccount(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	kst := new(mocks.KeyStoreInterface)
	kst.On("HasAccounts").Return(false)
	kst.On("Unlock", cltest.Password).Return(nil)
	kst.On("NewAccount").Return(accounts.Account{}, nil)
	store.KeyStore = kst

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

	kst.AssertExpectations(t)
}

func TestTerminalKeyStoreAuthenticator_WithNoAcctWithInitialPwdCreatesAcct(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	kst := new(mocks.KeyStoreInterface)
	kst.On("HasAccounts").Return(false)
	kst.On("Unlock", "somepassword").Return(nil)
	kst.On("NewAccount").Return(accounts.Account{}, nil)
	kst.On("Accounts").Return([]accounts.Account{})
	store.KeyStore = kst
	defer cleanup()

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}

	assert.Len(t, store.KeyStore.Accounts(), 0)
	_, err := auth.Authenticate(store, "somepassword")
	assert.NoError(t, err)

	kst.AssertExpectations(t)
}

func TestTerminalKeyStoreAuthenticator_WithAcctNoInitialPwdPromptLoop(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// prompt loop tries all in array
	prompt := &cltest.MockCountingPrompter{
		T:              t,
		EnteredStrings: []string{"wrongpassword", "wrongagain", cltest.Password},
	}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	_, err := auth.Authenticate(store, "")
	assert.NoError(t, err)
	assert.Equal(t, 3, prompt.Count)
}

func TestTerminalKeyStoreAuthenticator_WithAcctAndPwd(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	tests := []struct {
		password  string
		wantError bool
	}{
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
