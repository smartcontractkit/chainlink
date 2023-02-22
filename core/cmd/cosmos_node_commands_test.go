package cmd_test

import (
	"flag"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	coscfg "github.com/smartcontractkit/chainlink-terra/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/cosmostest"
)

func mustInsertCosmosChain(t *testing.T, ter cosmos.ChainSet, id string) types.DBChain {
	chain, err := ter.Add(testutils.Context(t), id, nil)
	require.NoError(t, err)
	return chain
}

func cosmosStartNewApplication(t *testing.T, cfgs ...*cosmos.CosmosConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Cosmos = cfgs
		c.EVM = nil
	})
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func cosmosStartNewLegacyApplication(t *testing.T) *cltest.TestApplication {
	return startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.CosmosEnabled = null.BoolFrom(true)
		c.Overrides.EVMEnabled = null.BoolFrom(false)
		c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	}))
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

func TestClient_CreateCosmosNode(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Cosmos
	ctx := testutils.Context(t)
	_, initialNodesCount, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainIDA := cosmostest.RandomChainID()
	chainIDB := cosmostest.RandomChainID()
	_ = mustInsertCosmosChain(t, ter, chainIDA)
	_ = mustInsertCosmosChain(t, ter, chainIDB)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewCosmosNodeClient(client).CreateNode, set, "cosmos")

	require.NoError(t, set.Set("name", "first"))
	require.NoError(t, set.Set("tendermint-url", "http://tender.mint.test/columbus-5"))
	require.NoError(t, set.Set("chain-id", chainIDA))

	c := cli.NewContext(nil, set, nil)
	err = cmd.NewCosmosNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	set = flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewCosmosNodeClient(client).CreateNode, set, "cosmos")

	require.NoError(t, set.Set("name", "second"))
	require.NoError(t, set.Set("tendermint-url", "http://tender.mint.test/bombay-12"))
	require.NoError(t, set.Set("chain-id", chainIDB))

	c = cli.NewContext(nil, set, nil)
	err = cmd.NewCosmosNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	nodes, _, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialNodesCount+2)
	n := nodes[initialNodesCount]
	assertEqualNodesCosmos(t, types.NewNode{
		Name:          "first",
		CosmosChainID: chainIDA,
		TendermintURL: "http://tender.mint.test/columbus-5",
	}, n)
	n = nodes[initialNodesCount+1]
	assertEqualNodesCosmos(t, types.NewNode{
		Name:          "second",
		CosmosChainID: chainIDB,
		TendermintURL: "http://tender.mint.test/bombay-12",
	}, n)

	assertTableRenders(t, r)
}

func TestClient_RemoveCosmosNode(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Cosmos
	ctx := testutils.Context(t)
	_, initialCount, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainID := cosmostest.RandomChainID()
	_ = mustInsertCosmosChain(t, ter, chainID)

	params := db.Node{
		Name:          "first",
		CosmosChainID: chainID,
		TendermintURL: "http://tender.mint.test/columbus-5",
	}
	node, err := ter.CreateNode(ctx, params)
	require.NoError(t, err)
	chains, _, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewCosmosNodeClient(client).RemoveNode, set, "cosmos")

	require.NoError(t, set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.NewCosmosNodeClient(client).RemoveNode(c)
	require.NoError(t, err)

	chains, _, err = ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

func assertEqualNodesCosmos(t *testing.T, newNode types.NewNode, gotNode db.Node) {
	t.Helper()

	assert.Equal(t, newNode.Name, gotNode.Name)
	assert.Equal(t, newNode.CosmosChainID, gotNode.CosmosChainID)
	assert.Equal(t, newNode.TendermintURL, gotNode.TendermintURL)
}
