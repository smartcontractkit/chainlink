package cosmos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cosmosdb "github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink/core/services/pg"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestSetupNodes(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	// Insert existing node which will be erased
	pgtest.MustExec(t, db, `INSERT INTO cosmos_chains (id, created_at, updated_at) VALUES ('test-setup',NOW(),NOW())`)
	pgtest.MustExec(t, db, `INSERT INTO cosmos_nodes (name, cosmos_chain_id, tendermint_url, created_at, updated_at) VALUES ('foo','test-setup','ws://example.com',NOW(),NOW())`)

	s := `
[
	{
		"name": "bombay-one",
		"cosmosChainId": "bombay",
		"tendermintURL": "ws://test1.invalid"
	},
	{
		"name": "bombay-two",
		"cosmosChainId": "bombay",
		"tendermintURL": "https://test2.invalid"
	},
	{
		"name": "columbus-one",
		"cosmosChainId": "columbus",
		"tendermintURL": "http://test3.invalid"
	},
	{
		"name": "columbus-two",
		"cosmosChainId": "columbus",
		"tendermintURL": "http://test4.invalid"
	}
]
	`

	cfg := config{
		cosmosNodes: s,
		QConfig:     pgtest.NewQConfig(false),
	}

	err := cosmos.SetupNodes(db, cfg, logger.TestLogger(t))
	require.NoError(t, err)

	cltest.AssertCount(t, db, "cosmos_nodes", 4)

	var nodes []cosmosdb.Node
	err = db.Select(&nodes, `SELECT * FROM cosmos_nodes ORDER BY name ASC`)
	require.NoError(t, err)

	require.Len(t, nodes, 4)

	assert.Equal(t, "bombay-one", nodes[0].Name)
	assert.Equal(t, "bombay-two", nodes[1].Name)
	assert.Equal(t, "columbus-one", nodes[2].Name)
	assert.Equal(t, "columbus-two", nodes[3].Name)

}

type config struct {
	cosmosNodes string
	pg.QConfig
}

func (c config) CosmosNodes() string {
	return c.cosmosNodes
}
