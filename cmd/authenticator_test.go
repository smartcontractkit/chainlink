package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestTerminalAuthenticatorWithAcctWithPwd(t *testing.T) {
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
			auth := cmd.TerminalAuthenticator{func(i int) {
				exited = true
				rval = i
			}}

			auth.Authenticate(app.Store, test.password)
			assert.Equal(t, test.wantExited, exited)
			assert.Equal(t, test.wantRval, rval)
		})
	}
}
