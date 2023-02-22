package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestClient_CosmosInit(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	newNode := types.NewNode{
		Name:          "first",
		CosmosChainID: "Columbus-5",
		TendermintURL: "http://tender.mint.test/columbus-5",
	}
	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewCosmosNodeClient(client).CreateNode, set, "cosmos")

	require.NoError(t, set.Set("name", newNode.Name))
	require.NoError(t, set.Set("tendermint-url", newNode.TendermintURL))
	require.NoError(t, set.Set("chain-id", newNode.CosmosChainID))

	// Try to add node
	c := cli.NewContext(nil, set, nil)
	err := cmd.NewCosmosNodeClient(client).CreateNode(c)
	require.Error(t, err)

	// Chain first
	setCh := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.CosmosChainClient(client).CreateChain, setCh, "cosmos")

	require.NoError(t, setCh.Set("id", newNode.CosmosChainID))
	require.NoError(t, setCh.Parse([]string{`{}`}))

	cCh := cli.NewContext(nil, setCh, nil)
	err = cmd.CosmosChainClient(client).CreateChain(cCh)
	require.NoError(t, err)

	// Then node
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewCosmosNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	assertTableRenders(t, r)
}
