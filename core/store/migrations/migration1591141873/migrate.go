package migration1591141873

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`
        CREATE TABLE flux_monitor_round_stats (
            id bigserial primary key,
            aggregator bytea not null,
            round_id integer not null,
            num_new_round_logs integer not null default 0,
            num_submissions integer not null default 0,

            UNIQUE (aggregator, round_id)
        )
    `).Error
	if err != nil {
		return errors.Wrap(err, "failed to create FluxMonitorRoundStats table")
	}
	return nil
}
