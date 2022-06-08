// Note: These tests are isolated in a separate package and never run in parallel, since they modify they environment.

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

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

func TestConfig_loadChainsAndNodes(t *testing.T) {
	//TODO
}

func TestConfig_loadLegacyEVMEnv(t *testing.T) {
	//TODO
}

func TestConfig_loadLegacyCoreEnv(t *testing.T) {
	//TODO
}

//go:embed testdata/dump/*.toml
//go:embed testdata/dump/*.env
var dumpTestFiles embed.FS

func TestChainlinkApplication_ConfigDump(t *testing.T) {
	//TODO add db later - EVM_NODES insufficient...
	fes, err := dumpTestFiles.ReadDir("testdata/dump")
	require.NoError(t, err)
	for _, fe := range fes {
		if fe.IsDir() {
			continue
		}
		if filepath.Ext(fe.Name()) != ".toml" {
			continue
		}
		name := strings.TrimSuffix(fe.Name(), ".toml")
		t.Run(name, func(t *testing.T) {
			exp, err := dumpTestFiles.ReadFile(filepath.Join("testdata/dump", fe.Name()))
			require.NoError(t, err)

			env, err := dumpTestFiles.ReadFile(filepath.Join("testdata/dump", name) + ".env")
			require.NoError(t, err)

			os.Clearenv()

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
				require.NoError(t, os.Setenv(k, v))
			}
			var app chainlink.ChainlinkApplication
			got, err := app.ConfigDump(testutils.TestCtx(t))
			require.NoError(t, err)
			assert.Equal(t, string(exp), got, diff.Diff(string(exp), got))
		})
	}
}
