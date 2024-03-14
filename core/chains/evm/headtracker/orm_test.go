package headtracker_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestORM_IdempotentInsertHead(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	// Returns nil when inserting first head
	head := cltest.Head(0)
	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), head))

	// Head is inserted
	foundHead, err := orm.LatestHead(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)

	// Returns nil when inserting same head again
	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), head))

	// Head is still inserted
	foundHead, err = orm.LatestHead(testutils.Context(t))
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)
}

func TestORM_TrimOldHeads(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	for i := 0; i < 10; i++ {
		head := cltest.Head(i)
		require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), head))
	}

	uncleHead := cltest.Head(5)
	require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), uncleHead))

	err := orm.TrimOldHeads(testutils.Context(t), 5)
	require.NoError(t, err)

	heads, err := orm.LatestHeads(testutils.Context(t), 0)
	require.NoError(t, err)

	// uncle block was loaded too
	require.Equal(t, 6, len(heads))
	for i := 0; i < 5; i++ {
		require.LessOrEqual(t, int64(5), heads[i].Number)
	}
}

func TestORM_HeadByHash(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	var hash common.Hash
	for i := 0; i < 10; i++ {
		head := cltest.Head(i)
		if i == 5 {
			hash = head.Hash
		}
		require.NoError(t, orm.IdempotentInsertHead(testutils.Context(t), head))
	}

	head, err := orm.HeadByHash(testutils.Context(t), hash)
	require.NoError(t, err)
	require.Equal(t, hash, head.Hash)
	require.Equal(t, int64(5), head.Number)
}

func TestORM_HeadByHash_NotFound(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	hash := cltest.Head(123).Hash
	head, err := orm.HeadByHash(testutils.Context(t), hash)

	require.Nil(t, head)
	require.NoError(t, err)
}

func TestORM_LatestHeads_NoRows(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(cltest.FixtureChainID, db)

	heads, err := orm.LatestHeads(testutils.Context(t), 100)

	require.Zero(t, len(heads))
	require.NoError(t, err)
}
