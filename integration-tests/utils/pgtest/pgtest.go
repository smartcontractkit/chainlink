package pgtest

import (
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

func NewSqlxDB(t testing.TB) *sqlx.DB {
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		t.Errorf("you must provide a CL_DATABASE_URL environment variable")
		return nil
	}
	return sqltest.NewDB(t, dbURL)
}
