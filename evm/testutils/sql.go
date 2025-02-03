package testutils

import (
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func NewSqlxDB(t testing.TB) *sqlx.DB {
	SkipShortDB(t)
	dbURL := os.Getenv("CL_DATABASE_URL")
	if dbURL == "" {
		t.Errorf("you must provide a CL_DATABASE_URL environment variable")
		return nil
	}
	return sqltest.NewDB(t, dbURL)
}

// SkipShortDB skips tb during -short runs, and notes the DB dependency.
func SkipShortDB(tb testing.TB) {
	tests.SkipShort(tb, "DB dependency")
}

func MustExec(t *testing.T, ds sqlutil.DataSource, stmt string, args ...interface{}) {
	require.NoError(t, utils.JustError(ds.ExecContext(Context(t), stmt, args...)))
}

// pristineDBName is a clean copy of test DB with migrations.
const pristineDBName = "chainlink_test_pristine" //TODO update when splitting schemas

// NewIndependentSqlxDB return a new independent test database, which does not use txdb and therefore supports txs etc.
// Use this with caution, as it is much more costly than NewSqlxDB.
func NewIndependentSqlxDB(t testing.TB) *sqlx.DB {
	SkipShortDB(t)

	ctx := tests.Context(t)

	dbURL, err := url.Parse(sqltest.TestURL(t))
	require.NoError(t, err)
	dbName := "chainlink_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	newDBURL := sqltest.CreateOrReplace(t, *dbURL, dbName, pristineDBName)

	db, err := pg.DBConfig{
		IdleInTxSessionTimeout: time.Hour,
		LockTimeout:            15 * time.Second,
		MaxOpenConns:           100,
		MaxIdleConns:           10,
	}.New(ctx, newDBURL.String(), pg.DriverPostgres)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	return db
}
