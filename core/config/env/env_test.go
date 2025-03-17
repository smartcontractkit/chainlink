package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPlugin(t *testing.T) {
	for _, tt := range []struct {
		name string
		kind string
		exp  Plugin
	}{
		{"lower", "foo", Plugin{Cmd: "CL_FOO_CMD", Env: "CL_FOO_ENV", CmdDefault: "chainlink-foo"}},
		{"upper", "BAR", Plugin{Cmd: "CL_BAR_CMD", Env: "CL_BAR_ENV", CmdDefault: "chainlink-bar"}},
		{"mixed", "Baz", Plugin{Cmd: "CL_BAZ_CMD", Env: "CL_BAZ_ENV", CmdDefault: "chainlink-baz"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlugin(tt.kind)
			require.Equal(t, tt.exp, got)
		})
	}
}
