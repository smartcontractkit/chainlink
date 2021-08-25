package cmd_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v4"
)

func mustInsertChain(t *testing.T, orm evm.ORM) types.Chain {
	id := utils.NewBigI(99)
	config := types.ChainCfg{}
	chain, err := orm.CreateChain(*id, config)
	require.NoError(t, err)
	return chain
}

func TestClient_IndexNodes(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	chain := mustInsertChain(t, orm)

	params := evm.NewNode{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    "http://localhost:8546",
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)

	require.Nil(t, client.IndexNodes(cltest.EmptyCLIContext()))
	nodes := *r.Renders[0].(*cmd.NodePresenters)
	require.Len(t, nodes, 1)
	n := nodes[0]
	assert.Equal(t, strconv.FormatInt(node.ID, 10), n.ID)
	assert.Equal(t, params.Name, n.Name)
	assert.Equal(t, params.EVMChainID, n.EVMChainID)
	assert.Equal(t, params.WSURL, n.WSURL)
	assert.Equal(t, params.HTTPURL, n.HTTPURL)
}

func TestClient_CreateNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	orm := app.EVMORM()
	chain := mustInsertChain(t, orm)

	// successful primary
	set := flag.NewFlagSet("cli", 0)
	set.String("name", "Example", "")
	set.String("type", "primary", "")
	set.String("ws-url", "ws://", "")
	set.String("http-url", "http://", "")
	set.Int64("chain-id", chain.ID.ToInt().Int64(), "")
	c := cli.NewContext(nil, set, nil)
	err := client.CreateNode(c)
	require.NoError(t, err)

	// successful send-only
	set = flag.NewFlagSet("cli", 0)
	set.String("name", "Send only", "")
	set.String("type", "sendonly", "")
	set.String("http-url", "http://", "")
	set.Int64("chain-id", chain.ID.ToInt().Int64(), "")
	c = cli.NewContext(nil, set, nil)
	err = client.CreateNode(c)
	require.NoError(t, err)

	nodes, _, err := orm.Nodes(0, 25)
	require.Len(t, nodes, 2)
	n := nodes[0]
	assert.Equal(t, "Example", n.Name)
	assert.Equal(t, false, n.SendOnly)
	assert.Equal(t, null.StringFrom("ws://"), n.WSURL)
	assert.Equal(t, "http://", n.HTTPURL)
	assert.Equal(t, chain.ID, n.EVMChainID)
	n = nodes[1]
	assert.Equal(t, "Send only", n.Name)
	assert.Equal(t, true, n.SendOnly)
	assert.Equal(t, null.String{}, n.WSURL)
	assert.Equal(t, "http://", n.HTTPURL)
	assert.Equal(t, chain.ID, n.EVMChainID)
}

func TestClient_RemoveNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	orm := app.EVMORM()
	chain := mustInsertChain(t, orm)

	params := evm.NewNode{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    "http://localhost:8546",
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)
	chains, _, err := orm.Nodes(0, 25)
	require.Len(t, chains, 1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{strconv.FormatInt(node.ID, 10)})
	c := cli.NewContext(nil, set, nil)

	err = client.RemoveNode(c)
	require.NoError(t, err)

	chains, _, err = orm.Nodes(0, 25)
	require.Len(t, chains, 0)
}
