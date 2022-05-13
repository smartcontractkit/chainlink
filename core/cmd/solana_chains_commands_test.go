package cmd_test

import (
	"flag"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestClient_IndexSolanaChains(t *testing.T) {
	t.Parallel()

	app := solanaStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	chain, err := sol.Add(ctx, solanatest.RandomChainID(), nil)
	require.NoError(t, err)

	require.Nil(t, cmd.SolanaChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.SolanaChainPresenters)
	require.Len(t, chains, initialCount+1)
	c := chains[initialCount]
	assert.Equal(t, chain.ID, c.ID)
	assertTableRenders(t, r)
}

func TestClient_CreateSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	solanaChainID := solanatest.RandomChainID()
	set := flag.NewFlagSet("cli", 0)
	set.String("id", solanaChainID, "")
	set.Parse([]string{`{}`})
	c := cli.NewContext(nil, set, nil)

	err = cmd.SolanaChainClient(client).CreateChain(c)
	require.NoError(t, err)

	chains, _, err := sol.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)
	ch := chains[initialCount]
	assert.Equal(t, solanaChainID, ch.ID)
	assertTableRenders(t, r)
}

func TestClient_RemoveSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	solanaChainID := solanatest.RandomChainID()
	_, err = sol.Add(ctx, solanaChainID, nil)
	require.NoError(t, err)
	chains, _, err := sol.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.Parse([]string{solanaChainID})
	c := cli.NewContext(nil, set, nil)

	err = cmd.SolanaChainClient(client).RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = sol.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

func TestClient_ConfigureSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana

	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	solanaChainID := solanatest.RandomChainID()
	minute := models.MustMakeDuration(time.Minute)
	hour := models.MustMakeDuration(time.Hour)
	original := db.ChainCfg{
		ConfirmPollPeriod: &minute,
		TxTimeout:         &hour,
	}
	ctx := testutils.Context(t)
	_, err = sol.Add(ctx, solanaChainID, &original)
	require.NoError(t, err)
	chains, _, err := sol.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount+1)

	set := flag.NewFlagSet("cli", 0)
	set.String("id", solanaChainID, "param")
	set.Parse([]string{
		"TxTimeout=1h",
	})
	c := cli.NewContext(nil, set, nil)

	err = cmd.SolanaChainClient(client).ConfigureChain(c)
	require.NoError(t, err)

	chains, _, err = sol.Index(0, 25)
	require.NoError(t, err)
	ch := chains[initialCount]

	assert.Equal(t, solanaChainID, ch.ID)
	assert.Equal(t, original.ConfirmPollPeriod, ch.Cfg.ConfirmPollPeriod)
	assert.Equal(t, original.TxTimeout, ch.Cfg.TxTimeout)
	assertTableRenders(t, r)
}
