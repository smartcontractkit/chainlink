package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestTerminalAuthenticatorWithCorrectPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	var exited bool
	var rval int
	auth := cmd.TerminalAuthenticator{func(i int) {
		exited = true
		rval = i
	}}

	auth.Authenticate(app.Store, cltest.Password)
	assert.False(t, exited)
}

func TestTerminalAuthenticatorWithWrongPwd(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	var exited bool
	var rval int
	auth := cmd.TerminalAuthenticator{func(i int) {
		exited = true
		rval = i
	}}

	auth.Authenticate(app.Store, "wrongpassword")
	assert.True(t, exited)
	assert.Equal(t, 1, rval)
}
