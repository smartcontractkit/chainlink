//go:build integration

package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func Test_ReplayFromBlock_Solana(t *testing.T) {
	t.Parallel()

	chain := chainlink.RawConfig{
		"ChainID": "devnet",
		"Enabled": true,
		"Nodes": []map[string]any{{
			"Name": "primary",
			"URL":  "http://solana.example",
		}},
	}
	app := solanaStartNewApplication(t, chain)
	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.ReplayFromBlock, set, "")

	require.NoError(t, set.Set("block-number", "1"))
	require.NoError(t, set.Set("chain-id", "devnet"))
	require.NoError(t, set.Set("family", "solana"))
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ReplayFromBlock(c))
}
