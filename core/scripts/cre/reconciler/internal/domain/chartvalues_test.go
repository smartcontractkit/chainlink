package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempYAML(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0600))
}

func TestNodeInternalHost_PerNodeNamespace(t *testing.T) {
	t.Parallel()

	cv := &ChartValues{
		Namespace: "zone-c",
		Nodes: []ChartNodeInfo{
			{Name: "node-0", NodeType: RoleStandard, Namespace: "zone-c"},
			{Name: "node-bt-1", NodeType: RoleBootstrap, Namespace: "zone-c-bt"},
		},
	}

	require.Equal(t, "node-0.zone-c.svc.cluster.local", cv.NodeInternalHost("node-0"))
	require.Equal(t, "node-bt-1.zone-c-bt.svc.cluster.local", cv.NodeInternalHost("node-bt-1"),
		"must use the bootstrap node's own namespace, not the chart's primary namespace")
}

// writeGriddleRepo creates a temp dir with griddle.yaml + config files,
// mimicking a real consumer repo structure.
func writeGriddleRepo(t *testing.T, env string, sharedContent, envContent string) string {
	t.Helper()
	repoRoot := t.TempDir()
	configDir := filepath.Join(repoRoot, "deploy", "config", "my-repo")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	writeTempYAML(t, configDir, "shared.yaml", sharedContent)
	// Always write the env file, even if empty
	writeTempYAML(t, configDir, env+".yaml", envContent)

	griddleYAML := fmt.Sprintf(`
deploy:
  %s:
    instances:
      - name: my-repo
        namespace: my-repo-nodeset
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/my-repo/shared.yaml
          - deploy/config/my-repo/%s.yaml
`, env, env)
	writeTempYAML(t, repoRoot, "griddle.yaml", griddleYAML)

	return repoRoot
}

func TestLoadChartValues_Nodes(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  defaults:
    nodeType: standard
  instances:
    node-bt-0:
      nodeType: boot
    node-0: {}
    node-1: {}
    node-2: {}
    node-3: {}
    node-gw-0:
      nodeType: gateway
`, "")

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.Len(t, cv.Nodes, 6)
	require.Equal(t, "my-repo-nodeset", cv.Namespace)

	bt := cv.GetNode("node-bt-0")
	require.NotNil(t, bt)
	require.Equal(t, RoleBootstrap, bt.NodeType)

	gw := cv.GetNode("node-gw-0")
	require.NotNil(t, gw)
	require.Equal(t, RoleGateway, gw.NodeType)

	standard := cv.GetNode("node-0")
	require.NotNil(t, standard)
	require.Equal(t, RoleStandard, standard.NodeType)

	require.Equal(t, "node-bt-0", cv.FindBootstrap())

	gws := cv.FindGatewayNodes()
	require.Len(t, gws, 1)
	require.Equal(t, "node-gw-0", gws[0].Name)
}

func TestLoadChartValues_BootstrapByNaming(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-bt-0: {}
    node-0: {}
`, "")

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.Equal(t, "node-bt-0", cv.FindBootstrap())
}

func TestLoadChartValues_NoBootstrap(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-0: {}
    node-1: {}
`, "")

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.Empty(t, cv.FindBootstrap())
}

func TestLoadChartValues_DeepMerge(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  defaults:
    nodeType: standard
    image:
      repository: chainlink
      tag: "2.52.0"
  instances:
    node-0: {}
`, `
chainlink-node:
  defaults:
    image:
      tag: "2.53.0"
  instances:
    node-1: {}
`)

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.True(t, cv.HasNode("node-0"))
	require.True(t, cv.HasNode("node-1"))
}

