package cmd_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
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
	assert.Equal(t, chainID, n.ChainID)
	assert.Equal(t, *node.Name, n.ID)
	assert.Equal(t, *node.Name, n.Name)
	wantConfig, err := toml.Marshal(node)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n.Config)
	assertTableRenders(t, r)
}
