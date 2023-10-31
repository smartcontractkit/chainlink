package configtest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

// WriteTOMLFile is used in tests to output toml config types to a toml string.
// Secret values are redacted.
func WriteTOMLFile(t *testing.T, contents any, fileName string) string {
	d := t.TempDir()
	p := filepath.Join(d, fileName)

	b, err := toml.Marshal(contents)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(p, b, 0600))
	return p
}
