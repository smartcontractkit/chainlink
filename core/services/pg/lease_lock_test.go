package pg_test

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func newLeaseLock(t *testing.T, db *sqlx.DB) pg.LeaseLock {
	return pg.NewLeaseLock(db, uuid.NewV4(), logger.TestLogger(t), 1*time.Second, 5*time.Second)
}

func Test_LeaseLock(t *testing.T) {
	t.Run("on migrated database", func(t *testing.T) {
		_, db := heavyweight.FullTestDB(t, "leaselock", true, false)

		leaseLock1 := newLeaseLock(t, db)

		err := leaseLock1.TakeAndHold()
		require.NoError(t, err)

		var clientID uuid.UUID
		err = db.Get(&clientID, `SELECT client_id FROM lease_lock`)
		require.NoError(t, err)
		assert.Equal(t, leaseLock1.ClientID(), clientID)

		started2 := make(chan struct{})
		leaseLock2 := newLeaseLock(t, db)
		go func() {
			defer leaseLock2.Release()
			err := leaseLock2.TakeAndHold()
			require.NoError(t, err)
			close(started2)
		}()

		time.Sleep(2 * time.Second)

		leaseLock1.Release()

		select {
		case <-started2:
		case <-time.After(cltest.WaitTimeout(t)):
			t.Fatal("timed out waiting for leaseLock2 to start")
		}

		err = db.Get(&clientID, `SELECT client_id FROM lease_lock`)
		require.NoError(t, err)
		assert.Equal(t, leaseLock2.ClientID(), clientID)
	})

	t.Run("on virgin database", func(t *testing.T) {
		_, db := heavyweight.FullTestDB(t, "leaselock", false, false)

		leaseLock1 := newLeaseLock(t, db)

		err := leaseLock1.TakeAndHold()
		defer leaseLock1.Release()
		require.NoError(t, err)
	})
}
