package headtracker_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestORM_IdempotentInsertHead(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)

	// Returns nil when inserting first head
	head := testutils.Head(0)
	require.NoError(t, orm.IdempotentInsertHead(tests.Context(t), head))

	// Head is inserted
	foundHead, err := orm.LatestHead(tests.Context(t))
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)

	// Returns nil when inserting same head again
	require.NoError(t, orm.IdempotentInsertHead(tests.Context(t), head))

	// Head is still inserted
	foundHead, err = orm.LatestHead(tests.Context(t))
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)
}

func TestORM_TrimOldHeads(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)

	for i := 0; i < 10; i++ {
		head := testutils.Head(i)
		require.NoError(t, orm.IdempotentInsertHead(tests.Context(t), head))
	}

	uncleHead := testutils.Head(5)
	require.NoError(t, orm.IdempotentInsertHead(tests.Context(t), uncleHead))

	err := orm.TrimOldHeads(tests.Context(t), 5)
	require.NoError(t, err)

	heads, err := orm.LatestHeads(tests.Context(t), 0)
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
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)

	var hash common.Hash
	for i := 0; i < 10; i++ {
		head := testutils.Head(i)
		if i == 5 {
			hash = head.Hash
		}
		require.NoError(t, orm.IdempotentInsertHead(tests.Context(t), head))
	}

	head, err := orm.HeadByHash(tests.Context(t), hash)
	require.NoError(t, err)
	require.Equal(t, hash, head.Hash)
	require.Equal(t, int64(5), head.Number)
}

func TestORM_HeadByHash_NotFound(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)

	hash := testutils.Head(123).Hash
	head, err := orm.HeadByHash(tests.Context(t), hash)

	require.Nil(t, head)
	require.NoError(t, err)
}

func TestORM_LatestHeads_NoRows(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	orm := headtracker.NewORM(*testutils.FixtureChainID, db)

	heads, err := orm.LatestHeads(tests.Context(t), 100)

	require.Zero(t, len(heads))
	require.NoError(t, err)
}
