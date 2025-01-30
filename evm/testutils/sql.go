package testutils

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
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
