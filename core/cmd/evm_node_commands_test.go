package cmd_test

import (
	"bytes"
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func mustInsertEVMChain(t *testing.T, orm types.ORM) types.DBChain {
	id := utils.NewBig(testutils.NewRandomEVMChainID())
	chain, err := orm.CreateChain(*id, nil)
	require.NoError(t, err)
	return chain
}

func assertTableRenders(t *testing.T, r *cltest.RendererMock) {
	// Should be no error rendering any of the responses as tables
	b := bytes.NewBuffer([]byte{})
	tb := cmd.RendererTable{b}
	for _, rn := range r.Renders {
		require.NoError(t, tb.Render(rn))
	}
}

func TestClient_IndexEVMNodes(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	chain := mustInsertEVMChain(t, orm)

	params := types.Node{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    null.StringFrom("http://localhost:8546"),
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)

	require.Nil(t, cmd.NewEVMNodeClient(client).IndexNodes(cltest.EmptyCLIContext()))
	require.NotEmpty(t, r.Renders)
	nodes := *r.Renders[0].(*cmd.EVMNodePresenters)
	require.Len(t, nodes, initialCount+1)
	n := nodes[initialCount]
	assert.Equal(t, strconv.FormatInt(int64(node.ID), 10), n.ID)
	assert.Equal(t, params.Name, n.Name)
	assert.Equal(t, params.EVMChainID, n.EVMChainID)
	assert.Equal(t, params.WSURL, n.WSURL)
	assert.Equal(t, params.HTTPURL, n.HTTPURL)
	assertTableRenders(t, r)
}

func TestClient_CreateEVMNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	_, initialNodesCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)

	chain := mustInsertEVMChain(t, orm)

	// successful primary
	set := flag.NewFlagSet("cli", 0)
	set.String("name", "Example", "")
	set.String("type", "primary", "")
	set.String("ws-url", "ws://TestClient_CreateEVMNode1.invalid", "")
	set.String("http-url", "http://TestClient_CreateEVMNode2.invalid", "")
	set.Int64("chain-id", chain.ID.ToInt().Int64(), "")
	c := cli.NewContext(nil, set, nil)
	err = cmd.NewEVMNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	// successful send-only
	set = flag.NewFlagSet("cli", 0)
	set.String("name", "Send only", "")
	set.String("type", "sendonly", "")
	set.String("http-url", "http://TestClient_CreateEVMNode3.invalid", "")
	set.Int64("chain-id", chain.ID.ToInt().Int64(), "")
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewEVMNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	nodes, _, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, nodes, initialNodesCount+2)
	n := nodes[initialNodesCount]
	assert.Equal(t, "Example", n.Name)
	assert.Equal(t, false, n.SendOnly)
	assert.Equal(t, null.StringFrom("ws://TestClient_CreateEVMNode1.invalid"), n.WSURL)
	assert.Equal(t, null.StringFrom("http://TestClient_CreateEVMNode2.invalid"), n.HTTPURL)
	assert.Equal(t, chain.ID, n.EVMChainID)
	n = nodes[initialNodesCount+1]
	assert.Equal(t, "Send only", n.Name)
	assert.Equal(t, true, n.SendOnly)
	assert.Equal(t, null.String{}, n.WSURL)
	assert.Equal(t, null.StringFrom("http://TestClient_CreateEVMNode3.invalid"), n.HTTPURL)
	assert.Equal(t, chain.ID, n.EVMChainID)

	assertTableRenders(t, r)
}

func TestClient_RemoveEVMNode(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)

	chain := mustInsertEVMChain(t, orm)

	params := types.Node{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    null.StringFrom("http://localhost:8546"),
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)
	chains, _, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)})
	c := cli.NewContext(nil, set, nil)

	err = cmd.NewEVMNodeClient(client).RemoveNode(c)
	require.NoError(t, err)

	chains, _, err = orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}
