package pgtest

import (
	"database/sql"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ pg.QConfig = &qConfig{}

// qConfig implements pg.QCOnfig
type qConfig struct {
	logSQL              bool
	defaultQueryTimeout time.Duration
}

func NewQConfig(logSQL bool) pg.QConfig {
	return &qConfig{logSQL, pg.DefaultQueryTimeout}
}

func (p *qConfig) LogSQL() bool { return p.logSQL }

func (p *qConfig) DatabaseDefaultQueryTimeout() time.Duration { return p.defaultQueryTimeout }

func NewSqlDB(t *testing.T) *sql.DB {
	testutils.SkipShortDB(t)
	db, err := sql.Open(string(dialects.TransactionWrappedPostgres), uuid.NewV4().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	return db
}

func NewSqlxDB(t testing.TB) *sqlx.DB {
	testutils.SkipShortDB(t)
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), uuid.NewV4().String())
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
