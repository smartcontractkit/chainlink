package cmd_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

func terraStartNewApplication(t *testing.T, cfgs ...*terra.TerraConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Terra = cfgs
		c.EVM = nil
	})
}

func TestClient_IndexTerraNodes(t *testing.T) {
	t.Parallel()

	chainID := terratest.RandomChainID()
	node := tercfg.Node{
		Name:          ptr("second"),
		TendermintURL: utils.MustParseURL("http://tender.mint.test/bombay-12"),
	}
	chain := terra.TerraConfig{
		ChainID: ptr(chainID),
		Nodes:   terra.TerraNodes{&node},
	}
	app := terraStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewTerraNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.TerraNodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, "0", n.ID)
	assert.Equal(t, *node.Name, n.Name)
	assert.Equal(t, *chain.ChainID, n.TerraChainID)
	assert.Equal(t, (*url.URL)(node.TendermintURL).String(), n.TendermintURL)
	assertTableRenders(t, r)
}
