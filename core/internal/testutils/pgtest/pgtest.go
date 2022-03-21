package pgtest

import (
	"database/sql"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ pg.LogConfig = PGCfg{}

type PGCfg struct{ logSQL bool }

func NewPGCfg(logSQL bool) pg.LogConfig { return PGCfg{logSQL} }
func (p PGCfg) LogSQL() bool            { return p.logSQL }

func NewSqlDB(t *testing.T) *sql.DB {
	testutils.SkipShortDB(t)
	db, err := sql.Open("txdb", uuid.NewV4().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	return db
}

func NewSqlxDB(t *testing.T) *sqlx.DB {
	testutils.SkipShortDB(t)
	db, err := sqlx.Open("txdb", uuid.NewV4().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	db.MapperFunc(reflectx.CamelToSnakeASCII)

	return db
}

func MustExec(t *testing.T, db *sqlx.DB, stmt string, args ...interface{}) {
	require.NoError(t, utils.JustError(db.Exec(stmt, args...)))
}
