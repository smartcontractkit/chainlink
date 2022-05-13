package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setupORM(t *testing.T) (*sqlx.DB, types.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := evm.NewORM(db, logger.TestLogger(t), pgtest.PGCfg{})

	return db, orm
}

func mustInsertChain(t *testing.T, orm types.ORM) types.DBChain {
	t.Helper()

	id := utils.NewBigI(99)
	chain, err := orm.CreateChain(*id, nil)
	require.NoError(t, err)
	return chain
}

func mustInsertNode(t *testing.T, orm types.ORM, chainID utils.Big) types.Node {
	t.Helper()

	params := types.Node{
		Name:       "Test node",
		EVMChainID: chainID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    null.StringFrom("http://localhost:8546"),
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)

	return node
}

func Test_EVMORM_CreateChain(t *testing.T) {
	_, orm := setupORM(t)

	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	id := utils.NewBigI(99)
	chain, err := orm.CreateChain(*id, nil)
	require.NoError(t, err)
	require.Equal(t, chain.ID.ToInt().Int64(), id.ToInt().Int64())

	chains, count, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
	require.Equal(t, chains[initialCount], chain)
}

func Test_EVMORM_GetChainsByIDs(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)

	chains, err := orm.GetChainsByIDs([]utils.Big{chain.ID})
	require.NoError(t, err)
	require.Len(t, chains, 1)

	actual := chains[0]
	require.Equal(t, chain.ID, actual.ID)
	require.Equal(t, chain.Enabled, actual.Enabled)
	require.Equal(t, chain.Cfg, actual.Cfg)
}

func Test_EVMORM_CreateNode(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)

	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)

	params := types.Node{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    null.StringFrom("http://localhost:8546"),
		SendOnly:   false,
	}
	node, err := orm.CreateNode(params)
	require.NoError(t, err)
	require.Equal(t, params.EVMChainID, node.EVMChainID)
	require.Equal(t, params.WSURL, node.WSURL)
	require.Equal(t, params.HTTPURL, node.HTTPURL)
	require.Equal(t, params.SendOnly, node.SendOnly)

	nodes, count, err := orm.Nodes(0, 25)
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
	require.Equal(t, nodes[initialCount], node)

	assert.NoError(t, orm.DeleteChain(chain.ID))
}

func Test_EVMORM_GetNodesByChainIDs(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)
	node := mustInsertNode(t, orm, chain.ID)

	nodes, err := orm.GetNodesByChainIDs([]utils.Big{chain.ID})
	require.NoError(t, err)
	require.Len(t, nodes, 1)

	actual := nodes[0]

	require.Equal(t, node, actual)
}

func Test_EVMORM_Node(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)
	node := mustInsertNode(t, orm, chain.ID)

	actual, err := orm.Node(node.ID)
	assert.NoError(t, err)

	require.Equal(t, node, actual)
}
