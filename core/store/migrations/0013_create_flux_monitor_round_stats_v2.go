package migrations

import (
	"gorm.io/gorm"
)

const (
	up13 = `
CREATE TABLE flux_monitor_round_stats_v2 (
	id BIGSERIAL PRIMARY KEY,
	aggregator bytea NOT NULL,
	round_id integer NOT NULL,
	num_new_round_logs integer NOT NULL DEFAULT 0,
	num_submissions integer NOT NULL DEFAULT 0,
	pipeline_run_id bigint REFERENCES pipeline_runs(id) ON DELETE CASCADE,
	CONSTRAINT flux_monitor_round_stats_v2_aggregator_round_id_key UNIQUE (aggregator, round_id)
);
`
	down13 = `
DROP TABLE flux_monitor_round_stats_v2;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0013_create_flux_monitor_round_stats_v2",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up13).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down13).Error
		},
	})
}
