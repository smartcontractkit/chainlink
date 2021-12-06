package headtracker_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_Heads_IdempotentInsertHead(t *testing.T) {
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
