package cmd_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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
	node1 := solcfg.Node{
		Name: ptr("first"),
		URL:  utils.MustParseURL("https://solana1.example"),
	}
	node2 := solcfg.Node{
		Name: ptr("second"),
		URL:  utils.MustParseURL("https://solana2.example"),
	}
	chain := solana.SolanaConfig{
		ChainID: &id,
		Nodes:   solana.SolanaNodes{&node1, &node2},
	}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.NewSolanaNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.SolanaNodePresenters)
	require.Len(t, nodes, 2)
	n1 := nodes[0]
	n2 := nodes[1]
	assert.Equal(t, "first", n1.ID)
	assert.Equal(t, *node1.Name, n1.Name)
	assert.Equal(t, (*url.URL)(node1.URL).String(), n1.SolanaURL)
	assert.Equal(t, "second", n2.ID)
	assert.Equal(t, *node2.Name, n2.Name)
	assert.Equal(t, (*url.URL)(node2.URL).String(), n2.SolanaURL)
	assertTableRenders(t, r)
}
