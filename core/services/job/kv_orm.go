package job

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// KVStore is a simple KV store that can store and retrieve serializable data.
type KVStore interface {
	Store(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type kVStore struct {
	jobID int32
	ds    sqlutil.DataSource
}

var _ KVStore = (*kVStore)(nil)

func NewKVStore(jobID int32, ds sqlutil.DataSource) kVStore {
	return kVStore{
		jobID: jobID,
		ds:    ds,
	}
}

// Store saves []byte value by key.
func (kv kVStore) Store(ctx context.Context, key string, val []byte) error {
	sql := `INSERT INTO job_kv_store (job_id, key, val_bytea)
       	 	VALUES ($1, $2, $3)
        	ON CONFLICT (job_id, key) DO UPDATE SET
				val_bytea = EXCLUDED.val_bytea,
				updated_at = $4;`

	if _, err := kv.ds.ExecContext(ctx, sql, kv.jobID, key, val, time.Now()); err != nil {
		return fmt.Errorf("failed to store value: %s for key: %s for jobID: %d : %w", string(val), key, kv.jobID, err)
	}
	return nil
}

// Get retrieves []byte value by key.
func (kv kVStore) Get(ctx context.Context, key string) ([]byte, error) {
	var val []byte
	sql := "SELECT val_bytea FROM job_kv_store WHERE job_id = $1 AND key = $2"
	if err := kv.ds.GetContext(ctx, &val, sql, kv.jobID, key); err != nil {
		return nil, fmt.Errorf("failed to get value by key: %s for jobID: %d : %w", key, kv.jobID, err)
	}

	return val, nil
}
