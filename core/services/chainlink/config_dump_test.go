// Note: These tests are isolated in a separate package and never run in parallel, since they modify the environment.

package chainlink_test

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
)

//go:embed cfgtest/dump/*.env
//go:embed cfgtest/dump/*.json
//go:embed cfgtest/dump/*.toml
var dumpTestFiles embed.FS

func TestChainlinkApplication_ConfigDump(t *testing.T) {
	dir := "cfgtest/dump"
	fes, err := dumpTestFiles.ReadDir(dir)
	require.NoError(t, err)
	for _, fe := range fes {
		if filepath.Ext(fe.Name()) != ".toml" {
			continue
		}
		name := strings.TrimSuffix(fe.Name(), ".toml")
		t.Run(name, func(t *testing.T) {
			exp, err := dumpTestFiles.ReadFile(filepath.Join(dir, fe.Name()))
			require.NoError(t, err)

			env, err := dumpTestFiles.ReadFile(filepath.Join(dir, name+".env"))
			if !os.IsNotExist(err) { // optional
				require.NoError(t, err)
			}

			chainsJSON, err := dumpTestFiles.ReadFile(filepath.Join(dir, name+".json"))
			if !os.IsNotExist(err) { // optional
				require.NoError(t, err)
			}

			cfgtest.Clearenv(t)

			seen := map[string]struct{}{}
			for _, kv := range strings.Split(string(env), "\n") {
				if strings.TrimSpace(kv) == "" {
					continue
				}
				i := strings.Index(kv, "=")
				require.NotEqual(t, -1, i, "invalid kv: %s", kv)
				k, v := kv[:i], kv[i+1:]
				_, ok := seen[k]
				require.False(t, ok, "duplicate key: %s", k)
				seen[k] = struct{}{}
				t.Setenv(k, v)
			}

			got, err := chainlink.FakeConfigDump(chainsJSON)
			require.NoError(t, err)
			assert.Equal(t, string(exp), got, diff.Diff(string(exp), got))
		})
	}
}
