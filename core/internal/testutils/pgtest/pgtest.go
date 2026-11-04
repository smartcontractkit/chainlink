package pgtest

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

func NewSqlxDB(t testing.TB) *sqlx.DB {
	testutils.SkipShortDB(t)
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		t.Fatalf("you must provide a CL_DATABASE_URL environment variable")
	}
	db := sqltest.NewDB(t, dbURL)

	// Prevent parallel txdb tests from blocking indefinitely on lock contention.
	// sqltest.NewDB does not run any init SQL, so without this a session will wait
	// forever for locks held by other txdb-wrapped tests (whose transactions stay
	// open for the full test lifetime).
	_, err := db.Exec(`SET lock_timeout = '15s';
SET idle_in_transaction_session_timeout = '30s';
SET statement_timeout = '30s';`)
	require.NoError(t, err, "failed to set session timeouts on test DB")

	opened := time.Now()
	t.Cleanup(func() {
		if elapsed := time.Since(opened); elapsed > 2*time.Minute {
			t.Logf("pgtest: txdb connection held for a long time: %s (opened at %s). If tests are failing or hanging, there might be issues with how you're accessing the DB that lock out others. You can also consider increasing the lock timeout.", elapsed.Round(time.Second), opened.Format(time.RFC3339))
		}
	})

	return db
}

func MustExec(t *testing.T, ds sqlutil.DataSource, stmt string, args ...any) {
	ctx := t.Context()
	require.NoError(t, utils.JustError(ds.ExecContext(ctx, stmt, args...)))
}

// NewIndependentSqlxSessions returns n handles to a fresh, migrated test database, each pinned to its own
// connection (pool size 1). Unlike NewSqlxDB, operations are not wrapped in a single test transaction and
// session state is not shared, so concurrent sessions can observe and contend for each other's locks.
// Provisioning a database is expensive, so use this only for tests that need real concurrency.
func NewIndependentSqlxSessions(t testing.TB, n int) []*sqlx.DB {
	testutils.SkipShortDB(t)

	ctx := testutils.Context(t)
	dbURL := testdb.New(t, true)

	sessions := make([]*sqlx.DB, n)
	for i := range sessions {
		db, err := pg.DBConfig{
			IdleInTxSessionTimeout: time.Hour,
			// Fail instead of hanging when a query waits on a lock, so tests that
			// check for lock avoidance regress loudly rather than timing out.
			LockTimeout:  10 * time.Second,
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		}.New(ctx, dbURL.String(), pg.DriverPostgres)
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })
		sessions[i] = db
	}

	return sessions
}

func MustCount(t *testing.T, db *sqlx.DB, stmt string, args ...any) (cnt int) {
	require.NoError(t, db.Get(&cnt, stmt, args...))
	return
}
