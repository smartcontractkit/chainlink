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
families = ["family-a"]

[infra]
  type = "griddle"

[jd]
  grpc = "grpc-job-distributor.main.stage.cldev.sh:443"
  domain = "cre"
  environment = "dev"

[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true

[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "evm-1337"]
  nodes = ["node-0", "node-1"]
  family = "family-a"

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
families = ["family-a"]

[infra]
  type = "griddle"

[jd]
  grpc = "grpc-jd:443"
  domain = "cre"
  environment = "dev"

[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true

[[dons]]
  name = "workflow"
  don_types = ["workflow"]
  capabilities = ["cron", "http-action", "evm-1337"]
  nodes = ["node-0", "node-1", "node-2", "node-3"]
  family = "family-a"

[[dons]]
  name = "capabilities"
  don_types = ["capabilities"]
  exposes_remote_capabilities = true
  capabilities = ["vault"]
  nodes = ["node-0", "node-1"]
  family = "family-a"

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
	require.Equal(t, "family-a", ds.DONs[0].Family)
	require.ElementsMatch(t, []string{"family-a"}, ds.Families)
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
			name: "no DONs",
			toml: `[infra]
  type = "griddle"
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
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
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
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 1337
  family = "evm"
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
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 11155111
  family = "evm"
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
		{
			name: "missing family",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "family is required",
		},
		{
			name: "unsupported family",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "stellar"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "unsupported family",
		},
		{
			name: "non-evm chain cannot be registry",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 4
  family = "aptos"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "registry chain must be evm",
		},
		{
			name: "aptos capability references undeclared chain",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["aptos-4"]
  nodes = ["n"]
[capability_configs.aptos]
  binary_name = "aptos"
`,
			errSub: "not declared in [[chains]]",
		},
		{
			name: "aptos chain missing http_url",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 4
  family = "aptos"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "family aptos): http_url is required",
		},
		{
			name: "aptos chain cannot set ws_url",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 4
  family = "aptos"
  ws_url = "wss://aptos-4.example.com"
  http_url = "https://aptos-4.example.com"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "ws_url is not used by aptos",
		},
		{
			name: "solana chain missing genesis_hash",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 5
  family = "solana"
  ws_url = "wss://solana-5.example.com"
  http_url = "https://solana-5.example.com"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "chain_specific.genesis_hash is required",
		},
		{
			name: "solana chain missing ws_url",
			toml: `[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 5
  family = "solana"
  http_url = "https://solana-5.example.com"
  [chains.chain_specific]
    genesis_hash = "22222222222222222222222222222222222222222222"
[[dons]]
  name = "w"
  capabilities = ["cron"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "family solana): ws_url is required",
		},
		{
			name: "workflow DON missing family",
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
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "w"
  don_types = ["workflow"]
  capabilities = ["cron"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "family is required for workflow/capabilities/gateway DONs",
		},
		{
			name: "family not declared",
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
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "w"
  don_types = ["workflow"]
  capabilities = ["cron"]
  family = "family-a"
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "family \"family-a\" is not declared in [families]",
		},
		{
			name: "duplicate family name",
			toml: `families = ["family-a", "family-a"]
[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "w"
  capabilities = ["cron"]
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "duplicate family \"family-a\"",
		},
		{
			name: "gateway DON's family serves nothing",
			toml: `families = ["family-a"]
[infra]
  type = "griddle"
  chart_values = "x"
  namespace = "ns"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[dons]]
  name = "gw"
  don_types = ["gateway"]
  family = "family-a"
[capability_configs.cron]
  binary_name = "cron"
`,
			errSub: "has gateway DON(s)",
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

func TestLoadDesiredState_NonEVMChainDeclared(t *testing.T) {
	t.Parallel()

	path := writeTempTOML(t, `
[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 4
  family = "aptos"
  http_url = "https://aptos-4.example.com"
[[dons]]
  name = "w"
  capabilities = ["cron", "aptos-4"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
[capability_configs.aptos]
  binary_name = "aptos"
`)
	ds, err := LoadDesiredState(path)
	require.NoError(t, err)
	require.Len(t, ds.Chains, 2)
	require.Equal(t, "aptos", ds.Chains[1].Family)
	require.Equal(t, "https://aptos-4.example.com", ds.Chains[1].HTTPURL)
	require.ElementsMatch(t, []uint64{1337}, ds.EVMChainIDs())
}

func TestLoadDesiredState_SolanaChainDeclared(t *testing.T) {
	t.Parallel()

	path := writeTempTOML(t, `
[infra]
  type = "griddle"
[jd]
  grpc = "x"
  domain = "cre"
  environment = "dev"
[[chains]]
  chain_id = 1337
  family = "evm"
  ws_url = "wss://anvil-1337.example.com"
  http_url = "https://anvil-1337.example.com"
  registry = true
[[chains]]
  chain_id = 5
  family = "solana"
  ws_url = "wss://solana-5.example.com"
  http_url = "https://solana-5.example.com"
  [chains.chain_specific]
    genesis_hash = "22222222222222222222222222222222222222222222"
[[dons]]
  name = "w"
  capabilities = ["cron", "solana-5"]
  nodes = ["n"]
[capability_configs.cron]
  binary_name = "cron"
[capability_configs.solana]
  binary_name = "solana"
`)
	ds, err := LoadDesiredState(path)
	require.NoError(t, err)
	require.Len(t, ds.Chains, 2)
	require.Equal(t, "solana", ds.Chains[1].Family)
	require.Equal(t, "22222222222222222222222222222222222222222222", ds.Chains[1].SolanaGenesisHash())
}

func TestDON_NonEVMFamilies(t *testing.T) {
	t.Parallel()

	don := DON{Capabilities: []string{"cron", "aptos-4", "evm-1337"}}
	require.Equal(t, []string{"aptos"}, don.NonEVMFamilies())

	don2 := DON{Capabilities: []string{"cron", "solana", "aptos-4"}}
	require.ElementsMatch(t, []string{"solana", "aptos"}, don2.NonEVMFamilies())

	don3 := DON{Capabilities: []string{"cron", "evm-1337"}}
	require.Empty(t, don3.NonEVMFamilies())
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

func TestGatewayDONForNode(t *testing.T) {
	t.Parallel()

	ds := &DesiredState{
		DONs: []DON{
			{Name: "capabilities", DONTypes: []string{"capabilities"}},
			{Name: "gateway-don", DONTypes: []string{"gateway"}},
		},
	}
	cv := &ChartValues{
		Nodes: []ChartNodeInfo{
			{Name: "node-gw-0", DONName: "gateway-don"},
			{Name: "node-cap-0", DONName: "capabilities"},
		},
	}

	require.Equal(t, "gateway-don", ds.GatewayDONForNode(cv, "node-gw-0"))
	require.Empty(t, ds.GatewayDONForNode(cv, "node-cap-0"))
	require.Empty(t, ds.GatewayDONForNode(cv, "node-unknown"))
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
