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
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"

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

func terraStartNewApplication(t *testing.T, cfgs ...*terra.TerraConfig) *cltest.TestApplication {
	for i := range cfgs {
		cfgs[i].SetDefaults()
	}
	return startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Terra = cfgs
		c.EVM = nil
	})
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func terraStartNewLegacyApplication(t *testing.T) *cltest.TestApplication {
	return startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.TerraEnabled = null.BoolFrom(true)
		c.Overrides.EVMEnabled = null.BoolFrom(false)
		c.Overrides.EVMRPCEnabled = null.BoolFrom(false)
	}))
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
		Enabled: ptr(true),
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

func TestClient_CreateTerraNode(t *testing.T) {
	t.Parallel()

	app := terraStartNewLegacyApplication(t)
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
	cltest.CopyFlagSetFromAction(cmd.NewTerraNodeClient(client).CreateNode, set, "terra")

	require.NoError(t, set.Set("name", "first"))
	require.NoError(t, set.Set("tendermint-url", "http://tender.mint.test/columbus-5"))
	require.NoError(t, set.Set("chain-id", chainIDA))

	c := cli.NewContext(nil, set, nil)
	err = cmd.NewTerraNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	set = flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.NewTerraNodeClient(client).CreateNode, set, "terra")

	require.NoError(t, set.Set("name", "second"))
	require.NoError(t, set.Set("tendermint-url", "http://tender.mint.test/bombay-12"))
	require.NoError(t, set.Set("chain-id", chainIDB))

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

	app := terraStartNewLegacyApplication(t)
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
	cltest.CopyFlagSetFromAction(cmd.NewTerraNodeClient(client).RemoveNode, set, "terra")

	require.NoError(t, set.Parse([]string{strconv.FormatInt(int64(node.ID), 10)}))

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
