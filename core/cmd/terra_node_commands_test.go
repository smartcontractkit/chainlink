package cmd_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
)

func mustInsertTerraChain(t *testing.T, ter terra.ChainSet, id string) types.DBChain {
	chain, err := ter.Add(testutils.Context(t), id, nil)
	require.NoError(t, err)
	return chain
}

func terraStartNewApplication(t *testing.T) *cltest.TestApplication {
	return startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.TerraEnabled = null.BoolFrom(true)
		c.Overrides.EVMEnabled = null.BoolFrom(false)
		c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	}))
}

func TestClient_IndexTerraNodes(t *testing.T) {
	t.Parallel()

	app := terraStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra
	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)
	chainID := terratest.RandomChainID()
	_ = mustInsertTerraChain(t, ter, chainID)

	params := db.Node{
		Name:          "second",
		TerraChainID:  chainID,
		TendermintURL: "http://tender.mint.test/bombay-12",
	}
	ctx := testutils.Context(t)
	node, err := ter.CreateNode(ctx, params)
	require.NoError(t, err)

	require.Nil(t, cmd.NewTerraNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.TerraNodePresenters)
	require.Len(t, nodes, initialCount+1)
	n := nodes[initialCount]
	assert.Equal(t, strconv.FormatInt(int64(node.ID), 10), n.ID)
	assert.Equal(t, params.Name, n.Name)
	assert.Equal(t, params.TerraChainID, n.TerraChainID)
	assert.Equal(t, params.TendermintURL, n.TendermintURL)
	assertTableRenders(t, r)
}

func TestClient_CreateTerraNode(t *testing.T) {
	t.Parallel()

	app := terraStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra
	ctx := testutils.Context(t)
	_, initialNodesCount, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainIDA := terratest.RandomChainID()
	chainIDB := terratest.RandomChainID()
	_ = mustInsertTerraChain(t, ter, chainIDA)
	_ = mustInsertTerraChain(t, ter, chainIDB)

	set := flag.NewFlagSet("cli", 0)
	set.String("name", "first", "")
	set.String("tendermint-url", "http://tender.mint.test/columbus-5", "")
	set.String("fcd-url", "http://fcd.test/columbus-5", "")
	set.String("chain-id", chainIDA, "")
	c := cli.NewContext(nil, set, nil)
	err = cmd.NewTerraNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	set = flag.NewFlagSet("cli", 0)
	set.String("name", "second", "")
	set.String("tendermint-url", "http://tender.mint.test/bombay-12", "")
	set.String("fcd-url", "http://fcd.test/bombay-12", "")
	set.String("chain-id", chainIDB, "")
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewTerraNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	nodes, _, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialNodesCount+2)
	n := nodes[initialNodesCount]
	assertEqualNodesTerra(t, types.NewNode{
		Name:          "first",
		TerraChainID:  chainIDA,
		TendermintURL: "http://tender.mint.test/columbus-5",
	}, n)
	n = nodes[initialNodesCount+1]
	assertEqualNodesTerra(t, types.NewNode{
		Name:          "second",
		TerraChainID:  chainIDB,
		TendermintURL: "http://tender.mint.test/bombay-12",
	}, n)

	assertTableRenders(t, r)
}

func TestClient_RemoveTerraNode(t *testing.T) {
	t.Parallel()

	app := terraStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra
	ctx := testutils.Context(t)
	_, initialCount, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	chainID := terratest.RandomChainID()
	_ = mustInsertTerraChain(t, ter, chainID)

	params := db.Node{
		Name:          "first",
		TerraChainID:  chainID,
		TendermintURL: "http://tender.mint.test/columbus-5",
	}
	node, err := ter.CreateNode(ctx, params)
	require.NoError(t, err)
	chains, _, err := ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)})
	c := cli.NewContext(nil, set, nil)

	err = cmd.NewTerraNodeClient(client).RemoveNode(c)
	require.NoError(t, err)

	chains, _, err = ter.GetNodes(ctx, 0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

func assertEqualNodesTerra(t *testing.T, newNode types.NewNode, gotNode db.Node) {
	t.Helper()

	assert.Equal(t, newNode.Name, gotNode.Name)
	assert.Equal(t, newNode.TerraChainID, gotNode.TerraChainID)
	assert.Equal(t, newNode.TendermintURL, gotNode.TendermintURL)
}
