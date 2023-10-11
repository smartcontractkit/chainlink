package db

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type ScopedDB struct {
	*sqlx.DB
}

var schema = "evm"

func NewScopedDB(uuid uuid.UUID, cfg config.Database) (*ScopedDB, error) {
	db, err := pg.OpenUnlockedDB(uuid, cfg, pg.WithSchema(schema))
	if err != nil {
		return nil, err
	}
	return &ScopedDB{DB: db}, nil
}

func (s *ScopedDB) SqlxDB() *sqlx.DB {
	return s.DB
}

func ScopedConnection(dbURL url.URL) (evmScopedConnection url.URL) {
	// hacking, include public schema
	return pg.SchemaScopedConnection(dbURL, schema, "public")
}

func UseEVMSchema() pg.ConnectionOpt {
	return func(u *url.URL) error {
		return pg.WithSchema(schema)(u)
	}
}
