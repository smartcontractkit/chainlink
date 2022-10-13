package cfgtest

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed dump/empty-strings.env
var emptyStringsEnv string

func Clearenv(t *testing.T) {
	for _, kv := range strings.Split(emptyStringsEnv, "\n") {
		if strings.TrimSpace(kv) == "" {
			continue
		}
		i := strings.Index(kv, "=")
		require.NotEqual(t, -1, i, "invalid kv: %s", kv)
		os.Unsetenv(kv[:i])
	}
}
