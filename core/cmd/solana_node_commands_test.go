package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func solanaStartNewApplication(t *testing.T, cfgs ...*solana.TOMLConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
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
	chain := solana.TOMLConfig{
		ChainID: &id,
		Nodes:   solana.SolanaNodes{&node1, &node2},
	}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()

	require.Nil(t, cmd.NewSolanaNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.SolanaNodePresenters)
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

	//Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Equal(t, 17, len(renderLines))
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
