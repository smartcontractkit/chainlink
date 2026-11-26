package helpers

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const testEnvironmentDir = "../../../core/scripts/cre/environment"

func topologyPath(name string) string {
	return filepath.Join(testEnvironmentDir, "configs", name)
}

func TestModuleCacheSettingsFromTopology(t *testing.T) {
	t.Parallel()

	capacity := moduleCacheSettingsFromTopology(t, []string{topologyPath("workflow-gateway-don-capacity-test.toml")}, "workflow")
	require.Len(t, capacity, 1)
	require.NotNil(t, capacity[0].Enabled)
	assert.True(t, *capacity[0].Enabled)
	require.NotNil(t, capacity[0].MaxLoaded)
	assert.Equal(t, 50, *capacity[0].MaxLoaded)
	require.NotNil(t, capacity[0].IdleTimeout)
	assert.Equal(t, "1m", *capacity[0].IdleTimeout)

	plain := moduleCacheSettingsFromTopology(t, []string{topologyPath("workflow-gateway-don.toml")}, "workflow")
	require.Len(t, plain, 1)
	assert.Nil(t, plain[0].Enabled, "default topology leaves the cache at the node default (off)")

	// later files win, as when `env start` merges CTF_CONFIGS
	merged := moduleCacheSettingsFromTopology(t, []string{
		topologyPath("workflow-gateway-don-capacity-test.toml"),
		topologyPath("workflow-gateway-don-cache-test.toml"),
	}, "workflow")
	require.Len(t, merged, 1)
	require.NotNil(t, merged[0].MaxLoaded)
	assert.Equal(t, 1, *merged[0].MaxLoaded)

	gateway := moduleCacheSettingsFromTopology(t, []string{topologyPath("workflow-gateway-don.toml")}, "bootstrap-gateway")
	assert.Len(t, gateway, 1)
}

func TestRequestedTopologyFiles(t *testing.T) {
	tconf := &ttypes.TestConfig{
		EnvironmentDirPath:    testEnvironmentDir,
		EnvironmentConfigPath: topologyPath("workflow-gateway-don-capacity-test.toml"),
	}

	t.Setenv("CTF_CONFIGS", "")
	assert.Equal(t, []string{tconf.EnvironmentConfigPath}, requestedTopologyFiles(tconf))

	t.Setenv("CTF_CONFIGS", "configs/a.toml, /abs/b.toml")
	assert.Equal(t, []string{filepath.Join(testEnvironmentDir, "configs/a.toml"), "/abs/b.toml"}, requestedTopologyFiles(tconf))
}

func TestModuleCacheStatusFromConfigs(t *testing.T) {
	t.Parallel()

	// shaped like the files a local CRE node container starts with
	const base = "[Log]\nLevel = 'info'\n"
	const frameworkOverrides = "[Capabilities.WorkflowRegistry]\nAddress = '0x1'\n\n[Capabilities.WorkflowRegistry.ModuleCache]\n"
	const userOverrides = "[Capabilities.WorkflowRegistry.ModuleCache]\nEnabled = true\nMaxLoaded = 50\nIdleTimeout = '1m'\n"

	got, err := moduleCacheStatusFromConfigs("workflow-node0", []string{base, frameworkOverrides, userOverrides})
	require.NoError(t, err)
	assert.Equal(t, ModuleCacheStatus{
		Node: "workflow-node0", Enabled: true, IdleEviction: true, IdleTimeout: time.Minute, MaxLoaded: 50,
	}, got)

	got, err = moduleCacheStatusFromConfigs("workflow-node0", []string{base, frameworkOverrides, ""})
	require.NoError(t, err)
	assert.Equal(t, ModuleCacheStatus{
		Node: "workflow-node0", IdleEviction: true, IdleTimeout: 10 * time.Minute, MaxLoaded: 200,
	}, got, "nothing set: node defaults")

	got, err = moduleCacheStatusFromConfigs("workflow-node0", []string{userOverrides, "[Capabilities.WorkflowRegistry.ModuleCache]\nEnabled = false\n"})
	require.NoError(t, err)
	assert.False(t, got.Enabled, "later files win")

	_, err = moduleCacheStatusFromConfigs("workflow-node0", []string{"[Capabilities.WorkflowRegistry.ModuleCache]\nIdleTimeout = 'soon'\n"})
	require.ErrorContains(t, err, "invalid ModuleCache.IdleTimeout")
}

func TestModuleCacheMismatch(t *testing.T) {
	t.Parallel()

	enabled, disabled := true, false
	maxLoaded := 50
	idleTimeout := "1m"
	wantOn := moduleCacheSettings{Enabled: &enabled, MaxLoaded: &maxLoaded, IdleTimeout: &idleTimeout}
	on := ModuleCacheStatus{Node: "n", Enabled: true, IdleEviction: true, IdleTimeout: time.Minute, MaxLoaded: 50}

	for _, tc := range []struct {
		name    string
		got     ModuleCacheStatus
		want    moduleCacheSettings
		wantErr string
	}{
		{name: "matches", got: on, want: wantOn},
		{name: "topology on, node off", got: ModuleCacheStatus{Node: "n"}, want: wantOn, wantErr: "enabled=false, topology wants enabled=true"},
		{name: "topology unset, node on", got: on, want: moduleCacheSettings{}, wantErr: "enabled=true, topology wants enabled=false"},
		{name: "topology off, node off", got: ModuleCacheStatus{Node: "n"}, want: moduleCacheSettings{Enabled: &disabled}},
		{name: "different MaxLoaded", got: ModuleCacheStatus{Node: "n", Enabled: true, IdleTimeout: time.Minute, MaxLoaded: 200}, want: wantOn, wantErr: "MaxLoaded=200, topology wants 50"},
		{name: "different IdleTimeout", got: ModuleCacheStatus{Node: "n", Enabled: true, IdleTimeout: 10 * time.Minute, MaxLoaded: 50}, want: wantOn, wantErr: "IdleTimeout=10m0s, topology wants 1m0s"},
		{name: "unset fields are not compared", got: ModuleCacheStatus{Node: "n", Enabled: true, MaxLoaded: 7}, want: moduleCacheSettings{Enabled: &enabled}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := moduleCacheMismatch(tc.got, tc.want)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}
