package pg

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

var _ Getter = &mockGetter{}

type mockGetter struct {
	version int
	err     error
}

func (m *mockGetter) Get(dest interface{}, query string, args ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	*(dest.(*int)) = m.version
	return nil
}

func Test_checkVersion(t *testing.T) {
	if time.Now().Year() > 2027 {
		t.Fatal("Postgres version numbers only registered until 2028, please update the postgres version check using: https://www.postgresql.org/support/versioning/ then fix this test")
	}
	t.Run("when the version is too low", func(t *testing.T) {
		m := &mockGetter{version: 100000}
		err := checkVersion(m, 110000)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "The minimum required Postgres server version is 11, you are running: 10")
	})
	t.Run("when the version is at minimum", func(t *testing.T) {
		m := &mockGetter{version: 110000}
		err := checkVersion(m, 110000)
		require.NoError(t, err)
	})
	t.Run("when the version is above minimum", func(t *testing.T) {
		m := &mockGetter{version: 110001}
		err := checkVersion(m, 110000)
		require.NoError(t, err)
		m = &mockGetter{version: 120000}
		err = checkVersion(m, 110001)
		require.NoError(t, err)
	})
	t.Run("ignores wildly small versions, 0 etc", func(t *testing.T) {
		m := &mockGetter{version: 9000}
		err := checkVersion(m, 110001)
		require.NoError(t, err)
	})
	t.Run("ignores errors", func(t *testing.T) {
		m := &mockGetter{err: errors.New("some error")}
		err := checkVersion(m, 110001)
		require.NoError(t, err)
	})
}

func Test_disallowReplica(t *testing.T) {
	testutils.SkipShortDB(t)
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), uuid.New().String())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err = db.Exec("SET session_replication_role= 'origin'")
	require.NoError(t, err)
	err = disallowReplica(db)
	require.NoError(t, err)

	_, err = db.Exec("SET session_replication_role= 'replica'")
	require.NoError(t, err)
	err = disallowReplica(db)
	require.Error(t, err, "replica role should be disallowed")

	_, err = db.Exec("SET session_replication_role= 'not_valid_role'")
	require.Error(t, err)
}
