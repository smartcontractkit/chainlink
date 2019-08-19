package migration1560881855

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&LinkEarned{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate link_earned")
	}

	var fillLinkEarned string
	if dbutil.IsPostgres(tx) {
		fillLinkEarned = `
		INSERT INTO link_earned
		SELECT ROW_NUMBER() OVER (ORDER BY job_spec_id) AS id, job_spec_id, jr.id AS job_run_id, amount, finished_at
		FROM job_runs jr INNER JOIN run_results rr ON jr.overrides_id  = rr.id
		WHERE amount IS NOT NULL
		`
	} else {
		fillLinkEarned = `
		INSERT INTO link_earned
		SELECT ROW_NUMBER() OVER (ORDER BY job_spec_id) AS id, job_spec_id, job_runs.id AS job_run_id, amount, finished_at
		FROM job_runs INNER JOIN run_results ON job_runs.overrides_id  = run_results.id
		WHERE amount IS NOT NULL
		`
	}
	if err := tx.Exec(fillLinkEarned).Error; err != nil {
		return errors.Wrap(err, "failed to fill existing run rewards to link_earned table")
	}
	return nil
}

// LinkEarned is a capture of the model before migration1565291711
type LinkEarned struct {
	ID        uint64       `gorm:"primary_key;not null;auto_increment"`
	JobSpecID string       `gorm:"index;not null;type:varchar(36) REFERENCES job_specs(id)"`
	JobRunID  string       `gorm:"unique;not null;type:varchar(36) REFERENCES job_runs(id)"`
	Earned    *assets.Link `gorm:"type:varchar(255)"`
	EarnedAt  time.Time    `gorm:"index"`
}
