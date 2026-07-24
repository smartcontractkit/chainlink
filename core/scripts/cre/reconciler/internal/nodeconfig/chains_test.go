package nodeconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
)

func writeChartYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "dev.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestDiscoverChains_ReadsEVMConfigLayer(t *testing.T) {
	t.Parallel()

	path := writeChartYAML(t, `
chainlink-node:
  instances:
    node-0:
      configuration:
        - 10-network: |
            [[EVM]]
            ChainID = '11155111'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://sepolia.example.com/ws'
              HTTPURL = 'https://sepolia.example.com'
`)

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", ConfigFile: path},
		},
	}

	chains, err := DiscoverChains(cv)
	require.NoError(t, err)
	require.Len(t, chains, 1)
	require.Equal(t, uint64(11155111), chains[0].ChainID)
	require.Equal(t, "wss://sepolia.example.com/ws", chains[0].WSURL)
	require.Equal(t, "https://sepolia.example.com", chains[0].HTTPURL)
	require.False(t, chains[0].Registry, "discovery never sets the registry flag — the user designates it")
}

func TestDiscoverChains_SkipsManagedLayer(t *testing.T) {
	t.Parallel()

	path := writeChartYAML(t, `
chainlink-node:
  instances:
    node-0:
      configuration:
        - 30-cre: |
            [[EVM]]
            ChainID = '1337'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://should-not-appear.example.com'
              HTTPURL = 'https://should-not-appear.example.com'
`)

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", ConfigFile: path},
		},
	}

	chains, err := DiscoverChains(cv)
	require.NoError(t, err)
	require.Empty(t, chains, "the reconciler's own 30-cre layer must never be treated as chart-declared chain config")
}

func TestDiscoverChains_DedupesAcrossNodes(t *testing.T) {
	t.Parallel()

	path := writeChartYAML(t, `
chainlink-node:
  instances:
    node-0:
      configuration:
        - 10-network: |
            [[EVM]]
            ChainID = '1337'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://anvil-1337.example.com'
              HTTPURL = 'https://anvil-1337.example.com'
    node-1:
      configuration:
        - 10-network: |
            [[EVM]]
            ChainID = '1337'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://anvil-1337.example.com'
              HTTPURL = 'https://anvil-1337.example.com'
`)

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", ConfigFile: path},
			{Name: "node-1", ConfigFile: path},
		},
	}

	chains, err := DiscoverChains(cv)
	require.NoError(t, err)
	require.Len(t, chains, 1)
}

func TestDiscoverChains_KeepsDistinctURLVariantsForSameChainID(t *testing.T) {
	t.Parallel()

	// A gateway node can be configured with a different (sometimes unusable)
	// RPC for the same chain ID than worker nodes. Both variants must survive
	// discovery so the user can see and prune the wrong one in the UI, instead
	// of one silently overwriting the other.
	path := writeChartYAML(t, `
chainlink-node:
  instances:
    node-0:
      configuration:
        - 10-network: |
            [[EVM]]
            ChainID = '1337'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://worker-rpc.example.com'
              HTTPURL = 'https://worker-rpc.example.com'
    node-gw-0:
      configuration:
        - 10-network: |
            [[EVM]]
            ChainID = '1337'
              [[EVM.Nodes]]
              Name = 'primary'
              WSURL = 'wss://gateway-rpc.example.com'
              HTTPURL = 'https://gateway-rpc.example.com'
`)

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", ConfigFile: path},
			{Name: "node-gw-0", ConfigFile: path},
		},
	}

	chains, err := DiscoverChains(cv)
	require.NoError(t, err)
	require.Len(t, chains, 2)
	for _, ch := range chains {
		require.Equal(t, uint64(1337), ch.ChainID)
	}
	require.ElementsMatch(t, []string{"wss://worker-rpc.example.com", "wss://gateway-rpc.example.com"},
		[]string{chains[0].WSURL, chains[1].WSURL})
}

func TestDiscoverChains_NoConfiguration(t *testing.T) {
	t.Parallel()

	path := writeChartYAML(t, `
chainlink-node:
  instances:
    node-0: {}
`)

	cv := &domain.ChartValues{
		Nodes: []domain.ChartNodeInfo{
			{Name: "node-0", ConfigFile: path},
		},
	}

	chains, err := DiscoverChains(cv)
	require.NoError(t, err)
	require.Empty(t, chains)
}
