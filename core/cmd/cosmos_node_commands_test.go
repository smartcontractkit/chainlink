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

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

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
	require.Nil(t, cmd.NewCosmosNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.CosmosNodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, cltest.FormatWithPrefixedChainID(chainID, *node.Name), n.ID)
	assert.Equal(t, chainID, n.ChainID)
	assert.Equal(t, *node.Name, n.Name)
	wantConfig, err := toml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n.Config)
	assertTableRenders(t, r)

	//Render table and check the fields order
	b := new(bytes.Buffer)
	rt := cmd.RendererTable{b}
	require.NoError(t, nodes.RenderTable(rt))
	renderLines := strings.Split(b.String(), "\n")
	assert.Equal(t, 10, len(renderLines))
	assert.Contains(t, renderLines[2], "Name")
	assert.Contains(t, renderLines[2], n.Name)
	assert.Contains(t, renderLines[3], "Chain ID")
	assert.Contains(t, renderLines[3], n.ChainID)
	assert.Contains(t, renderLines[4], "State")
	assert.Contains(t, renderLines[4], n.State)
}
