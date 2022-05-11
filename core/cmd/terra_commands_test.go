package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
)

func TestClient_TerraInit(t *testing.T) {
	t.Parallel()

	app := terraStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	newNode := types.NewNode{
		Name:          "first",
		TerraChainID:  "Columbus-5",
		TendermintURL: "http://tender.mint.test/columbus-5",
	}
	set := flag.NewFlagSet("cli", 0)
	set.String("name", newNode.Name, "")
	set.String("tendermint-url", newNode.TendermintURL, "")
	set.String("chain-id", newNode.TerraChainID, "")

	// Try to add node
	c := cli.NewContext(nil, set, nil)
	err := cmd.NewTerraNodeClient(client).CreateNode(c)
	require.Error(t, err)

	// Chain first
	setCh := flag.NewFlagSet("cli", 0)
	setCh.String("id", newNode.TerraChainID, "")
	setCh.Parse([]string{`{}`})
	cCh := cli.NewContext(nil, setCh, nil)
	err = cmd.TerraChainClient(client).CreateChain(cCh)
	require.NoError(t, err)

	// Then node
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewTerraNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	assertTableRenders(t, r)
}
