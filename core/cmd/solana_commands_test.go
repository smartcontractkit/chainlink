package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestClient_SolanaInit(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	newNode := db.NewNode{
		Name:          "first",
		SolanaChainID: "Columbus-5",
		SolanaURL:     "https://solana.example",
	}
	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.NewSolanaNodeClient(client).CreateNode, set, "solana")

	require.NoError(t, set.Set("name", newNode.Name))
	require.NoError(t, set.Set("url", newNode.SolanaURL))
	require.NoError(t, set.Set("chain-id", newNode.SolanaChainID))

	// Try to add node
	c := cli.NewContext(nil, set, nil)
	err := cmd.NewSolanaNodeClient(client).CreateNode(c)
	require.Error(t, err)

	// Chain first
	setCh := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.SolanaChainClient(client).CreateChain, setCh, "Solana")

	require.NoError(t, setCh.Set("id", newNode.SolanaChainID))
	require.NoError(t, setCh.Parse([]string{`{}`}))

	cCh := cli.NewContext(nil, setCh, nil)
	err = cmd.SolanaChainClient(client).CreateChain(cCh)
	require.NoError(t, err)

	// Then node
	c = cli.NewContext(nil, set, nil)
	err = cmd.NewSolanaNodeClient(client).CreateNode(c)
	require.NoError(t, err)

	assertTableRenders(t, r)
}
