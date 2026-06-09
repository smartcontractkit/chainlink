package pgtest

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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
	ctx := testutils.Context(t)
	require.NoError(t, utils.JustError(ds.ExecContext(ctx, stmt, args...)))
}

func MustCount(t *testing.T, db *sqlx.DB, stmt string, args ...any) (cnt int) {
	require.NoError(t, db.Get(&cnt, stmt, args...))
	return
}
