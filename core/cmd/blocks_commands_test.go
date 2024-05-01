package cmd_test

import (
	"flag"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func Test_ReplayFromBlock(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = (*ubig.Big)(big.NewInt(5))
		c.EVM[0].Enabled = ptr(true)
	})

	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.ReplayFromBlock, set, "")

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

func Test_FindLCA(t *testing.T) {
	t.Parallel()

	//ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = (*ubig.Big)(big.NewInt(5))
		c.EVM[0].Enabled = ptr(true)
	})

	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.FindLCA, set, "")

	//Incorrect chain ID
	require.NoError(t, set.Set("evm-chain-id", "1"))
	c := cli.NewContext(nil, set, nil)
	require.ErrorContains(t, client.FindLCA(c), "does not match any local chains")

	//Correct chain ID
	require.NoError(t, set.Set("evm-chain-id", "5"))
	c = cli.NewContext(nil, set, nil)
	require.ErrorContains(t, client.FindLCA(c), "FindLCA is only available if LogPoller is enabled")
}
