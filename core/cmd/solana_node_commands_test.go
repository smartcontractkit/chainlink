package cmd_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

func solanaStartNewApplication(t *testing.T, cfgs ...*solana.SolanaConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = cfgs
		c.EVM = nil
	})
}

func TestClient_IndexSolanaNodes(t *testing.T) {
	t.Parallel()

	id := solanatest.RandomChainID()
	node := solcfg.Node{
		Name: ptr("second"),
		URL:  utils.MustParseURL("https://solana.example"),
	}
	chain := solana.SolanaConfig{
		ChainID: &id,
		Nodes:   solana.SolanaNodes{&node},
	}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewSolanaNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.SolanaNodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, "0", n.ID)
	assert.Equal(t, *node.Name, n.Name)
	assert.Equal(t, (*url.URL)(node.URL).String(), n.SolanaURL)
	assertTableRenders(t, r)
}
