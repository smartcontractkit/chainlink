package migrations

import (
	"gorm.io/gorm"
)

const up55 = `
ALTER TABLE eth_txes ADD COLUMN pipeline_task_run_id uuid REFERENCES pipeline_task_runs(id);
ALTER TABLE eth_txes ADD COLUMN min_confirmations integer;
`

const down55 = `
ALTER TABLE eth_txes DROP COLUMN pipeline_task_run_id;
ALTER TABLE eth_txes DROP COLUMN min_confirmations;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0055_add_pipeline_task_runs_id_to_eth_txs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up55).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down55).Error
		},
	})
}
