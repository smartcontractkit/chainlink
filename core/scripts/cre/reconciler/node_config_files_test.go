package reconciler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
)

func TestNodeIdentity(t *testing.T) {
	t.Parallel()

	require.Equal(t, "zone-a-gateway/gateway-0", domain.NodeIdentity("zone-a-gateway", "gateway-0"))
}

func TestBuildNodeConfigFileMap_DuplicateNamesDifferentNamespaces(t *testing.T) {
	t.Parallel()

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", Namespace: "zone-a", ConfigFile: "/a/dev.yaml"},
			{Name: "node-0", Namespace: "zone-b", ConfigFile: "/b/dev.yaml"},
		},
	}

	mapping, err := cv.BuildNodeConfigFileMap()
	require.NoError(t, err)
	require.Len(t, mapping, 2)
	require.Equal(t, "/a/dev.yaml", mapping[domain.NodeIdentity("zone-a", "node-0")])
	require.Equal(t, "/b/dev.yaml", mapping[domain.NodeIdentity("zone-b", "node-0")])
}

func TestGroupPatchesByConfigFile_SplitsByTargetFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zoneA := filepath.Join(dir, "zone-a", "dev.yaml")
	zoneGw := filepath.Join(dir, "zone-a-gateway", "dev.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(zoneA), 0755))
	require.NoError(t, os.MkdirAll(filepath.Dir(zoneGw), 0755))
	require.NoError(t, os.WriteFile(zoneA, []byte(`
chainlink-node:
  instances:
    node-0: {}
`), 0600))
	require.NoError(t, os.WriteFile(zoneGw, []byte(`
chainlink-node:
  instances:
    gateway-0:
      nodeType: gateway
`), 0600))

	mapping := map[string]string{
		domain.NodeIdentity("zone-a", "node-0"):            zoneA,
		domain.NodeIdentity("zone-a-gateway", "gateway-0"): zoneGw,
	}
	patches := []nodeTOMLPatch{
		{Namespace: "zone-a", Name: "node-0", TOML: "toml-node-0"},
		{Namespace: "zone-a-gateway", Name: "gateway-0", TOML: "toml-gateway-0"},
	}

	grouped, err := groupPatchesByConfigFile(mapping, patches)
	require.NoError(t, err)
	require.Len(t, grouped, 2)
	require.Equal(t, "toml-node-0", grouped[zoneA]["node-0"])
	require.Equal(t, "toml-gateway-0", grouped[zoneGw]["gateway-0"])
}

func TestGroupPatchesByConfigFile_MissingMapping(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zoneGw := filepath.Join(dir, "dev.yaml")
	require.NoError(t, os.WriteFile(zoneGw, []byte("chainlink-node:\n  instances:\n    gateway-0: {}\n"), 0600))

	_, err := groupPatchesByConfigFile(map[string]string{
		domain.NodeIdentity("zone-a", "node-0"): zoneGw,
	}, []nodeTOMLPatch{
		{Namespace: "zone-a-gateway", Name: "gateway-0", TOML: "x"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "zone-a-gateway/gateway-0")
	require.Contains(t, err.Error(), "no config file mapping")
}

func TestPatchChartValues_GatewayNodeInOwnFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	gwYAML := filepath.Join(dir, "dev.yaml")
	require.NoError(t, os.WriteFile(gwYAML, []byte(`
chainlink-node:
  instances:
    gateway-0:
      nodeType: gateway
`), 0600))

	err := infra.PatchChartValues(gwYAML, map[string]string{
		"gateway-0": "[Capabilities]\nDonID = 'workflow'\n",
	})
	require.NoError(t, err)

	updated, err := os.ReadFile(gwYAML)
	require.NoError(t, err)
	require.Contains(t, string(updated), "30-cre")
	require.Contains(t, string(updated), "DonID = 'workflow'")
}

func TestRefreshNodeConfigFiles_PersistsMapping(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	zoneGw := filepath.Join(dir, "zone-a-gateway", "dev.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(zoneGw), 0755))
	require.NoError(t, os.WriteFile(zoneGw, []byte("chainlink-node:\n  instances:\n    gateway-0: {}\n"), 0600))

	r := &Reconciler{
		cv: &domain.ChartValues{
			Nodes: []domain.ChartNodeInfo{
				{Name: "gateway-0", Namespace: "zone-a-gateway", ConfigFile: zoneGw},
			},
		},
		state: &domain.StateFile{},
	}
	require.NoError(t, r.refreshNodeConfigFiles())
	require.Equal(t, zoneGw, r.state.NodeConfigFiles[domain.NodeIdentity("zone-a-gateway", "gateway-0")])
}
