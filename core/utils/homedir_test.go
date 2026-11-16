package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandHomeDir(t *testing.T) {
	t.Setenv("HOME", filepath.Join(t.TempDir(), "home"))
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty", input: "", want: ""},
		{name: "absolute", input: "/var/log/chainlink", want: "/var/log/chainlink"},
		{name: "relative", input: "relative/path", want: "relative/path"},
		{name: "tilde", input: "~", want: home},
		{name: "tilde path", input: "~/chainlink", want: filepath.Join(home, "chainlink")},
		{name: "user specific", input: "~someone/chainlink", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ExpandHomeDir(test.input)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
