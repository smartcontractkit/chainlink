package mercurytransmitter

import (
	"sort"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func bootstrapPersistenceManager(t *testing.T, donID uint32, db *sqlx.DB, maxTransmitQueueSize int) (*persistenceManager, *observer.ObservedLogs) {
	t.Helper()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	orm := NewORM(db, donID)
	return NewPersistenceManager(lggr, orm, "wss://example.com/mercury", maxTransmitQueueSize, 5*time.Millisecond, 5*time.Millisecond, 30*24*time.Hour), observedLogs
}

func TestPersistenceManager(t *testing.T) {
	donID1 := uint32(1234)
	donID2 := uint32(2345)

	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)

	t.Run("loads transmissions", func(t *testing.T) {
		pm, _ := bootstrapPersistenceManager(t, donID1, db, 2)
		transmissions := makeSampleTransmissions(3, sURL)
		err := pm.orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		sort.Slice(transmissions, func(i, j int) bool {
			// sort by seqnr desc to match return of Get
			return transmissions[i].SeqNr > transmissions[j].SeqNr
		})
		result, err := pm.Load(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, transmissions[0:2], result)

		err = pm.orm.Delete(ctx, [][32]byte{transmissions[0].Hash()})
		require.NoError(t, err)
	})

	t.Run("scopes load to only transmissions with matching don ID", func(t *testing.T) {
		pm, _ := bootstrapPersistenceManager(t, donID1, db, 2)
		transmissions := makeSampleTransmissions(3, sURL)
		err := pm.orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		pm2, _ := bootstrapPersistenceManager(t, donID2, db, 3)
		result, err := pm2.Load(ctx)
		require.NoError(t, err)

		assert.Empty(t, result)
	})

	t.Run("does not load records older than maxAge", func(t *testing.T) {
		pm, _ := bootstrapPersistenceManager(t, donID1, db, 3)
		transmissions := makeSampleTransmissions(3, sURL)
		err := pm.orm.Insert(ctx, transmissions)
		require.NoError(t, err)

		pgtest.MustExec(t, db, `UPDATE llo_mercury_transmit_queue SET inserted_at = NOW() - INTERVAL '1 year' WHERE seq_nr = 0`)

		result, err := pm.Load(ctx)
		require.NoError(t, err)

		assert.Len(t, result, 2)
		assert.Equal(t, uint64(2), result[0].SeqNr)
		assert.Equal(t, uint64(1), result[1].SeqNr)
	})
}

func TestPersistenceManagerAsyncDelete(t *testing.T) {
	ctx := testutils.Context(t)
	donID := uint32(1234)
	db := pgtest.NewSqlxDB(t)
	pm, observedLogs := bootstrapPersistenceManager(t, donID, db, 1000)

	transmissions := makeSampleTransmissions(3, sURL)
	err := pm.orm.Insert(ctx, transmissions)
	require.NoError(t, err)

	servicetest.Run(t, pm)

	pm.AsyncDelete(transmissions[0].Hash())

	// Wait for next poll.
	observedLogs.TakeAll()
	testutils.WaitForLogMessage(t, observedLogs, "Flushed delete queue")

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
		transmissions[i] = makeSampleTransmission(i, sURL, ocrtypes.Report{byte(i)})
	}

	// cut 25 down to 2
	pm, observedLogs := bootstrapPersistenceManager(t, donID1, db, 2)
	err := pm.orm.Insert(ctx, transmissions[:25])
	require.NoError(t, err)

	pm2, _ := bootstrapPersistenceManager(t, donID2, db, 20)
	err = pm2.orm.Insert(ctx, transmissions[25:])
	require.NoError(t, err)

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
	require.Len(t, result, 2)

	t.Run("prune was scoped to don ID", func(t *testing.T) {
		result, err = pm2.Load(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 20)
	})
}

func Test_PersistenceManager_deleteTransmissions(t *testing.T) {
	donID1 := uint32(123456)
	db := pgtest.NewSqlxDB(t)

	ctx := testutils.Context(t)

	transmissions := make([]*Transmission, 45)
	for i := uint64(0); i < 45; i++ {
		transmissions[i] = makeSampleTransmission(i, sURL, ocrtypes.Report{byte(i)})
	}

	pm, _ := bootstrapPersistenceManager(t, donID1, db, 1000)
	require.NoError(t, pm.orm.Insert(ctx, transmissions))

	hashesToDelete := make([][32]byte, 20)
	for i := 0; i < 20; i++ {
		hashesToDelete[i] = transmissions[i].Hash()
	}
	pm.deleteTransmissions(ctx, hashesToDelete, 7)

	ts, err := pm.Load(ctx)
	require.NoError(t, err)

	require.Len(t, ts, 25)
	for i := 0; i < 20; i++ {
		assert.NotContains(t, ts, transmissions[i])
	}
	for i := 20; i < 45; i++ {
		assert.Contains(t, ts, transmissions[i])
	}
}
