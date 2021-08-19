package migrations

import (
	"gorm.io/gorm"
)

const up57 = `
ALTER TABLE eth_txes ADD COLUMN pipeline_task_run_id uuid UNIQUE;
ALTER TABLE eth_txes ADD COLUMN min_confirmations integer;
CREATE INDEX pipeline_runs_suspended ON pipeline_runs (id) WHERE state = 'suspended' ;
`

const down57 = `
DROP INDEX eth_txes_pipeline_task_run_id_idx;
ALTER TABLE eth_txes DROP COLUMN pipeline_task_run_id;
ALTER TABLE eth_txes DROP COLUMN min_confirmations;
DROP INDEX pipeline_runs_suspended;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0057_add_pipeline_task_runs_id_idx_to_eth_txs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up57).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down57).Error
		},
	})
}
