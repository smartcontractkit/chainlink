package db

import (
	"testing"

	evmdb "github.com/smartcontractkit/chainlink/v2/core/chains/evm/db"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func NewScopedDB(t testing.TB, cfg config.Database) *evmdb.ScopedDB {
	evmURL := evmdb.ScopedConnection(cfg.URL())
	return &evmdb.ScopedDB{pgtest.NewSqlxDB(t, pg.WithURL(evmURL))}
}
