package orm_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/orm"
	"github.com/stretchr/testify/require"
)

func TestORM_NewLockingStrategy(t *testing.T) {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	c := tc.Config

	tests := []struct {
		name        string
		dialectName orm.DialectName
		path        string
		expect      reflect.Type
		wantError   bool
	}{
		{"sqlite", orm.DialectSqlite, c.RootDir(), reflect.ValueOf(&orm.FileLockingStrategy{}).Type(), false},
		{"sqlite bad path", orm.DialectSqlite, ":/\\fd/8970382094", nil, true},
		{"postgres", orm.DialectPostgres, "postgres://something:5432", reflect.ValueOf(&orm.PostgresLockingStrategy{}).Type(), false},
	}

	for _, test := range tests {
		t.Run(string(test.name), func(t *testing.T) {
			rval, err := orm.NewLockingStrategy(test.dialectName, test.path)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				rtype := reflect.ValueOf(rval).Type()
				require.Equal(t, test.expect, rtype)
			}
		})
	}
}

const delay = 10 * time.Millisecond

func TestORM_FileLockingStrategy_Lock(t *testing.T) {
	tc, cleanup := cltest.NewConfig()
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

	require.NoError(t, ls.Unlock())

	require.NoError(t, ls2.Lock(delay), "allow another to lock after unlock")
	require.NoError(t, ls2.Unlock())
}

func TestORM_PostgresLockingStrategy_Lock(t *testing.T) {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	c := tc.Config

	if c.DatabaseURL() == "" {
		t.Skip("No postgres DatabaseURL set.")
	}

	ls, err := orm.NewPostgresLockingStrategy(c.DatabaseURL())
	require.NoError(t, err)
	require.NoError(t, ls.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls.Lock(delay), "relocking on same instance is noop")

	ls2, err := orm.NewPostgresLockingStrategy(c.DatabaseURL())
	require.NoError(t, err)
	require.Error(t, ls2.Lock(delay), "should not get 2nd exclusive lock")
	require.NoError(t, ls2.Unlock())

	require.NoError(t, ls.Unlock())
	require.NoError(t, ls2.Lock(delay), "should get exclusive lock")
	require.NoError(t, ls2.Unlock())
}
