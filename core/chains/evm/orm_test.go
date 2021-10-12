package evm_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func setupORM(t *testing.T) (*sqlx.DB, types.ORM) {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm := evm.NewORM(db)

	return db, orm
}

func mustInsertChain(t *testing.T, orm types.ORM) types.Chain {
	id := utils.NewBigI(99)
	config := types.ChainCfg{}
	chain, err := orm.CreateChain(*id, config)
	require.NoError(t, err)
	return chain
}

func Test_EVMORM_CreateChain(t *testing.T) {
	_, orm := setupORM(t)

	_, initialCount, err := orm.Chains(0, 25)
	require.NoError(t, err)

	id := utils.NewBigI(99)
	config := types.ChainCfg{}
	chain, err := orm.CreateChain(*id, config)
	require.NoError(t, err)
	require.Equal(t, chain.ID.ToInt().Int64(), id.ToInt().Int64())

	chains, count, err := orm.Chains(0, 25)
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
	require.Equal(t, chains[initialCount], chain)
}

func Test_EVMORM_CreateNode(t *testing.T) {
	_, orm := setupORM(t)
	chain := mustInsertChain(t, orm)

	_, initialCount, err := orm.Nodes(0, 25)
	require.NoError(t, err)

	params := types.NewNode{
		Name:       "Test node",
		EVMChainID: chain.ID,
		WSURL:      null.StringFrom("ws://localhost:8546"),
		HTTPURL:    null.StringFrom("http://localhost:8546"),
		SendOnly:   false,
	}

	t.Run("with successful callback", func(t *testing.T) {
		var called bool
		f := func(n types.Node) error {
			assert.Equal(t, "Test node", n.Name)
			called = true
			return nil
		}
		node, err := orm.CreateNode(context.Background(), params, f)
		require.NoError(t, err)
		assert.True(t, called)
		require.Equal(t, params.EVMChainID, node.EVMChainID)
		require.Equal(t, params.WSURL, node.WSURL)
		require.Equal(t, params.HTTPURL, node.HTTPURL)
		require.Equal(t, params.SendOnly, node.SendOnly)

		nodes, count, err := orm.Nodes(0, 25)
		require.NoError(t, err)
		require.Equal(t, initialCount+1, count)
		require.Equal(t, nodes[initialCount], node)
	})
}
