package cmd_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func starknetStartNewApplication(t *testing.T, cfgs ...*starknet.StarknetConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = cfgs
		c.EVM = nil
		c.Solana = nil
	})
}

func TestClient_IndexStarkNetNodes(t *testing.T) {
	t.Parallel()

	id := "starknet chain ID"
	node1 := config.Node{
		Name: ptr("first"),
		URL:  utils.MustParseURL("https://starknet1.example"),
	}
	node2 := config.Node{
		Name: ptr("second"),
		URL:  utils.MustParseURL("https://starknet2.example"),
	}
	chain := starknet.StarknetConfig{
		ChainID: &id,
		Nodes:   starknet.StarknetNodes{&node1, &node2},
	}
	app := starknetStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewStarkNetNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.StarkNetNodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, id, n1.ChainID)
	assert.Equal(t, *node1.Name, n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	wantConfig, err := toml.Marshal(node1)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig), n1.Config)
	assert.Equal(t, id, n2.ChainID)
	assert.Equal(t, *node2.Name, n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	wantConfig2, err := toml.Marshal(node2)
	require.NoError(t, err)
	assert.Equal(t, string(wantConfig2), n2.Config)
	assertTableRenders(t, r)
}
