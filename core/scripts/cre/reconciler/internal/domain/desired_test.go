package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempTOML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "desired.toml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestLoadDesiredState_Minimal(t *testing.T) {
	t.Parallel()

	path := writeTempTOML(t, `
[infra]
  type = "griddle"
  chart_values = "deploy/config/my-repo"
  namespace = "my-repo-nodeset"

[jd]
  grpc = "grpc-job-distributor.main.stage.cldev.sh:443"
  domain = "cre"
  environment = "dev"

[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true

[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "evm-1337"]
  nodes = ["node-0", "node-1"]

[capability_configs.cron]
  binary_name = "cron"

[capability_configs.evm]
  binary_name = "evm"
[capability_configs.evm.values]
  LogTriggerPollInterval = 1500000000
`)
	ds, err := LoadDesiredState(path)
	require.NoError(t, err)
	require.Equal(t, "griddle", ds.Infra.Type)
	require.Equal(t, "deploy/config/my-repo", ds.Infra.ChartValues)
	require.Len(t, ds.DONs, 1)
	require.Equal(t, "workflow", ds.DONs[0].Name)
	require.ElementsMatch(t, []string{"cron", "evm-1337"}, ds.DONs[0].Capabilities)
	// The legacy "nodes" key is present in the TOML above but the schema no
	// longer has a field for it — go-toml silently ignores unknown keys, and
	// membership is derived from the chart instead (see chartvalues_test.go).
}

func TestLoadDesiredState_MultipleDONs(t *testing.T) {
	t.Parallel()

	path := writeTempTOML(t, `
[infra]
  type = "griddle"
  chart_values = "deploy/config/my-repo"
  namespace = "my-repo-nodeset"

[jd]
  grpc = "grpc-jd:443"
  domain = "cre"
  environment = "dev"

[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true

[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "http-action", "evm-1337"]
  nodes = ["node-0", "node-1", "node-2", "node-3"]

[[dons]]
  name = "capabilities"
  don_types = ["capabilities"]
  exposes_remote_capabilities = true
  capabilities = ["vault"]
  nodes = ["node-0", "node-1"]

[[gateway_nodes]]
  node = "node-gw-0"
  don = "workflow"

[capability_configs.cron]
  binary_name = "cron"
[capability_configs.http-action]
  binary_name = "http_action"
[capability_configs.evm]
  binary_name = "evm"
[capability_configs.vault]
  binary_name = "vault"
`)
	ds, err := LoadDesiredState(path)
	require.NoError(t, err)
	require.Len(t, ds.DONs, 2)
	require.Equal(t, "capabilities", ds.DONs[1].Name)
	require.True(t, ds.DONs[1].ExposesRemoteCaps)
	require.Len(t, ds.GatewayNodes, 1)
	require.Equal(t, "node-gw-0", ds.GatewayNodes[0].Node)
}

func TestLoadDesiredState_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		toml   string
		errSub string
	}{
		{
			name: "wrong infra type",
			toml: `[infra]
  type = "docker"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "must be \"griddle\"",
		},
		{
			name: "missing chart_values",
			toml: `[infra]
  type = "griddle"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "chart_values is required",
		},
		{
			name: "no DONs",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
`,
			errSub: "at least one",
		},
		{
			name: "duplicate DON name",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[dons]]
  name = "dup"
  capabilities = ["cron"]
  nodes = ["n"]
[[dons]]
  name = "dup"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "duplicate DON name",
		},
		{
			name: "missing capability config",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[dons]]
  name = "w"
  capabilities = ["nonexistent"]
  nodes = ["n"]
`,
			errSub: "no capability_configs",
		},
		{
			name: "no chains declared",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "at least one [[chains]]",
		},
		{
			name: "no registry chain",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "exactly one [[chains]] entry must set registry",
		},
		{
			name: "duplicate chain_id",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337-b.example.com"
  http_url = "https://anvil-1337-b.example.com"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "duplicate chain_id",
		},
		{
			name: "capability references undeclared chain",
			toml: `[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 11155111
  ws_url = "wss://sepolia.example.com"
  http_url = "https://sepolia.example.com"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["evm-1337"]
  nodes = ["n"]
[capability_configs.evm]
  binary_name = "evm"
`,
			errSub: "not declared in [[chains]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := writeTempTOML(t, tt.toml)
			_, err := LoadDesiredState(path)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errSub)
		})
	}
}

func TestNeedsGateway(t *testing.T) {
	t.Parallel()

	ds := &DesiredState{
		DONs: []DON{
			{Capabilities: []string{"cron", "evm-1337"}},
		},
	}
	require.False(t, ds.NeedsGateway())

	ds.DONs[0].Capabilities = []string{"cron", "vault"}
	require.True(t, ds.NeedsGateway())

	ds.DONs[0].Capabilities = []string{"http-action"}
	require.True(t, ds.NeedsGateway())
}

