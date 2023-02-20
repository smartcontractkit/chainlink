package solana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	solanadb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/services/pg"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestSetupNodes(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	// Insert existing node which will be erased
	pgtest.MustExec(t, db, `INSERT INTO solana_chains (id, created_at, updated_at) VALUES ('test-setup',NOW(),NOW())`)
	pgtest.MustExec(t, db, `INSERT INTO solana_nodes (name, solana_chain_id, solana_url, created_at, updated_at) VALUES ('foo','test-setup','ws://example.com',NOW(),NOW())`)

	s := `
[
	{
		"name": "mainnet-one",
		"solanaChainId": "mainnet",
		"solanaURL": "ws://test1.invalid"
	},
	{
		"name": "mainnet-two",
		"solanaChainId": "mainnet",
		"solanaURL": "https://test2.invalid"
	},
	{
		"name": "testnet-one",
		"solanaChainId": "testnet",
		"solanaURL": "http://test3.invalid"
	},
	{
		"name": "testnet-two",
		"solanaChainId": "testnet",
		"solanaURL": "http://test4.invalid"
	}
]
	`

	cfg := config{
		solanaNodes: s,
		QConfig:     pgtest.NewQConfig(false),
	}

	err := solana.SetupNodes(db, cfg, logger.TestLogger(t))
	require.NoError(t, err)

	cltest.AssertCount(t, db, "solana_nodes", 4)

	var nodes []solanadb.Node
	err = db.Select(&nodes, `SELECT * FROM solana_nodes ORDER BY name ASC`)
	require.NoError(t, err)

	require.Len(t, nodes, 4)

	assert.Equal(t, "mainnet-one", nodes[0].Name)
	assert.Equal(t, "mainnet-two", nodes[1].Name)
	assert.Equal(t, "testnet-one", nodes[2].Name)
	assert.Equal(t, "testnet-two", nodes[3].Name)

}

type config struct {
	solanaNodes string
	pg.QConfig
}

func (c config) SolanaNodes() string {
	return c.solanaNodes
}
