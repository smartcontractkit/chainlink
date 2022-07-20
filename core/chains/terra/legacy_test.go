package terra_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestSetupNodes(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	// Insert existing node which will be erased
	pgtest.MustExec(t, db, `INSERT INTO terra_chains (id, created_at, updated_at) VALUES ('test-setup',NOW(),NOW())`)
	pgtest.MustExec(t, db, `INSERT INTO terra_nodes (name, terra_chain_id, tendermint_url, created_at, updated_at) VALUES ('foo','test-setup','ws://example.com',NOW(),NOW())`)

	s := `
[
	{
		"name": "bombay-one",
		"terraChainId": "bombay",
		"tendermintURL": "ws://test1.invalid"
	},
	{
		"name": "bombay-two",
		"terraChainId": "bombay",
		"tendermintURL": "https://test2.invalid"
	},
	{
		"name": "columbus-one",
		"terraChainId": "columbus",
		"tendermintURL": "http://test3.invalid"
	},
	{
		"name": "columbus-two",
		"terraChainId": "columbus",
		"tendermintURL": "http://test4.invalid"
	}
]
	`

	cfg := config{
		terraNodes: s,
	}

	err := terra.SetupNodes(db, cfg, logger.TestLogger(t))
	require.NoError(t, err)

	cltest.AssertCount(t, db, "terra_nodes", 4)

	var nodes []terradb.Node
	err = db.Select(&nodes, `SELECT * FROM terra_nodes ORDER BY name ASC`)
	require.NoError(t, err)

	require.Len(t, nodes, 4)

	assert.Equal(t, "bombay-one", nodes[0].Name)
	assert.Equal(t, "bombay-two", nodes[1].Name)
	assert.Equal(t, "columbus-one", nodes[2].Name)
	assert.Equal(t, "columbus-two", nodes[3].Name)

}

type config struct {
	terraNodes string
}

func (c config) TerraNodes() string {
	return c.terraNodes
}

func (c config) LogSQL() bool { return false }
