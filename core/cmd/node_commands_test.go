package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	starkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func assertTableRenders(t *testing.T, r *cltest.RendererMock) {
	// Should be no error rendering any of the responses as tables
	b := bytes.NewBuffer([]byte{})
	tb := cmd.RendererTable{b}
	for _, rn := range r.Renders {
		require.NoError(t, tb.Render(rn))
	}
}

func TestShell_IndexEVMNodes(t *testing.T) {
	t.Parallel()

	chainID := newRandChainID()
	node1 := evmcfg.Node{
		Name:     ptr("Test node 1"),
		WSURL:    config.MustParseURL("ws://localhost:8546"),
		HTTPURL:  config.MustParseURL("http://localhost:8546"),
		SendOnly: ptr(false),
		Order:    ptr(int32(15)),
	}
	node2 := evmcfg.Node{
		Name:     ptr("Test node 2"),
		WSURL:    config.MustParseURL("ws://localhost:8547"),
		HTTPURL:  config.MustParseURL("http://localhost:8547"),
		SendOnly: ptr(false),
		Order:    ptr(int32(36)),
	}
	chain := evmcfg.EVMConfig{
		ChainID: chainID,
		Chain:   evmcfg.Defaults(chainID),
		Nodes:   evmcfg.EVMNodes{&node1, &node2},
	}
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM = evmcfg.EVMConfigs{&chain}
	})
	client, r := app.NewShellAndRenderer()

	require.NoError(t, cmd.NewNodeClient(client, "evm").IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.NodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, chainID.String(), n1.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID.String(), *node1.Name), n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, chainID.String(), n2.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID.String(), *node2.Name), n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)

	// Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Len(t, renderLines, 23)
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n1.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n1.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n1.State)
	assert.Contains(t, renderLines[12], "Name")
	assert.Contains(t, renderLines[12], n2.Name)
	assert.Contains(t, renderLines[13], "Chain ID")
	assert.Contains(t, renderLines[13], n2.ChainID)
	assert.Contains(t, renderLines[14], "State")
	assert.Contains(t, renderLines[14], n2.State)
}

func cosmosStartNewApplication(t *testing.T, cfgs ...*coscfg.TOMLConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cfgs
		c.EVM = nil
	})
}

func TestShell_IndexCosmosNodes(t *testing.T) {
	t.Parallel()

	chainID := cosmostest.RandomChainID()
	node := coscfg.Node{
		Name:          ptr("second"),
		TendermintURL: config.MustParseURL("http://tender.mint.test/bombay-12"),
	}
	chain := coscfg.TOMLConfig{
		ChainID: ptr(chainID),
		Enabled: ptr(true),
		Nodes:   coscfg.Nodes{&node},
	}
	app := cosmosStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()
	require.NoError(t, cmd.NewNodeClient(client, "cosmos").IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.NodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID, *node.Name), n.ID)
	assert.Equal(t, chainID, n.ChainID)
	assert.Equal(t, *node.Name, n.Name)
	wantConfig, err := toml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n.Config)
	assertTableRenders(t, r)

	// Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Len(t, renderLines, 10)
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n.State)
}
func solanaStartNewApplication(t *testing.T, cfgs ...*solcfg.TOMLConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].Chain.SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = cfgs
		c.EVM = nil
	})
}

func TestShell_IndexSolanaNodes(t *testing.T) {
	t.Parallel()

	id := solanatest.RandomChainID()
	node1 := solcfg.Node{
		Name: ptr("first"),
		URL:  config.MustParseURL("https://solana1.example"),
	}
	node2 := solcfg.Node{
		Name: ptr("second"),
		URL:  config.MustParseURL("https://solana2.example"),
	}
	chain := solcfg.TOMLConfig{
		ChainID: &id,
		Nodes:   solcfg.Nodes{&node1, &node2},
	}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()

	require.NoError(t, cmd.NewNodeClient(client, "solana").IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.NodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, id, n1.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node1.Name), n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, id, n2.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node2.Name), n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)

	// Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Len(t, renderLines, 19)
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n1.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n1.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n1.State)
	assert.Contains(t, renderLines[10], "Name")
	assert.Contains(t, renderLines[10], n2.Name)
	assert.Contains(t, renderLines[11], "Chain ID")
	assert.Contains(t, renderLines[11], n2.ChainID)
	assert.Contains(t, renderLines[12], "State")
	assert.Contains(t, renderLines[12], n2.State)
}

func starknetStartNewApplication(t *testing.T, cfgs ...*starkcfg.TOMLConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = cfgs
		c.EVM = nil
		c.Solana = nil
	})
}

func TestShell_IndexStarkNetNodes(t *testing.T) {
	t.Parallel()

	id := "starknet chain ID"
	node1 := starkcfg.Node{
		Name: ptr("first"),
		URL:  config.MustParseURL("https://starknet1.example"),
	}
	node2 := starkcfg.Node{
		Name: ptr("second"),
		URL:  config.MustParseURL("https://starknet2.example"),
	}
	chain := starkcfg.TOMLConfig{
		ChainID: &id,
		Nodes:   starkcfg.Nodes{&node1, &node2},
	}
	app := starknetStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()

	require.NoError(t, cmd.NewNodeClient(client, "starknet").IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.NodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, id, n1.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node1.Name), n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, id, n2.ChainID)
	assert.Equal(t, cltest.FormatWithPrefixedChainID(id, *node2.Name), n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)

	// Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Len(t, renderLines, 17)
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n1.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n1.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n1.State)
	assert.Contains(t, renderLines[9], "Name")
	assert.Contains(t, renderLines[9], n2.Name)
	assert.Contains(t, renderLines[10], "Chain ID")
	assert.Contains(t, renderLines[10], n2.ChainID)
	assert.Contains(t, renderLines[11], "State")
	assert.Contains(t, renderLines[11], n2.State)
}
