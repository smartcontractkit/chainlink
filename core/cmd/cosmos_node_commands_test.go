package cmd_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/cosmostest"
)

func cosmosStartNewApplication(t *testing.T, cfgs ...*cosmos.CosmosConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cfgs
		c.EVM = nil
	})
}

func TestClient_IndexCosmosNodes(t *testing.T) {
	t.Parallel()

	chainID := cosmostest.RandomChainID()
	node := coscfg.Node{
		Name:          ptr("second"),
		TendermintURL: utils.MustParseURL("http://tender.mint.test/bombay-12"),
	}
	chain := cosmos.CosmosConfig{
		ChainID: ptr(chainID),
		Enabled: ptr(true),
		Nodes:   cosmos.CosmosNodes{&node},
	}
	app := cosmosStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewCosmosNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.CosmosNodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, "0", n.ID)
	assert.Equal(t, *node.Name, n.Name)
	assert.Equal(t, *chain.ChainID, n.CosmosChainID)
	assert.Equal(t, (*url.URL)(node.TendermintURL).String(), n.TendermintURL)
	assertTableRenders(t, r)
}