func TestLoadChartValues_HasNode(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-0: {}
`, "")

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.True(t, cv.HasNode("node-0"))
	require.False(t, cv.HasNode("node-99"))
}

func TestLoadChartValues_MultipleNodeSetInstances(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	zoneADir := filepath.Join(repoRoot, "deploy", "config", "zone-a")
	zoneGwDir := filepath.Join(repoRoot, "deploy", "config", "zone-a-gateway")
	require.NoError(t, os.MkdirAll(zoneADir, 0755))
	require.NoError(t, os.MkdirAll(zoneGwDir, 0755))

	writeTempYAML(t, zoneADir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "zone-a"
  instances:
    node-0: {}
    node-bt-0:
      nodeType: boot
`)
	writeTempYAML(t, zoneADir, "dev.yaml", `
anvil:
  instances:
    anvil-1337:
      chainID: 1337
      gateway:
        hostnames:
          - anvil-1337.zone-a.dev.internal.griddle.sh
`)
	writeTempYAML(t, zoneGwDir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "zone-a-gateway"
  instances:
    gateway-0:
      nodeType: gateway
`)
	writeTempYAML(t, zoneGwDir, "dev.yaml", "")
	writeTempYAML(t, repoRoot, "griddle.yaml", `
deploy:
  dev:
    instances:
      - name: zone-a
        namespace: zone-a
        path: stack://node-set
        config:
          - deploy/config/zone-a/shared.yaml
          - deploy/config/zone-a/dev.yaml
      - name: zone-a-gateway
        namespace: zone-a-gateway
        path: stack://node-set
        config:
          - deploy/config/zone-a-gateway/shared.yaml
          - deploy/config/zone-a-gateway/dev.yaml
`)

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)
	require.Equal(t, "zone-a", cv.Namespace)
	require.True(t, cv.HasNode("node-0"))
	require.True(t, cv.HasNode("gateway-0"))
	require.Equal(t, "zone-a", cv.GetNodeNamespace("node-0"))
	require.Equal(t, "zone-a-gateway", cv.GetNodeNamespace("gateway-0"))

	node0 := cv.GetNodeInNamespace("zone-a", "node-0")
	require.NotNil(t, node0)
	require.Equal(t, filepath.Join(repoRoot, "deploy/config/zone-a/dev.yaml"), node0.ConfigFile)

	gateway0 := cv.GetNodeInNamespace("zone-a-gateway", "gateway-0")
	require.NotNil(t, gateway0)
	require.Equal(t, filepath.Join(repoRoot, "deploy/config/zone-a-gateway/dev.yaml"), gateway0.ConfigFile)

	mapping, err := cv.BuildNodeConfigFileMap()
	require.NoError(t, err)
	require.Equal(t, node0.ConfigFile, mapping[NodeIdentity("zone-a", "node-0")])
	require.Equal(t, gateway0.ConfigFile, mapping[NodeIdentity("zone-a-gateway", "gateway-0")])
	require.NotEqual(t, mapping[NodeIdentity("zone-a", "node-0")], mapping[NodeIdentity("zone-a-gateway", "gateway-0")])

	gws := cv.FindGatewayNodes()
	require.Len(t, gws, 1)
	require.Equal(t, "gateway-0", gws[0].Name)
}

func TestLoadChartValues_FallbackConfigDir(t *testing.T) {
	t.Parallel()

	// When there's no griddle.yaml, fall back to treating the dir as a config dir
	dir := t.TempDir()
	writeTempYAML(t, dir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-0: {}
`)
	writeTempYAML(t, dir, "dev.yaml", "")

	cv, err := LoadChartValues(dir, "dev")
	require.NoError(t, err)
	require.True(t, cv.HasNode("node-0"))
}

func TestLoadChartValues_EnvNotFound(t *testing.T) {
	t.Parallel()

	// Two environments exist, request a third that doesn't
	repoRoot := t.TempDir()
	configDir := filepath.Join(repoRoot, "deploy", "config", "my-repo")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	writeTempYAML(t, configDir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-0: {}
`)
	writeTempYAML(t, configDir, "dev.yaml", "")
	writeTempYAML(t, configDir, "stage.yaml", "")
	writeTempYAML(t, repoRoot, "griddle.yaml", `
deploy:
  dev:
    instances:
      - name: my-repo
        namespace: ns-dev
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/my-repo/shared.yaml
          - deploy/config/my-repo/dev.yaml
  stage:
    instances:
      - name: my-repo
        namespace: ns-stage
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/my-repo/shared.yaml
          - deploy/config/my-repo/stage.yaml
`)

	_, err := LoadChartValues(repoRoot, "prod")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found in griddle.yaml")
}

func TestLoadChartValues_SingleEnvFallback(t *testing.T) {
	t.Parallel()

	// If only one environment exists in griddle.yaml, use it regardless of env param
	repoRoot := t.TempDir()
	configDir := filepath.Join(repoRoot, "deploy", "config", "my-repo")
	require.NoError(t, os.MkdirAll(configDir, 0755))
	writeTempYAML(t, configDir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  instances:
    node-0: {}
`)
	writeTempYAML(t, repoRoot, "griddle.yaml", `
deploy:
  dev:
    instances:
      - name: my-repo
        namespace: ns
        path: stack://node-set
        version: 0.11.0
        config:
          - deploy/config/my-repo/shared.yaml
`)

	// Request "prod" but only "dev" exists — should use dev
	cv, err := LoadChartValues(repoRoot, "prod")
	require.NoError(t, err)
	require.True(t, cv.HasNode("node-0"))
}

