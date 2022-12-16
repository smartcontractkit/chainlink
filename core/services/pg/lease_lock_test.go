package pg_test

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func newLeaseLock(t *testing.T, db *sqlx.DB, cfg config.GeneralConfig) pg.LeaseLock {
	return pg.NewLeaseLock(db, uuid.NewV4(), logger.TestLogger(t), cfg)
}

func Test_LeaseLock(t *testing.T) {
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, "leaselock", func(c *chainlink.Config, s *chainlink.Secrets) {
		t := true
		c.Database.Lock.Enabled = &t
		c.Database.Lock.LeaseDuration = models.MustNewDuration(15 * time.Second)
		c.Database.Lock.LeaseRefreshInterval = models.MustNewDuration(100 * time.Millisecond)
	})

	t.Run("on migrated database", func(t *testing.T) {
		leaseLock1 := newLeaseLock(t, db, cfg)

		err := leaseLock1.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		var clientID uuid.UUID
		err = db.Get(&clientID, `SELECT client_id FROM lease_lock`)
		require.NoError(t, err)
		assert.Equal(t, leaseLock1.ClientID(), clientID)

		started2 := make(chan struct{})
		leaseLock2 := newLeaseLock(t, db, cfg)
		go func() {
			defer leaseLock2.Release()
			err := leaseLock2.TakeAndHold(testutils.Context(t))
			require.NoError(t, err)
			close(started2)
		}()

		// Give it plenty of time to have a few tries at getting the lease
		time.Sleep(cfg.LeaseLockRefreshInterval() * 5)

		leaseLock1.Release()

		select {
		case <-started2:
		case <-time.After(testutils.WaitTimeout(t)):
			t.Fatal("timed out waiting for leaseLock2 to start")
		}

		err = db.Get(&clientID, `SELECT client_id FROM lease_lock`)
		require.NoError(t, err)
		assert.Equal(t, leaseLock2.ClientID(), clientID)
	})

	t.Run("recovers and re-opens connection if it's closed externally on initial take wait", func(t *testing.T) {
		leaseLock := newLeaseLock(t, db, cfg)

		otherAppID := uuid.NewV4()

		// simulate another application holding lease to force it to retry
		res, err := db.Exec(`UPDATE lease_lock SET client_id=$1,expires_at=NOW()+'1 day'::interval`, otherAppID)
		require.NoError(t, err)
		rowsAffected, err := res.RowsAffected()
		require.NoError(t, err)
		require.EqualValues(t, 1, rowsAffected)

		conn, err := db.Connx(testutils.Context(t))
		require.NoError(t, err)

		pg.SetConn(leaseLock, conn)

		// Simulate the connection being closed (leaseLock should automatically check out a new one)
		require.NoError(t, conn.Close())

		gotLease := make(chan struct{})
		go func() {
			errInternal := leaseLock.TakeAndHold(testutils.Context(t))
			require.NoError(t, errInternal)
			close(gotLease)
		}()

		// Give it plenty of time to have a few tries at getting the lease
		time.Sleep(cfg.LeaseLockRefreshInterval() * 5)

		// Release the dummy lease lock to allow the lease locker to take it now
		_, err = db.Exec(`DELETE FROM lease_lock WHERE client_id=$1`, otherAppID)
		require.NoError(t, err)

		select {
		case <-gotLease:
		case <-time.After(testutils.WaitTimeout(t)):
			t.Fatal("timed out waiting for lease lock to start")
		}

		// check that the lease lock was actually taken
		var exists bool
		err = db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM lease_lock)`)
		require.NoError(t, err)

		assert.True(t, exists)

		leaseLock.Release()
	})

	t.Run("recovers and re-opens connection if it's closed externally while holding", func(t *testing.T) {
		leaseLock := newLeaseLock(t, db, cfg)

		err := leaseLock.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)
		defer leaseLock.Release()

		conn := pg.GetConn(leaseLock)

		var prevExpiresAt time.Time

		err = conn.Close()
		require.NoError(t, err)

		err = db.Get(&prevExpiresAt, `SELECT expires_at FROM lease_lock`)
		require.NoError(t, err)

		time.Sleep(cfg.LeaseLockRefreshInterval() + 1*time.Second)

		var expiresAt time.Time

		err = db.Get(&expiresAt, `SELECT expires_at FROM lease_lock`)
		require.NoError(t, err)

		// The lease lock must have recovered and re-opened the connection if the second expires_at is later
		assert.Greater(t, expiresAt.Unix(), prevExpiresAt.Unix())
	})

	t.Run("release lock with Release() func", func(t *testing.T) {
		leaseLock := newLeaseLock(t, db, cfg)

		err := leaseLock.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		leaseLock.Release()

		leaseLock2 := newLeaseLock(t, db, cfg)
		err = leaseLock2.TakeAndHold(testutils.Context(t))
		defer leaseLock2.Release()
		require.NoError(t, err)
	})

	t.Run("cancel TakeAndHold with ctx", func(t *testing.T) {
		leaseLock1 := newLeaseLock(t, db, cfg)
		leaseLock2 := newLeaseLock(t, db, cfg)

		err := leaseLock1.TakeAndHold(testutils.Context(t))
		require.NoError(t, err)

		awaiter := cltest.NewAwaiter()
		go func() {
			ctx, cancel := context.WithCancel(testutils.Context(t))
			go func() {
				<-time.After(3 * time.Second)
				cancel()
			}()
			err := leaseLock2.TakeAndHold(ctx)
			require.Error(t, err)
			awaiter.ItHappened()
		}()

		awaiter.AwaitOrFail(t)
		leaseLock1.Release()
	})

	require.NoError(t, db.Close())

	t.Run("on virgin database", func(t *testing.T) {
		_, db := heavyweight.FullTestDBEmptyV2(t, "leaselock", nil)

		leaseLock1 := newLeaseLock(t, db, cfg)

		err := leaseLock1.TakeAndHold(testutils.Context(t))
		defer leaseLock1.Release()
		require.NoError(t, err)
	})
}
