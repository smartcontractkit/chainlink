package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerminalKeyStoreAuthenticator_WithNoAcctNoPwdCreatesKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	kst := cltest.NewKeyStore(t, store.DB).Eth

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
	has, err := kst.HasDBSendingKeys()
	require.NoError(t, err)
	assert.False(t, has)
	_, err = auth.AuthenticateEthKey(kst, "")
	assert.NoError(t, err)
	assert.Equal(t, 4, prompt.Count)
}

func TestTerminalKeyStoreAuthenticator_WithNoAcctWithInitialPwdCreatesAcct(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	kst := cltest.NewKeyStore(t, store.DB).Eth

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}

	kst.Unlock(cltest.Password)
	sendingKeys, err := kst.SendingKeys()
	require.NoError(t, err)
	assert.Len(t, sendingKeys, 0)
	_, err = auth.AuthenticateEthKey(kst, cltest.Password)
	assert.NoError(t, err)
	sendingKeys, err = kst.SendingKeys()
	require.NoError(t, err)
	assert.Len(t, sendingKeys, 1)
}

func TestTerminalKeyStoreAuthenticator_WithAcctNoInitialPwdPromptLoop(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	// prompt loop tries all in array
	prompt := &cltest.MockCountingPrompter{
		T:              t,
		EnteredStrings: []string{"wrongpassword", "wrongagain", cltest.Password},
	}

	auth := cmd.TerminalKeyStoreAuthenticator{Prompter: prompt}
	_, err := auth.AuthenticateEthKey(ethKeyStore, "")
	assert.NoError(t, err)
	assert.Equal(t, 3, prompt.Count)
}

func TestTerminalKeyStoreAuthenticator_WithAcctAndPwd(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

	cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	tests := []struct {
		password  string
		wantError bool
	}{
		{"wrongpassword", true},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			auth := cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{T: t}}
			_, err := auth.AuthenticateEthKey(ethKeyStore, test.password)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}

func TestTerminalKeyStoreAuthenticator_ValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name               string
		failingPassword    string
		succeedingPassword string
		errorString        string
	}{
		{
			"not long enough",
			"password",
			"passwordpassword",
			"must be longer than 12 characters",
		},
		{
			"not enough lowercase",
			"paSSWORD",
			"password",
			"must contain at least 3 lowercase characters",
		},
		{
			"not enough uppercase",
			"PAssword",
			"PASsword",
			"must contain at least 3 uppercase characters",
		},
		{
			"not enough numbers",
			"password",
			"password123",
			"must contain at least 3 numbers",
		},
		{
			"not enough symbols",
			"password",
			"password!@#",
			"must contain at least 3 symbols",
		},
		{
			"identical consecutive characters",
			"paaaasword",
			"password",
			"must not contain more than 3 identical consecutive characters",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store, cleanup := cltest.NewStore(t)
			defer cleanup()
			ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth

			auth := cmd.TerminalKeyStoreAuthenticator{}
			err := auth.ExportedValidatePasswordStrength(ethKeyStore, test.failingPassword)
			require.Error(t, err)
			require.Contains(t, err.Error(), test.errorString)
			err = auth.ExportedValidatePasswordStrength(ethKeyStore, test.succeedingPassword)
			if err != nil {
				require.NotContains(t, err.Error(), test.errorString)
			}
		})
	}
}
