package store

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"

	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
)

// insertFixtures must execute the given fixture SQL without depending on
// on-disk paths derived from runtime.Caller, so that pre-compiled binaries
// (e.g. preparetest) work regardless of checkout location.
//
//nolint:paralleltest // provisions and mutates test database
func TestInsertFixtures(t *testing.T) {
	dbURL := testdb.New(t, true)

	require.NoError(t, insertFixtures(*dbURL, fixturesSQL))

	db, err := sqlx.Open(pgcommon.DriverPostgres, dbURL.String())
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	var count int
	require.NoError(t, db.Get(&count, "SELECT count(*) FROM users"))
	require.Positive(t, count, "expected fixture users to be inserted")
}

//nolint:paralleltest // provisions and mutates test database
func TestInsertFixturesUserOnly(t *testing.T) {
	dbURL := testdb.New(t, true)

	require.NoError(t, insertFixtures(*dbURL, usersOnlyFixtureSQL))

	db, err := sqlx.Open(pgcommon.DriverPostgres, dbURL.String())
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	var count int
	require.NoError(t, db.Get(&count, "SELECT count(*) FROM users"))
	require.Positive(t, count, "expected fixture users to be inserted")
}
