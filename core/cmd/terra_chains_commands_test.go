package cmd_test

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/terratest"
)

func TestClient_IndexTerraChains(t *testing.T) {
	t.Parallel()

	chainID := terratest.RandomChainID()
	chain := terra.TerraConfig{
		ChainID: ptr(chainID),
		Enabled: ptr(true),
	}
	app := terraStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.TerraChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.TerraChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, chainID, c.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_CreateTerraChain(t *testing.T) {
	t.Parallel()

	app := terraStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra
	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	terraChainID := terratest.RandomChainID()
	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.TerraChainClient(client).CreateChain, set, "terra")

	require.NoError(t, set.Set("id", terraChainID))
	require.NoError(t, set.Parse([]string{`{}`}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.TerraChainClient(client).CreateChain(c)
	require.NoError(t, err)

	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, terraChainID, ch.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_RemoveTerraChain(t *testing.T) {
	t.Parallel()

	app := terraStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra
	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	terraChainID := terratest.RandomChainID()
	_, err = ter.Add(ctx, terraChainID, nil)
	require.NoError(t, err)
	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.TerraChainClient(client).RemoveChain, set, "terra")

	require.NoError(t, set.Parse([]string{terraChainID}))
	c := cli.NewContext(nil, set, nil)

	err = cmd.TerraChainClient(client).RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_ConfigureTerraChain(t *testing.T) {
	t.Parallel()

	app := terraStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Terra

	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	terraChainID := terratest.RandomChainID()
	minute, err := utils.NewDuration(time.Minute)
	require.NoError(t, err)
	original := db.ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("99.07"),
		GasLimitMultiplier:    null.FloatFrom(1.111),
		ConfirmPollPeriod:     &minute,
	}
	ctx := testutils.Context(t)
	_, err = ter.Add(ctx, terraChainID, &original)
	require.NoError(t, err)
	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.TerraChainClient(client).ConfigureChain, set, "terra")

	require.NoError(t, set.Set("id", terraChainID))
	require.NoError(t, set.Parse([]string{
		"BlocksUntilTxTimeout=7",
		"FallbackGasPriceULuna=\"9.999\"",
		"GasLimitMultiplier=1.55555",
	}))
	c := cli.NewContext(nil, set, nil)

	err = cmd.TerraChainClient(client).ConfigureChain(c)
	require.NoError(t, err)

	chains, _, err = ter.Index(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, terraChainID, ch.ID)
	assert.Equal(t, null.IntFrom(7), ch.Cfg.BlocksUntilTxTimeout)
	assert.Equal(t, null.StringFrom("9.999"), ch.Cfg.FallbackGasPriceULuna)
	assert.Equal(t, null.FloatFrom(1.55555), ch.Cfg.GasLimitMultiplier)
	assert.Equal(t, original.ConfirmPollPeriod, ch.Cfg.ConfirmPollPeriod)
	assertTableRenders(t, r)
}

func ptr[T any](t T) *T { return &t }
