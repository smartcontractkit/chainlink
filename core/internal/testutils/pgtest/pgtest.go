package pgtest

import (
	"database/sql"
	"fmt"
	"math/rand"
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

	require.NoError(t, randomiseAndRestartSequences(db))
	return db
}

func NewSqlxDB(t testing.TB) *sqlx.DB {
	testutils.SkipShortDB(t)
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), uuid.New().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	db.MapperFunc(reflectx.CamelToSnakeASCII)

	require.NoError(t, randomiseAndRestartSequences(db))
	return db
}

type failedToSetupTestPGError struct{}

func (m *failedToSetupTestPGError) Error() string {
	return "failed to setup pgtest"
}

type sdlDBCommon interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

// randomiseAndRestartSequences randomizes sequenced table columns sequence
// This is necessary as to avoid false positives in some test cases.
func randomiseAndRestartSequences(db sdlDBCommon) error {
	seqRows, err := db.Query(`SELECT sequence_schema, sequence_name FROM information_schema.sequences WHERE sequence_schema = $1`, "public")
	if err != nil {
		return fmt.Errorf("failed to setup pgtest, error fetching sequences: %s", err)
	}

	defer seqRows.Close()
	for seqRows.Next() {
		var sequenceSchema, sequenceName string
		if err = seqRows.Scan(&sequenceSchema, &sequenceName); err != nil {
			return fmt.Errorf("%s: failed to setup pgtest, failed scanning sequence rows: %s", failedToSetupTestPGError{}, err)
		}
		if _, err = db.Exec(fmt.Sprintf("ALTER SEQUENCE %s.%s RESTART WITH %d", sequenceSchema, sequenceName, rand.Intn(10000))); err != nil {
			return fmt.Errorf("%s: failed to setup pgtest, failed to alter and restart %s sequence: %w", failedToSetupTestPGError{}, sequenceName, err)
		}
	}

	if err = seqRows.Err(); err != nil {
		return fmt.Errorf("%s: failed to setup pgtest, failed to iterate through sequences: %w", failedToSetupTestPGError{}, err)
	}

	return nil
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
