package pgtest

import (
	"database/sql"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewQConfig(logSQL bool) pg.QConfig {
	return pg.NewQConfig(logSQL)
}

func NewSqlDB(t *testing.T) *sql.DB {
	testutils.SkipShortDB(t)
	db, err := sql.Open(string(dialects.TransactionWrappedPostgres), uuid.New().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	return db
}

func uniqueConnection(t testing.TB) *url.URL {
	url := testutils.MustParseURL(t, defaultDBURL.String())
	// inject uuid by default because the transaction wrapped driver requires it
	q := url.Query()
	q.Add("uuid", uuid.New().String())
	url.RawQuery = q.Encode()
	return url
}

func NewSqlxDB(t testing.TB, opts ...pg.ConnectionOpt) *sqlx.DB {
	testutils.SkipShortDB(t)

	url := uniqueConnection(t)
	for _, opt := range opts {
		assert.NoError(t, opt(url))
	}
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), url.String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}

func MustExec(t *testing.T, db *sqlx.DB, stmt string, args ...interface{}) {
	require.NoError(t, utils.JustError(db.Exec(stmt, args...)))
}

func MustSelect(t *testing.T, db *sqlx.DB, dest interface{}, stmt string, args ...interface{}) {
	require.NoError(t, db.Select(dest, stmt, args...))
}

func MustCount(t *testing.T, db *sqlx.DB, stmt string, args ...interface{}) (cnt int) {
	require.NoError(t, db.Get(&cnt, stmt, args...))
	return
}
