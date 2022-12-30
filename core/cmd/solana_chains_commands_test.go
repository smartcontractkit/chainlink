package cmd_test

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains/solana"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/solanatest"
)

func TestClient_IndexSolanaChains(t *testing.T) {
	t.Parallel()

	id := solanatest.RandomChainID()
	chain := solana.SolanaConfig{
		ChainID: &id,
		Enabled: ptr(true),
	}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.SolanaChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.SolanaChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, id, c.ID)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_CreateSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana
	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	solanaChainID := solanatest.RandomChainID()
	set := flag.NewFlagSet("cli", 0)
	cltest.CopyFlagSetFromAction(cmd.SolanaChainClient(client).CreateChain, set, "Solana")

	require.NoError(t, set.Set("id", solanaChainID))
	require.NoError(t, set.Parse([]string{`{}`}))

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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_RemoveSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
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
	cltest.CopyFlagSetFromAction(cmd.SolanaChainClient(client).RemoveChain, set, "")

	require.NoError(t, set.Parse([]string{solanaChainID}))

	c := cli.NewContext(nil, set, nil)

	err = cmd.SolanaChainClient(client).RemoveChain(c)
	require.NoError(t, err)

	chains, _, err = sol.Index(0, 25)
	require.NoError(t, err)
	require.Len(t, chains, initialCount)
	assertTableRenders(t, r)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_ConfigureSolanaChain(t *testing.T) {
	t.Parallel()

	app := solanaStartNewLegacyApplication(t)
	client, r := app.NewClientAndRenderer()

	sol := app.Chains.Solana

	_, initialCount, err := sol.Index(0, 25)
	require.NoError(t, err)

	solanaChainID := solanatest.RandomChainID()
	minute, err := utils.NewDuration(time.Minute)
	require.NoError(t, err)
	hour, err := utils.NewDuration(time.Hour)
	require.NoError(t, err)
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
	cltest.CopyFlagSetFromAction(cmd.SolanaChainClient(client).ConfigureChain, set, "Solana")

	require.NoError(t, set.Set("id", solanaChainID))
	require.NoError(t, set.Parse([]string{"TxTimeout=1h"}))

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
