package cmd_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

const terraChainID = "Chainlinktest-99"

func TestClient_IndexTerraChains(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	chain, err := orm.CreateChain(terraChainID, db.ChainCfg{})
	require.NoError(t, err)

	require.Nil(t, client.IndexTerraChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.TerraChainPresenters)
	require.Len(t, chains, initialCount+1)
	c := chains[initialCount]
	assert.Equal(t, chain.ID, c.ID)
	assertTableRenders(t, r)
}

func TestClient_CreateTerraChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	set := flag.NewFlagSet("cli", 0)
	set.String("id", terraChainID, "")
	set.Parse([]string{`{}`})
	c := cli.NewContext(nil, set, nil)

	err = client.CreateTerraChain(c)
	require.NoError(t, err)

	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, terraChainID, ch.ID)
	assertTableRenders(t, r)
}

func TestClient_RemoveTerraChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()
	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	_, err = orm.CreateChain(terraChainID, db.ChainCfg{})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{terraChainID})
	c := cli.NewContext(nil, set, nil)

	err = client.RemoveTerraChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

func TestClient_ConfigureTerraChain(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	orm := app.TerraORM()

	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	_, err = orm.CreateChain(terraChainID, db.ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("99.07"),
		GasLimitMultiplier:    null.FloatFrom(1.111),
	})
	require.NoError(t, err)
	chains, _, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.String("id", terraChainID, "param")
	set.Parse([]string{"FallbackGasPriceULuna=\"9.999\"", "GasLimitMultiplier=1.55555"}) //TODO more
	c := cli.NewContext(nil, set, nil)

	err = client.ConfigureTerraChain(c)
	require.NoError(t, err)

	chains, _, err = orm.Chains(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, terraChainID, ch.ID)
	assert.Equal(t, null.StringFrom("9.999"), ch.Cfg.FallbackGasPriceULuna) // this key was changed
	assert.Equal(t, null.FloatFrom(1.55555), ch.Cfg.GasLimitMultiplier)     // this key was unchanged
	assertTableRenders(t, r)
}
