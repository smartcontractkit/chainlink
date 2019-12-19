package orm_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestNewLockingStrategy(t *testing.T) {
	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()
	c := tc.Config

	tests := []struct {
		name        string
		dialectName orm.DialectName
		path        string
		expect      reflect.Type
	}{
		{"sqlite", orm.DialectSqlite, c.RootDir(), reflect.ValueOf(&orm.FileLockingStrategy{}).Type()},
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

func TestFileLockingStrategy_Lock(t *testing.T) {
	const delay = 10 * time.Millisecond

	tc, cleanup := cltest.NewConfig(t)
	defer cleanup()
	c := tc.Config

	require.NoError(t, os.MkdirAll(c.RootDir(), 0700))
	defer os.RemoveAll(c.RootDir())

	dbpath := filepath.ToSlash(filepath.Join(c.RootDir(), "db.sqlite3"))
	ls, err := orm.NewFileLockingStrategy(dbpath)
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")

	ls2, err := orm.NewFileLockingStrategy(dbpath)
	require.NoError(t, err)
	require.Error(t, ls2.Lock(delay), "should not get 2nd exclusive lock")

	require.NoError(t, ls.Unlock(delay))

	require.NoError(t, ls2.Lock(delay), "allow another to lock after unlock")
	require.NoError(t, ls2.Unlock(delay))
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
		t.Skip("No postgres DatabaseURL set.")
	}

	connErr, dbErr := store.ORM.LockingStrategyHelperSimulateDisconnect()
	require.NoError(t, connErr)
	require.NoError(t, dbErr)

	orm2, err := orm.NewORM(orm.NormalizedDatabaseURL(store.Config), store.Config.DatabaseTimeout())
	require.NoError(t, err)
	defer orm2.Close()

	err = orm2.RawDB(func(db *gorm.DB) error {
		return db.Save(&models.JobSpec{ID: models.NewID()}).Error
	})
	require.NoError(t, err)

	require.Panics(t, func() {
		_ = store.ORM.RawDB(func(db *gorm.DB) error { return nil })
	})
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

	require.Panics(t, func() {
		_ = store.ORM.RawDB(func(db *gorm.DB) error { return nil })
	})
}
