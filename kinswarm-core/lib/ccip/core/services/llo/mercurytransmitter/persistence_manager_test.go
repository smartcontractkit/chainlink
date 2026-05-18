package mercurytransmitter

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func bootstrapPersistenceManager(t *testing.T, donID uint32, db *sqlx.DB) (*persistenceManager, *observer.ObservedLogs) {
	t.Helper()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	orm := NewORM(db, donID)
	return NewPersistenceManager(lggr, orm, "wss://example.com/mercury", 2, 5*time.Millisecond, 5*time.Millisecond), observedLogs
}

func TestPersistenceManager(t *testing.T) {
	donID1 := uint32(1234)
	donID2 := uint32(2345)

	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	pm, _ := bootstrapPersistenceManager(t, donID1, db)

	transmissions := makeSampleTransmissions()
	err := pm.orm.Insert(ctx, transmissions)
	require.NoError(t, err)

	result, err := pm.Load(ctx)
	require.NoError(t, err)
	assert.ElementsMatch(t, transmissions, result)

	err = pm.orm.Delete(ctx, [][32]byte{transmissions[0].Hash()})
	require.NoError(t, err)

	t.Run("scopes load to only transmissions with matching don ID", func(t *testing.T) {
		pm2, _ := bootstrapPersistenceManager(t, donID2, db)
		result, err = pm2.Load(ctx)
		require.NoError(t, err)

		assert.Len(t, result, 0)
	})
}

func TestPersistenceManagerAsyncDelete(t *testing.T) {
	ctx := testutils.Context(t)
	donID := uint32(1234)
	db := pgtest.NewSqlxDB(t)
	pm, observedLogs := bootstrapPersistenceManager(t, donID, db)

	transmissions := makeSampleTransmissions()
	err := pm.orm.Insert(ctx, transmissions)
	require.NoError(t, err)

	servicetest.Run(t, pm)

	pm.AsyncDelete(transmissions[0].Hash())

	// Wait for next poll.
	observedLogs.TakeAll()
	testutils.WaitForLogMessage(t, observedLogs, "Deleted queued transmit requests")

	result, err := pm.Load(ctx)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.ElementsMatch(t, transmissions[1:], result)
}

func TestPersistenceManagerPrune(t *testing.T) {
	donID1 := uint32(123456)
	donID2 := uint32(654321)
	db := pgtest.NewSqlxDB(t)

	ctx := testutils.Context(t)

	transmissions := make([]*Transmission, 45)
	for i := uint64(0); i < 45; i++ {
		transmissions[i] = makeSampleTransmission(i)
	}

	pm, _ := bootstrapPersistenceManager(t, donID1, db)
	err := pm.orm.Insert(ctx, transmissions[:25])
	require.NoError(t, err)

	pm2, _ := bootstrapPersistenceManager(t, donID2, db)
	err = pm2.orm.Insert(ctx, transmissions[25:])
	require.NoError(t, err)

	pm, observedLogs := bootstrapPersistenceManager(t, donID1, db)

	err = pm.Start(ctx)
	require.NoError(t, err)

	// Wait for next poll.
	observedLogs.TakeAll()
	testutils.WaitForLogMessage(t, observedLogs, "Pruned transmit requests table")

	result, err := pm.Load(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, transmissions[23:25], result)

	// Test pruning stops after Close.
	err = pm.Close()
	require.NoError(t, err)

	err = pm.orm.Insert(ctx, transmissions)
	require.NoError(t, err)

	result, err = pm.Load(ctx)
	require.NoError(t, err)
	require.Len(t, result, 25)

	t.Run("prune was scoped to don ID", func(t *testing.T) {
		result, err = pm2.Load(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 20)
	})
}
