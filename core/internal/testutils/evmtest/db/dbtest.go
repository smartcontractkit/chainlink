package db

import (
	"testing"

	evmdb "github.com/smartcontractkit/chainlink/v2/core/chains/evm/db"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func NewScopedDB(t testing.TB, cfg config.Database) *evmdb.ScopedDB {
	evmURL := evmdb.ScopedConnection(cfg.URL())
	return &evmdb.ScopedDB{pgtest.NewSqlxDB(t, pg.WithURL(evmURL))}
}

// krr note to self: consolidate with above
func NewScopedDBWithOpts(t testing.TB, opts ...TestDBOpt) *evmdb.ScopedDB {
	cfg := configtest.NewTestGeneralConfig(t)
	evmURL := evmdb.ScopedConnection(cfg.Database().URL())
	d := &evmdb.ScopedDB{DB: pgtest.NewSqlxDB(t, pg.WithURL(evmURL))}
	for _, opt := range opts {
		opt(t, d)
	}
	return d
}

type TestDBOpt = func(t testing.TB, d *evmdb.ScopedDB)

func DBConfig(cfg config.Database) TestDBOpt {
	return func(t testing.TB, d *evmdb.ScopedDB) {
		evmURL := evmdb.ScopedConnection(cfg.URL())
		d.DB = pgtest.NewSqlxDB(t, pg.WithURL(evmURL))
	}
}
