package headtracker_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_IdempotentInsertHead(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

	// Returns nil when inserting first head
	head := cltest.Head(0)
	require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))

	// Head is inserted
	foundHead, err := orm.LatestHead(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)

	// Returns nil when inserting same head again
	require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))

	// Head is still inserted
	foundHead, err = orm.LatestHead(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, head.Hash, foundHead.Hash)
}

func TestORM_TrimOldHeads(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

	for i := 0; i < 10; i++ {
		head := cltest.Head(i)
		require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))
	}

	err := orm.TrimOldHeads(context.TODO(), 5)
	require.NoError(t, err)

	heads, err := orm.LatestHeads(context.TODO(), 10)
	require.NoError(t, err)

	require.Equal(t, 5, len(heads))
	for i := 0; i < 5; i++ {
		require.LessOrEqual(t, int64(5), heads[i].Number)
	}
}

func TestORM_HeadByHash(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

	var hash common.Hash
	for i := 0; i < 10; i++ {
		head := cltest.Head(i)
		if i == 5 {
			hash = head.Hash
		}
		require.NoError(t, orm.IdempotentInsertHead(context.TODO(), head))
	}

	head, err := orm.HeadByHash(context.TODO(), hash)
	require.NoError(t, err)
	require.Equal(t, hash, head.Hash)
	require.Equal(t, int64(5), head.Number)
}

func TestORM_HeadByHash_NotFound(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

	hash := cltest.Head(123).Hash
	head, err := orm.HeadByHash(context.TODO(), hash)

	require.Nil(t, head)
	require.NoError(t, err)
}

func TestORM_LatestHeads_NoRows(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logger := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := headtracker.NewORM(db, logger, cfg, cltest.FixtureChainID)

	heads, err := orm.LatestHeads(context.TODO(), 100)

	require.Zero(t, len(heads))
	require.NoError(t, err)
}
