package plugins

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnvFile(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := ParseEnvFile("testdata/valid.env")
		require.NoError(t, err)
		require.Equal(t, []string{"GOMEMLIMIT=1MiB"}, got)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseEnvFile("testdata/invalid.env")
		require.Error(t, err)
	})
	t.Run("missing", func(t *testing.T) {
		_, err := ParseEnvFile("testdata/missing.env")
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
