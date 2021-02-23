package orm_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

func TestNewLockingStrategy(t *testing.T) {
	tests := []struct {
		name        string
		dialectName dialects.DialectName
		path        string
		expect      reflect.Type
	}{
		{"postgres", dialects.Postgres, "postgres://something:5432", reflect.ValueOf(&orm.PostgresLockingStrategy{}).Type()},
	}

	for _, test := range tests {
		t.Run(string(test.name), func(t *testing.T) {
			connectionType, err := orm.NewConnection(dialects.Postgres, test.path, 42, 1*time.Second, 0, 0)
			require.NoError(t, err)
			rval, err := orm.NewLockingStrategy(connectionType)
			require.NoError(t, err)
			rtype := reflect.ValueOf(rval).Type()
			require.Equal(t, test.expect, rtype)
		})
	}
}

func TestPostgresLockingStrategy_Lock_withLock(t *testing.T) {
	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()

	tc.Config.Set("DATABASE_TIMEOUT", "500ms")
	delay := tc.DatabaseTimeout()
	if tc.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	withLock, err := orm.NewConnection(dialects.Postgres, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault(), tc.GlobalLockRetryInterval().Duration(), tc.ORMMaxOpenConns(), tc.ORMMaxIdleConns())
	require.NoError(t, err)
	ls, err := orm.NewPostgresLockingStrategy(withLock)
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls.Lock(delay), "relocking on same instance is reentrant")

	ls2, err := orm.NewPostgresLockingStrategy(withLock)
	require.NoError(t, err)
	require.Error(t, ls2.Lock(delay), "should not get 2nd exclusive lock")

	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls2.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls2.Unlock(delay))
}

func TestPostgresLockingStrategy_Lock_withoutLock(t *testing.T) {
	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()
	delay := tc.DatabaseTimeout()

	tc.Config.Set("DATABASE_TIMEOUT", "500ms")
	if tc.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	withLock, err := orm.NewConnection(dialects.Postgres, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault(), tc.GlobalLockRetryInterval().Duration(), tc.ORMMaxOpenConns(), tc.ORMMaxIdleConns())
	require.NoError(t, err)
	ls, err := orm.NewPostgresLockingStrategy(withLock)
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls.Lock(delay), "relocking on same instance is reentrant")

	withoutLock, err := orm.NewConnection(dialects.PostgresWithoutLock, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault(), tc.GlobalLockRetryInterval().Duration(), tc.ORMMaxOpenConns(), tc.ORMMaxIdleConns())
	require.NoError(t, err)
	ls2, err := orm.NewPostgresLockingStrategy(withoutLock)
	require.NoError(t, err)
	require.NoError(t, ls2.Lock(delay), "should not wait for lock")

	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls2.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls2.Unlock(delay))
}

func TestPostgresLockingStrategy_WhenLostIsReacquired(t *testing.T) {
	tc := cltest.NewTestConfig(t)
	tc.Config.Set("DATABASE_TIMEOUT", "500ms")

	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	delay := store.Config.DatabaseTimeout()

	// NewStore no longer takes a lock on opening, so do something that does...
	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewJobID()}).Error
	})
	require.NoError(t, err)

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	err = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewJobID()}).Error
	})
	require.NoError(t, err)

	ct, err := orm.NewConnection(dialects.Postgres, store.Config.DatabaseURL(), tc.Config.GetAdvisoryLockIDConfiguredOrDefault(), 10*time.Millisecond, 0, 0)
	require.NoError(t, err)
	lock2, err := orm.NewLockingStrategy(ct)
	require.NoError(t, err)
	err = lock2.Lock(delay)
	require.Equal(t, errors.Cause(err), orm.ErrNoAdvisoryLock)
	defer lock2.Unlock(delay)
}

func TestPostgresLockingStrategy_CanBeReacquiredByNewNodeAfterDisconnect(t *testing.T) {
	tc := cltest.NewTestConfig(t)
	tc.Config.Set("DATABASE_TIMEOUT", "500ms")
	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	// NewStore no longer takes a lock on opening, so do something that does...
	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewJobID()}).Error
	})
	require.NoError(t, err)

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	orm2ShutdownSignal := gracefulpanic.NewSignal()
	orm2, err := orm.NewORM(store.Config.DatabaseURL(), store.Config.DatabaseTimeout(), orm2ShutdownSignal, dialects.TransactionWrappedPostgres, tc.Config.GetAdvisoryLockIDConfiguredOrDefault(), tc.Config.GlobalLockRetryInterval().Duration(), tc.ORMMaxOpenConns(), tc.ORMMaxIdleConns())
	require.NoError(t, err)
	defer orm2.Close()

	err = orm2.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewJobID()}).Error
	})
	require.NoError(t, err)

	_ = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error { return nil })
	gomega.NewGomegaWithT(t).Eventually(store.ORM.ShutdownSignal().Wait()).Should(gomega.BeClosed())
}

func TestPostgresLockingStrategy_WhenReacquiredOriginalNodeErrors(t *testing.T) {
	tc := cltest.NewTestConfig(t)
	tc.Config.Set("DATABASE_TIMEOUT", "500ms")
	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	delay := store.Config.DatabaseTimeout()

	// NewStore no longer takes a lock on opening, so do something that does...
	err := store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewJobID()}).Error
	})
	require.NoError(t, err)

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	ct, err := orm.NewConnection(dialects.Postgres, store.Config.DatabaseURL(), tc.Config.GetAdvisoryLockIDConfiguredOrDefault(), tc.Config.GlobalLockRetryInterval().Duration(), tc.ORMMaxOpenConns(), tc.ORMMaxIdleConns())
	require.NoError(t, err)
	lock, err := orm.NewLockingStrategy(ct)
	require.NoError(t, err)
	defer lock.Unlock(delay)

	err = lock.Lock(models.MustMakeDuration(1 * time.Second))
	require.NoError(t, err)
	defer lock.Unlock(delay)

	_ = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error { return nil })
	gomega.NewGomegaWithT(t).Eventually(store.ORM.ShutdownSignal().Wait()).Should(gomega.BeClosed())
}
