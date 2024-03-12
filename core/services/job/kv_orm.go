package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// KVStore is a simple KV store that can store and retrieve serializable data.
//
//go:generate mockery --quiet --name KVStore --output ./mocks/ --case=underscore
type KVStore interface {
	Store(key string, val interface{}) error
	Get(key string, dest interface{}) error
}

type kVStore struct {
	jobID int32
	q     pg.Q
	lggr  logger.SugaredLogger
}

var _ KVStore = (*kVStore)(nil)

func NewKVStore(jobID int32, db *sqlx.DB, cfg pg.QConfig, lggr logger.Logger) kVStore {
	namedLogger := logger.Sugared(lggr.Named("JobORM"))
	return kVStore{
		jobID: jobID,
		q:     pg.NewQ(db, namedLogger, cfg),
		lggr:  namedLogger,
	}
}

// Store saves serializable value by key.
func (kv kVStore) Store(key string, val interface{}) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	sql := `INSERT INTO job_kv_store (job_id, key, val)
       	 	VALUES ($1, $2, $3)
        	ON CONFLICT (job_id, key) DO UPDATE SET
				val = EXCLUDED.val,
				updated_at = $4;`

	if err = kv.q.ExecQ(sql, kv.jobID, key, types.JSONText(jsonVal), time.Now()); err != nil {
		return fmt.Errorf("failed to store value: %s for key: %s for jobID: %d : %w", string(jsonVal), key, kv.jobID, err)
	}
	return nil
}

// Get retrieves serializable value by key.
func (kv kVStore) Get(key string, dest interface{}) error {
	var ret json.RawMessage
	sql := "SELECT val FROM job_kv_store WHERE job_id = $1 AND key = $2"
	if err := kv.q.Get(&ret, sql, kv.jobID, key); err != nil {
		return fmt.Errorf("failed to get value by key: %s for jobID: %d : %w", key, kv.jobID, err)
	}

	return json.Unmarshal(ret, dest)
}