func TestGatewayDONFor(t *testing.T) {
	t.Parallel()

	ds := &DesiredState{
		DONs: []DON{
			{Name: "capabilities", DONTypes: []string{"capabilities"}},
			{Name: "workflow", DONTypes: []string{"workflow"}},
		},
	}
	// No explicit assignment -> first workflow DON
	require.Equal(t, "workflow", ds.GatewayDONFor("node-gw-0"))

	// Explicit assignment
	ds.GatewayNodes = []GatewayNodeAssignment{
		{Node: "node-gw-0", DON: "capabilities"},
	}
	require.Equal(t, "capabilities", ds.GatewayDONFor("node-gw-0"))
}

func TestStripChainSuffix(t *testing.T) {
	t.Parallel()

	require.Equal(t, "evm", stripChainSuffix("evm-1337"))
	require.Equal(t, "solana", stripChainSuffix("solana-123"))
	require.Equal(t, "cron", stripChainSuffix("cron"))
	require.Equal(t, "http-action", stripChainSuffix("http-action"))
}

func TestIsConfiglessCapability(t *testing.T) {
	t.Parallel()

	require.True(t, isConfiglessCapability("don-time"))
	require.True(t, isConfiglessCapability("consensus"))
	require.False(t, isConfiglessCapability("cron"))
	require.False(t, isConfiglessCapability("evm-1337"))
	require.False(t, isConfiglessCapability("vault"))
}

func TestDON_HasDONType(t *testing.T) {
	t.Parallel()

	don := DON{DONTypes: []string{"workflow", "capabilities"}}
	require.True(t, don.HasDONType("workflow"))
	require.True(t, don.HasDONType("capabilities"))
	require.False(t, don.HasDONType("gateway"))
}

func TestDON_IsWorkflowDon(t *testing.T) {
	t.Parallel()

	require.True(t, (&DON{DONTypes: []string{"workflow"}}).IsWorkflowDon())
	require.False(t, (&DON{DONTypes: []string{"capabilities"}}).IsWorkflowDon())
}

func TestDON_IsGatewayDon(t *testing.T) {
	t.Parallel()

	require.True(t, (&DON{DONTypes: []string{"gateway"}}).IsGatewayDon())
	require.False(t, (&DON{DONTypes: []string{"workflow"}}).IsGatewayDon())
}

func TestDON_IsBootstrapDon(t *testing.T) {
	t.Parallel()

	require.True(t, (&DON{DONTypes: []string{"bootstrap"}}).IsBootstrapDon())
	require.False(t, (&DON{DONTypes: []string{"gateway"}}).IsBootstrapDon())
}

func TestDON_NeedsGatewayAccess(t *testing.T) {
	t.Parallel()

	require.True(t, (&DON{Capabilities: []string{"vault"}}).NeedsGatewayAccess())
	require.True(t, (&DON{Capabilities: []string{"http-action"}}).NeedsGatewayAccess())
	require.False(t, (&DON{Capabilities: []string{"cron"}}).NeedsGatewayAccess())
}

func TestLoadDesiredState_GatewayOnlyDONRegression(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "..", "tmp", "desired.toml")
	if _, err := os.Stat(path); err != nil {
		t.Skip("tmp/desired.toml not present")
	}

	ds, err := LoadDesiredState(path)
	require.NoError(t, err)

	cv := &ChartValues{
		Nodes: []ChartNodeInfo{
			{Name: "gateway-0", NodeType: RoleGateway, Namespace: "zone-a-gateway"},
		},
	}

	var gatewayDon *DON
	for i := range ds.DONs {
		if ds.DONs[i].Name == "don-3" {
			gatewayDon = &ds.DONs[i]
			break
		}
	}
	require.NotNil(t, gatewayDon)
	require.True(t, gatewayDon.IsGatewayDon())
	require.False(t, gatewayDon.IsBootstrapOnly(cv))
}

