package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
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