func TestLoadChartValues_DONNameLabel(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	zoneADir := filepath.Join(repoRoot, "deploy", "config", "zone-a")
	zoneGwDir := filepath.Join(repoRoot, "deploy", "config", "zone-a-gateway")
	require.NoError(t, os.MkdirAll(zoneADir, 0755))
	require.NoError(t, os.MkdirAll(zoneGwDir, 0755))

	writeTempYAML(t, zoneADir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "zone-a-workflow"
  instances:
    node-0: {}
    node-bt-0:
      nodeType: boot
`)
	writeTempYAML(t, zoneADir, "dev.yaml", "")
	writeTempYAML(t, zoneGwDir, "shared.yaml", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "zone-a-gateway"
  instances:
    gateway-0:
      nodeType: gateway
`)
	writeTempYAML(t, zoneGwDir, "dev.yaml", "")
	writeTempYAML(t, repoRoot, "griddle.yaml", `
deploy:
  dev:
    instances:
      - name: zone-a
        namespace: zone-a
        path: stack://node-set
        config:
          - deploy/config/zone-a/shared.yaml
          - deploy/config/zone-a/dev.yaml
      - name: zone-a-gateway
        namespace: zone-a-gateway
        path: stack://node-set
        config:
          - deploy/config/zone-a-gateway/shared.yaml
          - deploy/config/zone-a-gateway/dev.yaml
`)

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)

	require.Equal(t, "zone-a-workflow", cv.GetNode("node-0").DONName)
	require.Equal(t, "zone-a-workflow", cv.GetNode("node-bt-0").DONName)
	require.Equal(t, "zone-a-gateway", cv.GetNode("gateway-0").DONName)

	require.ElementsMatch(t, []string{"node-0", "node-bt-0"}, cv.NodeNamesForDONName("zone-a-workflow"))
	require.Equal(t, []string{"gateway-0"}, cv.NodeNamesForDONName("zone-a-gateway"))
	require.Empty(t, cv.NodeNamesForDONName("no-such-don"))
}

func TestLoadChartValues_DONNameLabelMissing(t *testing.T) {
	t.Parallel()

	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  instances:
    node-0: {}
`, "")

	_, err := LoadChartValues(repoRoot, "dev")
	require.Error(t, err)
	require.Contains(t, err.Error(), "don-name is required")
}

func TestLoadChartValues_BootstrapInferredByNaming(t *testing.T) {
	t.Parallel()

	// node-bt-0 doesn't have nodeType: boot in the chart, but should be
	// upgraded to boot by naming convention
	repoRoot := writeGriddleRepo(t, "dev", `
chainlink-node:
  registerNodes:
    labels:
      don-name: "my-repo"
  defaults:
    nodeType: standard
  instances:
    node-bt-0: {}
    node-0: {}
`, "")

	cv, err := LoadChartValues(repoRoot, "dev")
	require.NoError(t, err)

	bt := cv.GetNode("node-bt-0")
	require.NotNil(t, bt)
	require.Equal(t, RoleBootstrap, bt.NodeType, "node-bt-0 should be inferred as boot")

	standard := cv.GetNode("node-0")
	require.NotNil(t, standard)
	require.Equal(t, RoleStandard, standard.NodeType)

	require.Equal(t, "node-bt-0", cv.FindBootstrap())
}