// DON membership is always chart-derived via cv.NodeNamesForDONName(don.Name),
// so these tests set ChartNodeInfo.DONName to match the DON under test instead
// of populating a (now-removed) DON.Nodes field.
func TestDON_ResolveBootstrap(t *testing.T) {
	t.Parallel()

	cv := &ChartValues{
		Nodes: []ChartNodeInfo{
			{Name: "node-bt-0", NodeType: RoleBootstrap, DONName: "workflow"},
			{Name: "node-0", NodeType: RoleStandard, DONName: "workflow"},
			{Name: "node-1", NodeType: RoleStandard, DONName: "workflow"},
		},
	}

	t.Run("explicit bootstrap_node", func(t *testing.T) {
		t.Parallel()

		don := DON{Name: "workflow", BootstrapNode: "node-bt-0"}
		require.Equal(t, "node-bt-0", don.ResolveBootstrap(cv))
	})

	t.Run("inferred from nodeType", func(t *testing.T) {
		t.Parallel()

		don := DON{Name: "workflow"}
		require.Equal(t, "node-bt-0", don.ResolveBootstrap(cv))
	})

	t.Run("inferred from naming convention", func(t *testing.T) {
		t.Parallel()

		cvNoType := &ChartValues{
			Nodes: []ChartNodeInfo{
				{Name: "node-bt-0", NodeType: RoleStandard, DONName: "workflow"}, // no nodeType set
				{Name: "node-0", NodeType: RoleStandard, DONName: "workflow"},
			},
		}
		don := DON{Name: "workflow"}
		require.Equal(t, "node-bt-0", don.ResolveBootstrap(cvNoType))
	})

	t.Run("no bootstrap found", func(t *testing.T) {
		t.Parallel()

		cvNoBootstrap := &ChartValues{
			Nodes: []ChartNodeInfo{
				{Name: "node-0", NodeType: RoleStandard, DONName: "no-boot"},
				{Name: "node-1", NodeType: RoleStandard, DONName: "no-boot"},
			},
		}
		don := DON{Name: "no-boot"}
		require.Empty(t, don.ResolveBootstrap(cvNoBootstrap))
	})
}

func TestDON_WorkerNodes(t *testing.T) {
	t.Parallel()

	cv := &ChartValues{
		Nodes: []ChartNodeInfo{
			{Name: "node-bt-0", NodeType: RoleBootstrap, DONName: "workflow"},
			{Name: "node-0", NodeType: RoleStandard, DONName: "workflow"},
			{Name: "node-1", NodeType: RoleStandard, DONName: "workflow"},
			{Name: "node-gw-0", NodeType: RoleGateway, DONName: "workflow"},
		},
	}

	t.Run("excludes bootstrap and gateway", func(t *testing.T) {
		t.Parallel()

		don := DON{Name: "workflow"}
		workers := don.WorkerNodes(cv)
		require.ElementsMatch(t, []string{"node-0", "node-1"}, workers)
	})

	t.Run("bootstrap-only DON", func(t *testing.T) {
		t.Parallel()

		cvBootstrapOnly := &ChartValues{
			Nodes: []ChartNodeInfo{{Name: "node-bt-0", NodeType: RoleBootstrap, DONName: "bootstrap-don"}},
		}
		don := DON{Name: "bootstrap-don", DONTypes: []string{"bootstrap"}}
		workers := don.WorkerNodes(cvBootstrapOnly)
		require.Empty(t, workers)
		require.True(t, don.IsBootstrapOnly(cvBootstrapOnly))
		require.True(t, don.IsBootstrapDon())
	})

	t.Run("gateway-only DON", func(t *testing.T) {
		t.Parallel()

		cvWithGateway := &ChartValues{
			Nodes: []ChartNodeInfo{{Name: "gateway-0", NodeType: RoleGateway, DONName: "don-3"}},
		}
		don := DON{Name: "don-3", DONTypes: []string{"gateway"}}
		workers := don.WorkerNodes(cvWithGateway)
		require.Empty(t, workers)
		require.False(t, don.IsBootstrapOnly(cvWithGateway))
		require.True(t, don.IsGatewayDon())
	})

	t.Run("workflow DON with workers", func(t *testing.T) {
		t.Parallel()

		cvNoBootstrap := &ChartValues{
			Nodes: []ChartNodeInfo{
				{Name: "node-0", NodeType: RoleStandard, DONName: "workflow-only"},
				{Name: "node-1", NodeType: RoleStandard, DONName: "workflow-only"},
			},
		}
		don := DON{Name: "workflow-only", DONTypes: []string{"workflow"}}
		workers := don.WorkerNodes(cvNoBootstrap)
		require.ElementsMatch(t, []string{"node-0", "node-1"}, workers)
		require.False(t, don.IsBootstrapOnly(cvNoBootstrap))
	})

	t.Run("mixed DON", func(t *testing.T) {
		t.Parallel()

		cvMixed := &ChartValues{
			Nodes: []ChartNodeInfo{
				{Name: "node-bt-0", NodeType: RoleBootstrap, DONName: "mixed"},
				{Name: "node-0", NodeType: RoleStandard, DONName: "mixed"},
			},
		}
		don := DON{Name: "mixed"}
		workers := don.WorkerNodes(cvMixed)
		require.ElementsMatch(t, []string{"node-0"}, workers)
		require.False(t, don.IsBootstrapOnly(cvMixed))
	})
}
