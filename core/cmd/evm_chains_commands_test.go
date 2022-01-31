package cmd_test

import (
	"flag"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v4"
)

func TestClient_IndexEVMChains(t *testing.T) {
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

	id := utils.NewBigI(99)
	chain, err := orm.CreateChain(*id, types.ChainCfg{})
	require.NoError(t, err)

	require.Nil(t, client.IndexEVMChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.EVMChainPresenters)
	require.Len(t, chains, initialCount+1)
	c := chains[initialCount]
	assert.Equal(t, chain.ID.ToInt().String(), c.ID)
	assertTableRenders(t, r)
}

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
	set.Int64("id", 99, "")
	set.Parse([]string{`{}`})
	c := cli.NewContext(nil, set, nil)

	err = client.CreateEVMChain(c)
	require.NoError(t, err)

	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, int64(99), ch.ID.ToInt().Int64())
	assertTableRenders(t, r)
}

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

	id := utils.NewBigI(99)
	_, err = orm.CreateChain(*id, types.ChainCfg{})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{"99"})
	c := cli.NewContext(nil, set, nil)

	err = client.RemoveEVMChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

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

	id := utils.NewBigI(99)
	_, err = orm.CreateChain(*id, types.ChainCfg{
		BlockHistoryEstimatorBlockDelay: null.IntFrom(5),
		EvmFinalityDepth:                null.IntFrom(5),
		EvmGasBumpPercent:               null.IntFrom(3),
	})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Int64("id", 99, "param")
	set.Parse([]string{"BlockHistoryEstimatorBlockDelay=9", "EvmGasBumpPercent=null"})
	c := cli.NewContext(nil, set, nil)

	err = client.ConfigureEVMChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, null.IntFrom(int64(9)), ch.Cfg.BlockHistoryEstimatorBlockDelay) // this key was changed
	assert.Equal(t, null.IntFrom(int64(5)), ch.Cfg.EvmFinalityDepth)                // this key was unchanged
	assert.Equal(t, null.Int{}, ch.Cfg.EvmGasBumpPercent)                           // this key was unset
	assertTableRenders(t, r)
}
