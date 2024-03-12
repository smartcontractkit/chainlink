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

//go:generate mockery --quiet --name KVStore --output ./mocks/ --case=underscore

type KVStore interface {
	Store(key string, val interface{}) error
	Get(key string, dest interface{}) error
}

type jobKVStore struct {
	jobID int32
	q     pg.Q
	lggr  logger.SugaredLogger
}

var _ KVStore = (*jobKVStore)(nil)

func NewJobKVStore(jobID int32, db *sqlx.DB, cfg pg.QConfig, lggr logger.Logger) KVStore {
	namedLogger := logger.Sugared(lggr.Named("JobORM"))
	return &jobKVStore{
		jobID: jobID,
		q:     pg.NewQ(db, namedLogger, cfg),
		lggr:  namedLogger,
	}
}

func (kv jobKVStore) Store(key string, val interface{}) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	sql := `INSERT INTO job_kv_store (id, key, val)
       	 	VALUES ($1, $2, $3)
        	ON CONFLICT (id, key) DO UPDATE SET
				val = EXCLUDED.val,
				updated_at = $4
        	RETURNING id;`

	if err = kv.q.ExecQ(sql, kv.jobID, key, types.JSONText(jsonVal), time.Now()); err != nil {
		return fmt.Errorf("failed to store value: %s for key: %s for jobID: %d : %w", string(jsonVal), key, kv.jobID, err)
	}
	return nil
}

func (kv jobKVStore) Get(key string, dest interface{}) error {
	var ret json.RawMessage
	sql := "SELECT val FROM job_kv_store WHERE id = $1 AND key = $2"
	if err := kv.q.Get(&ret, sql, kv.jobID, key); err != nil {
		return fmt.Errorf("failed to get value by key: %s for jobID: %d : %w", key, kv.jobID, err)
	}
	
	return json.Unmarshal(ret, dest)
}
