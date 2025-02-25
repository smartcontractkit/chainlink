//go:build integration

package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/cosmostest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestShell_IndexCosmosChains(t *testing.T) {
	t.Parallel()

	chainID := cosmostest.RandomChainID()
	chain := chainlink.RawConfig{
		"ChainID": chainID,
		"Nodes": []map[string]any{{
			"Name":          "primary",
			"TendermintURL": "http://tender.mint",
		}},
	}
	app := cosmosStartNewApplication(t, chain)
	client, r := app.NewShellAndRenderer()

	require.NoError(t, cmd.NewChainClient(client, "cosmos").IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.ChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, chainID, c.ID)
	assertTableRenders(t, r)
}
