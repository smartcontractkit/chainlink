package migrations

import (
	"gorm.io/gorm"
)

const up33 = `
CREATE INDEX flux_monitor_round_stats_job_run_id_idx ON flux_monitor_round_stats (job_run_id);
CREATE INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx ON flux_monitor_round_stats_v2 (pipeline_run_id);
`
const down33 = `
DROP INDEX flux_monitor_round_stats_job_run_id_idx;
DROP INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0033_flux_monitor_round_stats_fk_index",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up33).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down33).Error
		},
	})
}
