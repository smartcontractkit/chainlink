package migrations

import (
	"gorm.io/gorm"
)

const up55 = `
CREATE INDEX eth_txes_pipeline_task_run_id_idx ON eth_txes ((meta ->> 'PipelineTaskRunID')) WHERE state = 'confirmed';
`

const down55 = `
DROP INDEX  eth_txes_pipeline_task_run_id_idx;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0055_add_pipeline_task_runs_id_idx_to_eth_txs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up55).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down55).Error
		},
	})
}
