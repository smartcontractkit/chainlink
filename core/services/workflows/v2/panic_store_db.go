package v2

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// dbModulePanicStore is a durable ModulePanicStore backed by a DataSource. The
// counter must survive process restarts — the loop guard exists to break a
// crash-restart loop and the process is destroyed between panics — so it lives
// in Postgres rather than memory.
type dbModulePanicStore struct {
	ds sqlutil.DataSource
}

// NewDBModulePanicStore returns a durable ModulePanicStore backed by ds.
func NewDBModulePanicStore(ds sqlutil.DataSource) ModulePanicStore {
	return &dbModulePanicStore{ds: ds}
}

func (s *dbModulePanicStore) Store(ctx context.Context, key string, val []byte) error {
	_, err := s.ds.ExecContext(ctx,
		`INSERT INTO cre.workflow_module_panics (key, val, updated_at) VALUES ($1, $2, now())
		 ON CONFLICT (key) DO UPDATE SET val = EXCLUDED.val, updated_at = now()`, key, val)
	if err != nil {
		return fmt.Errorf("store workflow module panic %q: %w", key, err)
	}
	return nil
}

func (s *dbModulePanicStore) Get(ctx context.Context, key string) ([]byte, error) {
	var val []byte
	// Returns sql.ErrNoRows for a missing key; callers treat that as count zero.
	if err := s.ds.GetContext(ctx, &val, `SELECT val FROM cre.workflow_module_panics WHERE key = $1`, key); err != nil {
		return nil, err
	}
	return val, nil
}

func (s *dbModulePanicStore) Delete(ctx context.Context, key string) error {
	if _, err := s.ds.ExecContext(ctx, `DELETE FROM cre.workflow_module_panics WHERE key = $1`, key); err != nil {
		return fmt.Errorf("delete workflow module panic %q: %w", key, err)
	}
	return nil
}
