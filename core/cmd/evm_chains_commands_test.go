package cmd_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	client2 "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newRandChainID() *utils.Big {
	return utils.NewBig(testutils.NewRandomEVMChainID())
}

func TestClient_IndexEVMChains(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	})
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.EVMChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.EVMChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, strconv.Itoa(client2.NullClientChainID), c.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_CreateEVMChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t,
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.EVMChainClient(client).CreateChain, set, "")

	id := newRandChainID()
	require.NoError(t, set.Set("id", id.ToInt().String()))
	require.NoError(t, set.Parse([]string{`{}`}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.EVMChainClient(client).CreateChain(c)
	require.NoError(t, err)

	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, id.ToInt().Int64(), ch.ID.ToInt().Int64())
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_RemoveEVMChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t,
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()
	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	id := newRandChainID()
	_, err = orm.CreateChain(*id, nil)
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.EVMChainClient(client).RemoveChain, set, "")

	require.NoError(t, set.Parse([]string{id.String()}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.EVMChainClient(client).RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_ConfigureEVMChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t,
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()

	orm := app.EVMORM()

	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	id := newRandChainID()
	_, err = orm.CreateChain(*id, &types.ChainCfg{
		BlockHistoryEstimatorBlockDelay: null.IntFrom(5),
		EvmFinalityDepth:                null.IntFrom(5),
		EvmGasBumpPercent:               null.IntFrom(3),
	})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.EVMChainClient(client).ConfigureChain, set, "")

	require.NoError(t, set.Set("id", id.ToInt().String()))
	require.NoError(t, set.Parse([]string{"BlockHistoryEstimatorBlockDelay=9", "EvmGasBumpPercent=null"}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.EVMChainClient(client).ConfigureChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, null.IntFrom(int64(9)), ch.Cfg.BlockHistoryEstimatorBlockDelay) // this key was changed
	assert.Equal(t, null.IntFrom(int64(5)), ch.Cfg.EvmFinalityDepth)                // this key was unchanged
	assert.Equal(t, null.Int{}, ch.Cfg.EvmGasBumpPercent)                           // this key was unset
	assertTableRenders(t, r)
}
