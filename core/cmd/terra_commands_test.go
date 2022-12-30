package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestClient_TerraInit(t *testing.T) {
	t.Parallel()

	app := terraStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	newNode := types.NewNode{
		Name:          "first",
		TerraChainID:  "Columbus-5",
		TendermintURL: "http://tender.mint.test/columbus-5",
	}
	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.NewTerraNodeClient(client).CreateNode, set, "terra")

	require.NoError(t, set.Set("name", newNode.Name))
	require.NoError(t, set.Set("tendermint-url", newNode.TendermintURL))
	require.NoError(t, set.Set("chain-id", newNode.TerraChainID))

	// Try to add node
	c := cli.NewContext(nil, set, nil)
	err := cmd.NewTerraNodeClient(client).CreateNode(c)
	require.Error(t, err)

	// Chain first
	setCh := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.TerraChainClient(client).CreateChain, setCh, "terra")

	require.NoError(t, setCh.Set("id", newNode.TerraChainID))
	require.NoError(t, setCh.Parse([]string{`{}`}))

	cCh := cli.NewContext(nil, setCh, nil)
	err = cmd.TerraChainClient(client).CreateChain(cCh)
	require.NoError(t, err)

	// Then node
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewTerraNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	assertTableRenders(t, r)
}
