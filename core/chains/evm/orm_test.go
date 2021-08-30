package evm_test

import (
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setupORM(t *testing.T) (*sqlx.DB, evm.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := evm.NewORM(db)

	return db, orm
}

func mustInsertChain(t *testing.T, orm evm.ORM) types.Chain {
	id := utils.NewBigI(99)
	config := types.ChainCfg{}
	chain, err := orm.CreateChain(*id, config)
	require.NoError(t, err)
	return chain
}

func Test_EVMORM_CreateChain(t *testing.T) {
	_, orm := setupORM(t)

	id := utils.NewBigI(99)
	config := types.ChainCfg{}
	chain, err := orm.CreateChain(*id, config)
	require.NoError(t, err)
	require.Equal(t, chain.ID.ToInt().Int64(), id.ToInt().Int64())

	chains, count, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Equal(t, 2, count) // it includes the default Ethereum chain already
	require.Equal(t, chains[1], chain)
}

func Test_EVMORM_CreateNode(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)

	params := evm.NewNode{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    "http://localhost:8546",
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
	require.Equal(t, 1, count)
	require.Equal(t, nodes[0], node)
}
