package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
)

func TestClient_SolanaInit(t *testing.T) {
	t.Parallel()

	app := solanaStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	newNode := db.NewNode{
		Name:          "first",
		SolanaChainID: "Columbus-5",
		SolanaURL:     "TODO",
	}
	set := flag.NewFlagSet("cli", 0)
	set.String("name", newNode.Name, "")
	set.String("url", newNode.SolanaURL, "")
	set.String("chain-id", newNode.SolanaChainID, "")

	// Try to add node
	c := cli.NewContext(nil, set, nil)
	err := client.CreateSolanaNode(c)
	require.Error(t, err)

	// Chain first
	setCh := flag.NewFlagSet("cli", 0)
	setCh.String("id", newNode.SolanaChainID, "")
	setCh.Parse([]string{`{}`})
	cCh := cli.NewContext(nil, setCh, nil)
	err = client.CreateSolanaChain(cCh)
	require.NoError(t, err)

	// Then node
	c = cli.NewContext(nil, set, nil)
	err = client.CreateSolanaNode(c)
	require.NoError(t, err)

	assertTableRenders(t, r)
}
