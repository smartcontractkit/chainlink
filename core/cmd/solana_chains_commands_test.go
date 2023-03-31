package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/solanatest"
)

func TestClient_IndexSolanaChains(t *testing.T) {
	t.Parallel()

	id := solanatest.RandomChainID()
	chain := solana.SolanaConfig{ChainID: &id}
	app := solanaStartNewApplication(t, &chain)
	client, r := app.NewClientAndRenderer()

	require.Nil(t, cmd.SolanaChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.SolanaChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, id, c.ID)
	assertTableRenders(t, r)
}
