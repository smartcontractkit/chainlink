package cmd_test

import (
	"flag"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_ReplayFromBlock(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = (*utils.Big)(big.NewInt(5))
		c.EVM[0].Enabled = ptr(true)
	})

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.ReplayFromBlock, set, "")

	//Incorrect block number
	require.NoError(t, set.Set("block-number", "0"))
	c := cli.NewContext(nil, set, nil)
	require.ErrorContains(t, client.ReplayFromBlock(c), "Must pass a positive value in")

	//Incorrect chain ID
	require.NoError(t, set.Set("block-number", "1"))
	require.NoError(t, set.Set("evm-chain-id", "1"))
	c = cli.NewContext(nil, set, nil)
	require.ErrorContains(t, client.ReplayFromBlock(c), "does not match any local chains")

	//Correct chain ID
	require.NoError(t, set.Set("evm-chain-id", "5"))
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ReplayFromBlock(c))
}
