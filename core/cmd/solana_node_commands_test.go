package cmd_test

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func mustInsertSolanaChain(t *testing.T, sol solana.ChainSet, id string) solana.DBChain {
	chain, err := sol.Add(testutils.Context(t), id, nil)
	require.NoError(t, err)
	return chain
}

func solanaStartNewApplication(t *testing.T, cfgs ...*solana.SolanaConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = cfgs
		c.EVM = nil
	})
}

// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func solanaStartNewLegacyApplication(t *testing.T) *cltest.TestApplication {
	return startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.SolanaEnabled = null.BoolFrom(true)
		c.Overrides.EVMEnabled = null.BoolFrom(false)
		c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	}))
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
		Enabled: ptr(true),
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

func TestClient_CreateSolanaNode(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	ctx := testutils.Context(t)
	_, initialNodesCount, err := sol.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainIDA := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	chainIDB := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	_ = mustInsertSolanaChain(t, sol, chainIDA)
	_ = mustInsertSolanaChain(t, sol, chainIDB)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewSolanaNodeClient(client).CreateNode, set, "solana")

	require.NoError(t, set.Set("name", "first"))
	require.NoError(t, set.Set("url", "http://tender.mint.test/columbus-5"))
	require.NoError(t, set.Set("chain-id", chainIDA))

	c := cli.NewContext(nil, set, nil)
	err = cmd.NewSolanaNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	set = flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewSolanaNodeClient(client).CreateNode, set, "solana")

	require.NoError(t, set.Set("name", "second"))
	require.NoError(t, set.Set("url", "http://tender.mint.test/bombay-12"))
	require.NoError(t, set.Set("chain-id", chainIDB))
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewSolanaNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	nodes, _, err := sol.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialNodesCount+2)
	n := nodes[initialNodesCount]
	assertEqualNodesSolana(t, db.Node{
		Name:          "first",
		SolanaChainID: chainIDA,
		SolanaURL:     "http://tender.mint.test/columbus-5",
	}, n)
	n = nodes[initialNodesCount+1]
	assertEqualNodesSolana(t, db.Node{
		Name:          "second",
		SolanaChainID: chainIDB,
		SolanaURL:     "http://tender.mint.test/bombay-12",
	}, n)

	assertTableRenders(t, r)
}

func TestClient_RemoveSolanaNode(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	ctx := testutils.Context(t)
	_, initialCount, err := sol.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainID := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	_ = mustInsertSolanaChain(t, sol, chainID)

	params := db.Node{
		Name:          "first",
		SolanaChainID: chainID,
		SolanaURL:     "http://tender.mint.test/columbus-5",
	}
	node, err := sol.CreateNode(ctx, params)
	require.NoError(t, err)
	nodes, _, err := sol.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewSolanaNodeClient(client).RemoveNode, set, "solana")

	require.NoError(t, set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)}))
	c := cli.NewContext(nil, set, nil)

	err = cmd.NewSolanaNodeClient(client).RemoveNode(c)
	require.NoError(t, err)

	nodes, _, err = sol.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialCount)
	assertTableRenders(t, r)
}

func assertEqualNodesSolana(t *testing.T, newNode db.Node, gotNode db.Node) {
	t.Helper()

	assert.Equal(t, newNode.Name, gotNode.Name)
	assert.Equal(t, newNode.SolanaChainID, gotNode.SolanaChainID)
	assert.Equal(t, newNode.SolanaURL, gotNode.SolanaURL)
}
