package pgtest

import (
	"database/sql"
	"encoding/json"
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

func NewEVMScopedDB(t testing.TB) *sqlx.DB {
	// hack to scope to evm schema. the value "evm" will need to be dynamic to support multiple relayers
	url := pg.SchemaScopedConnection(defaultDBURL, "evm")
	return NewSqlxDB(t, WithURL(url))
}

func NewSqlxDB(t testing.TB, opts ...ConnectionOpt) *sqlx.DB {
	testutils.SkipShortDB(t)
	conn := &pg.ConnectionScope{
		UUID: uuid.New().String(),
		URL:  defaultDBURL,
	}
	for _, opt := range opts {
		opt(conn)
	}
	enc, err := json.Marshal(conn)
	require.NoError(t, err)
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), string(enc))
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}

type ConnectionOpt func(conn *pg.ConnectionScope)

func WithURL(url string) ConnectionOpt {
	return func(conn *pg.ConnectionScope) {
		conn.URL = url
	}
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
