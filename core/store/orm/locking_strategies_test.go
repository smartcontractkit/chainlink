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
			rval, err := orm.NewLockingStrategy(test.dialectName, test.path)
			require.NoError(t, err)
			rtype := reflect.ValueOf(rval).Type()
			require.Equal(t, test.expect, rtype)
		})
	}
}

func TestPostgresLockingStrategy_Lock(t *testing.T) {
	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()

	cleanupDB := cltest.PrepareTestDB(tc)
	defer cleanupDB()

	c := tc.Config

	if c.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	delay := c.DatabaseTimeout()

	ls, err := orm.NewPostgresLockingStrategy(c.DatabaseURL())
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls.Lock(delay), "relocking on same instance is reentrant")

	ls2, err := orm.NewPostgresLockingStrategy(c.DatabaseURL())
	require.NoError(t, err)
	require.Error(t, ls2.Lock(delay), "should not get 2nd exclusive lock")

	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls.Unlock(delay))
	require.NoError(t, ls2.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls2.Unlock(delay))
}

func TestPostgresLockingStrategy_WhenLostIsReacquired(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	if store.Config.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	delay := store.Config.DatabaseTimeout()

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	err := store.ORM.RawDB(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewID()}).Error
	})
	require.NoError(t, err)

	lock2, err := orm.NewLockingStrategy("postgres", store.Config.DatabaseURL())
	require.NoError(t, err)
	err = lock2.Lock(delay)
	require.Equal(t, errors.Cause(err), orm.ErrNoAdvisoryLock)
	defer lock2.Unlock(delay)
}

func TestPostgresLockingStrategy_CanBeReacquiredByNewNodeAfterDisconnect(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	if store.Config.DatabaseURL() == "" {
		panic("No postgres DatabaseURL set.")
	}

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	orm2ShutdownSignal := gracefulpanic.NewSignal()
	orm2, err := orm.NewORM(store.Config.DatabaseURL(), store.Config.DatabaseTimeout(), orm2ShutdownSignal)
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
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	if store.Config.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	delay := store.Config.DatabaseTimeout()

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	lock, err := orm.NewLockingStrategy("postgres", store.Config.DatabaseURL())
	require.NoError(t, err)
	defer lock.Unlock(delay)

	err = lock.Lock(1 * time.Second)
	require.NoError(t, err)
	defer lock.Unlock(delay)

	_ = store.ORM.RawDB(func(db *gorm.DB) error { return nil })
	gomega.NewGomegaWithT(t).Eventually(store.ORM.ShutdownSignal().Wait()).Should(gomega.BeClosed())
}
