package migrations

import (
	"gorm.io/gorm"
)

const up56 = `
ALTER TABLE eth_txes ADD COLUMN pipeline_task_run_id uuid UNIQUE;
ALTER TABLE eth_txes ADD COLUMN min_confirmations integer;
`

const down56 = `
DROP INDEX eth_txes_pipeline_task_run_id_idx;
ALTER TABLE eth_txes DROP COLUMN pipeline_task_run_id;
ALTER TABLE eth_txes DROP COLUMN min_confirmations;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0056_add_pipeline_task_runs_id_idx_to_eth_txs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up56).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down56).Error
		},
	})
}
