package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func newAdvisoryLock(t *testing.T, db *sqlx.DB, cfg pg.AdvisoryLockConfig) pg.AdvisoryLock {
	return pg.NewAdvisoryLock(db, logger.TestLogger(t), cfg)
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func Test_AdvisoryLock(t *testing.T) {
	cfg, db := heavyweight.FullTestDBEmpty(t, "advisorylock")
	check := 1 * time.Second
	cfg.Overrides.AdvisoryLockCheckInterval = &check

	t.Run("takes lock", func(t *testing.T) {
		advLock1 := newAdvisoryLock(t, db, cfg)

		err := advLock1.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		var lockTaken bool
		err = db.Get(&lockTaken, `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND (classid::bigint << 32) | objid::bigint = $1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)
		assert.True(t, lockTaken)

		started2 := make(chan struct{})
		advLock2 := newAdvisoryLock(t, db, cfg)
		go func() {
			err := advLock2.TakeAndHold(testutils.Context(t))
			require.NoError(t, err)
			close(started2)
		}()

		// Give it plenty of time for advLock2 to have a few tries at getting the lease
		time.Sleep(cfg.AdvisoryLockCheckInterval() * 5)

		advLock1.Release()

		select {
		case <-started2:
		case <-time.After(testutils.WaitTimeout(t)):
			t.Fatal("timed out waiting for advLock2 to start")
		}

		err = db.Get(&lockTaken, `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND (classid::bigint << 32) | objid::bigint = $1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)
		assert.True(t, lockTaken)

		advLock2.Release()

		// pg_locks is not atomic
		time.Sleep(100 * time.Millisecond)

		err = db.Get(&lockTaken, `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND (classid::bigint << 32) | objid::bigint = $1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)
		assert.False(t, lockTaken)
	})

	t.Run("recovers and re-opens connection if it's closed externally on initial take wait", func(t *testing.T) {
		advLock := newAdvisoryLock(t, db, cfg)

		// simulate another application holding advisory lock to force it to retry
		ctx, cancel := context.WithTimeout(testutils.Context(t), cfg.DatabaseDefaultQueryTimeout())
		defer cancel()
		conn, err := db.Conn(ctx)
		require.NoError(t, err)
		_, err = conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)

		conn2, err := db.Connx(testutils.Context(t))
		require.NoError(t, err)

		pg.SetConn(advLock, conn2)

		// Simulate the connection being closed (advLock should automatically check out a new one)
		require.NoError(t, conn2.Close())

		gotLease := make(chan struct{})
		go func() {
			err := advLock.TakeAndHold(testutils.Context(t))
			require.NoError(t, err)
			close(gotLease)
		}()

		// Give it plenty of time to have a few tries at getting the lock
		time.Sleep(cfg.AdvisoryLockCheckInterval() * 5)

		// Release the dummy advisory lock to allow the lease locker to take it now
		_, err = conn.ExecContext(testutils.Context(t), `SELECT pg_advisory_unlock($1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)

		select {
		case <-gotLease:
		case <-time.After(testutils.WaitTimeout(t)):
			t.Fatal("timed out waiting for lease lock to start")
		}

		// check that the lease lock was actually taken
		var exists bool
		err = db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM pg_locks WHERE locktype = 'advisory' AND (classid::bigint << 32) | objid::bigint = $1)`, cfg.AdvisoryLockID())
		require.NoError(t, err)

		assert.True(t, exists)

		advLock.Release()
	})

	t.Run("release lock with Release() func", func(t *testing.T) {
		advisoryLock := newAdvisoryLock(t, db, cfg)

		err := advisoryLock.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		advisoryLock.Release()

		advisoryLock2 := newAdvisoryLock(t, db, cfg)
		defer advisoryLock2.Release()
		err = advisoryLock2.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)
	})

	t.Run("cancel TakeAndHold with ctx", func(t *testing.T) {
		advisoryLock1 := newAdvisoryLock(t, db, cfg)
		advisoryLock2 := newAdvisoryLock(t, db, cfg)

		err := advisoryLock1.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		awaiter := cltest.NewAwaiter()
		go func() {
			ctx, cancel := context.WithCancel(testutils.Context(t))
			go func() {
				<-time.After(3 * time.Second)
				cancel()
			}()
			err := advisoryLock2.TakeAndHold(ctx)
			require.Error(t, err)
			awaiter.ItHappened()
		}()

		awaiter.AwaitOrFail(t)
		advisoryLock1.Release()
	})

	require.NoError(t, db.Close())
}
