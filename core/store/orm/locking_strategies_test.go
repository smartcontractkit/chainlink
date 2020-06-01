package orm_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/jinzhu/gorm"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestNewLockingStrategy(t *testing.T) {
	tests := []struct {
		name        string
		dialectName orm.DialectName
		path        string
		expect      reflect.Type
	}{
		{"postgres", orm.DialectPostgres, "postgres://something:5432", reflect.ValueOf(&orm.PostgresLockingStrategy{}).Type()},
	}

	for _, test := range tests {
		t.Run(string(test.name), func(t *testing.T) {
			connectionType, err := orm.NewConnection(orm.DialectPostgres, test.path, 42)
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
	delay := tc.DatabaseTimeout()
	if tc.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	withLock, err := orm.NewConnection(orm.DialectPostgres, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault())
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
	if tc.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	withLock, err := orm.NewConnection(orm.DialectPostgres, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault())
	require.NoError(t, err)
	ls, err := orm.NewPostgresLockingStrategy(withLock)
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls.Lock(delay), "relocking on same instance is reentrant")

	withoutLock, err := orm.NewConnection(orm.DialectPostgresWithoutLock, tc.DatabaseURL(), tc.GetAdvisoryLockIDConfiguredOrDefault())
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
	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	delay := store.Config.DatabaseTimeout()

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	err := store.ORM.RawDB(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewID()}).Error
	})
	require.NoError(t, err)

	ct, err := orm.NewConnection(orm.DialectPostgres, store.Config.DatabaseURL(), tc.Config.GetAdvisoryLockIDConfiguredOrDefault())
	require.NoError(t, err)
	lock2, err := orm.NewLockingStrategy(ct)
	require.NoError(t, err)
	err = lock2.Lock(delay)
	require.Equal(t, errors.Cause(err), orm.ErrNoAdvisoryLock)
	defer lock2.Unlock(delay)
}

func TestPostgresLockingStrategy_CanBeReacquiredByNewNodeAfterDisconnect(t *testing.T) {
	tc := cltest.NewTestConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	orm2ShutdownSignal := gracefulpanic.NewSignal()
	orm2, err := orm.NewORM(store.Config.DatabaseURL(), store.Config.DatabaseTimeout(), orm2ShutdownSignal, orm.DialectTransactionWrappedPostgres, tc.Config.GetAdvisoryLockIDConfiguredOrDefault())
	require.NoError(t, err)
	defer orm2.Close()

	err = orm2.RawDB(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewID()}).Error
	})
	require.NoError(t, err)

	_ = store.ORM.RawDB(func(db *gorm.DB) error { return nil })
	gomega.NewGomegaWithT(t).Eventually(store.ORM.ShutdownSignal().Wait()).Should(gomega.BeClosed())
}

func TestPostgresLockingStrategy_WhenReacquiredOriginalNodeErrors(t *testing.T) {
	tc := cltest.NewTestConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(tc)
	defer cleanup()

	delay := store.Config.DatabaseTimeout()

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	ct, err := orm.NewConnection(orm.DialectPostgres, store.Config.DatabaseURL(), tc.Config.GetAdvisoryLockIDConfiguredOrDefault())
	require.NoError(t, err)
	lock, err := orm.NewLockingStrategy(ct)
	require.NoError(t, err)
	defer lock.Unlock(delay)

	err = lock.Lock(models.MustMakeDuration(1 * time.Second))
	require.NoError(t, err)
	defer lock.Unlock(delay)

	_ = store.ORM.RawDB(func(db *gorm.DB) error { return nil })
	gomega.NewGomegaWithT(t).Eventually(store.ORM.ShutdownSignal().Wait()).Should(gomega.BeClosed())
}
