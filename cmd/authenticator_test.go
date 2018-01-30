package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestTerminalAuthenticatorWithNoAcctWithInitialPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	var exited bool
	var rval int
	auth := cmd.TerminalAuthenticator{&cltest.MockCountingPrompt{}, func(i int) {
		exited = true
		rval = i
	}}

	auth.Authenticate(app.Store, "somepassword")
	assert.Equal(t, true, exited)
	assert.Equal(t, 1, rval)
}

func TestTerminalAuthenticatorWithAcctNoInitialPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	tests := []struct {
		password string
		prompts  int
	}{
		{cltest.Password, 1},
		{"wrongpassword", 2},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			var exited bool
			prompt := &cltest.MockCountingPrompt{
				EnteredStrings: []string{test.password, cltest.Password},
			}

			auth := cmd.TerminalAuthenticator{prompt, func(i int) { exited = true }}

			auth.Authenticate(app.Store, "")
			assert.False(t, exited)
			assert.Equal(t, test.prompts, prompt.Count)
		})
	}
}

func TestTerminalAuthenticatorWithAcctAndPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	tests := []struct {
		password   string
		wantExited bool
		wantRval   int
	}{
		{cltest.Password, false, 0},
		{"wrongpassword", true, 1},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			var exited bool
			var rval int
			auth := cmd.TerminalAuthenticator{&cltest.MockCountingPrompt{}, func(i int) {
				exited = true
				rval = i
			}}

			auth.Authenticate(app.Store, test.password)
			assert.Equal(t, test.wantExited, exited)
			assert.Equal(t, test.wantRval, rval)
		})
	}
}
