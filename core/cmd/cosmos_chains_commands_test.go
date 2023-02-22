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
	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/cosmostest"
)

func TestClient_IndexCosmosChains(t *testing.T) {
	t.Parallel()

	chainID := cosmostest.RandomChainID()
	chain := cosmos.CosmosConfig{
		ChainID: ptr(chainID),
		Enabled: ptr(true),
	}
	app := cosmosStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.CosmosChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.CosmosChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, chainID, c.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_CreateCosmosChain(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Cosmos
	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	cosmosChainID := cosmostest.RandomChainID()
	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.CosmosChainClient(client).CreateChain, set, "cosmos")

	require.NoError(t, set.Set("id", cosmosChainID))
	require.NoError(t, set.Parse([]string{`{}`}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.CosmosChainClient(client).CreateChain(c)
	require.NoError(t, err)

	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, cosmosChainID, ch.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_RemoveCosmosChain(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Cosmos
	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	cosmosChainID := cosmostest.RandomChainID()
	_, err = ter.Add(ctx, cosmosChainID, nil)
	require.NoError(t, err)
	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.CosmosChainClient(client).RemoveChain, set, "cosmos")

	require.NoError(t, set.Parse([]string{cosmosChainID}))
	c := cli.NewContext(nil, set, nil)

	err = cmd.CosmosChainClient(client).RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_ConfigureCosmosChain(t *testing.T) {
	t.Parallel()

	app := cosmosStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	ter := app.Chains.Cosmos

	_, initialCount, err := ter.Index(0, 25)
	require.NoError(t, err)

	cosmosChainID := cosmostest.RandomChainID()
	minute, err := utils.NewDuration(time.Minute)
	require.NoError(t, err)
	original := db.ChainCfg{
		FallbackGasPriceULuna: null.StringFrom("99.07"),
		GasLimitMultiplier:    null.FloatFrom(1.111),
		ConfirmPollPeriod:     &minute,
	}
	ctx := testutils.Context(t)
	_, err = ter.Add(ctx, cosmosChainID, &original)
	require.NoError(t, err)
	chains, _, err := ter.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	cltest.FlagSetApplyFromAction(cmd.CosmosChainClient(client).ConfigureChain, set, "cosmos")

	require.NoError(t, set.Set("id", cosmosChainID))
	require.NoError(t, set.Parse([]string{
		"BlocksUntilTxTimeout=7",
		"FallbackGasPriceULuna=\"9.999\"",
		"GasLimitMultiplier=1.55555",
	}))
	c := cli.NewContext(nil, set, nil)

	err = cmd.CosmosChainClient(client).ConfigureChain(c)
	require.NoError(t, err)

	chains, _, err = ter.Index(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, cosmosChainID, ch.ID)
	assert.Equal(t, null.IntFrom(7), ch.Cfg.BlocksUntilTxTimeout)
	assert.Equal(t, null.StringFrom("9.999"), ch.Cfg.FallbackGasPriceULuna)
	assert.Equal(t, null.FloatFrom(1.55555), ch.Cfg.GasLimitMultiplier)
	assert.Equal(t, original.ConfirmPollPeriod, ch.Cfg.ConfirmPollPeriod)
	assertTableRenders(t, r)
}
