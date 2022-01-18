package cmd_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func mustInsertTerraChain(t *testing.T, orm types.ORM, id string) types.Chain {
	chain, err := orm.CreateChain(id, types.ChainCfg{})
	require.NoError(t, err)
	return chain
}

func TestClient_IndexTerraNodes(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	_ = mustInsertTerraChain(t, orm, "Bombay-12")

	params := types.NewNode{
		Name:          "second",
		TerraChainID:  "Bombay-12",
		TendermintURL: "http://tender.mint.test/bombay-12",
		FCDURL:        "http://fcd.test/bombay-12",
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)

	require.Nil(t, client.IndexTerraNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.TerraNodePresenters)
	require.Len(t, nodes, initialCount+1)
	n := nodes[initialCount]
	assert.Equal(t, strconv.FormatInt(int64(node.ID), 10), n.ID)
	assert.Equal(t, params.Name, n.Name)
	assert.Equal(t, params.TerraChainID, n.TerraChainID)
	assert.Equal(t, params.TendermintURL, n.TendermintURL)
	assert.Equal(t, params.FCDURL, n.FCDURL)
	assertTableRenders(t, r)
}

func TestClient_CreateTerraNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialNodesCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	_ = mustInsertTerraChain(t, orm, "Columbus-5")
	_ = mustInsertTerraChain(t, orm, "Bombay-12")

	set := flag.NewFlagSet("cli", 0)
	set.String("name", "first", "")
	set.String("tendermint-url", "http://tender.mint.test/columbus-5", "")
	set.String("fcd-url", "http://fcd.test/columbus-5", "")
	set.String("chain-id", "Columbus-5", "")
	c := cli.NewContext(nil, set, nil)
	err = client.CreateTerraNode(c)
	require.NoError(t, err)

	set = flag.NewFlagSet("cli", 0)
	set.String("name", "second", "")
	set.String("tendermint-url", "http://tender.mint.test/bombay-12", "")
	set.String("fcd-url", "http://fcd.test/bombay-12", "")
	set.String("chain-id", "Bombay-12", "")
	c = cli.NewContext(nil, set, nil)
	err = client.CreateTerraNode(c)
	require.NoError(t, err)

	nodes, _, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialNodesCount+2)
	n := nodes[initialNodesCount]
	assertEqual(t, types.NewNode{
		Name:          "first",
		TerraChainID:  "Columbus-5",
		TendermintURL: "http://tender.mint.test/columbus-5",
		FCDURL:        "http://fcd.test/columbus-5",
	}, n)
	n = nodes[initialNodesCount+1]
	assertEqual(t, types.NewNode{
		Name:          "second",
		TerraChainID:  "Bombay-12",
		TendermintURL: "http://tender.mint.test/bombay-12",
		FCDURL:        "http://fcd.test/bombay-12",
	}, n)

	assertTableRenders(t, r)
}

func TestClient_RemoveTerraNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	_ = mustInsertTerraChain(t, orm, "Columbus-5")

	params := types.NewNode{
		Name:          "first",
		TerraChainID:  "Columbus-5",
		TendermintURL: "http://tender.mint.test/columbus-5",
		FCDURL:        "http://fcd.test/columbus-5",
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)
	chains, _, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)})
	c := cli.NewContext(nil, set, nil)

	err = client.RemoveTerraNode(c)
	require.NoError(t, err)

	chains, _, err = orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

func assertEqual(t *testing.T, newNode types.NewNode, gotNode types.Node) {
	t.Helper()

	assert.Equal(t, newNode.Name, gotNode.Name)
	assert.Equal(t, newNode.TerraChainID, gotNode.TerraChainID)
	assert.Equal(t, newNode.TendermintURL, gotNode.TendermintURL)
	assert.Equal(t, newNode.FCDURL, gotNode.FCDURL)
}
