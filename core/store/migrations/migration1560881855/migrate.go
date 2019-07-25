package migration1560881855

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.LinkEarned{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate link_earned")
	}
	var fillLinkEarned string
	if dbutil.IsPostgres(tx) {
		fillLinkEarned = `
		INSERT INTO link_earned
		SELECT * FROM
		(
		SELECT ROW_NUMBER() OVER (ORDER BY job_spec_id) AS id, job_spec_id, amount,  finished_at
		FROM
		(
		job_runs
		INNER JOIN run_results ON job_runs.overrides_id  = run_results.id
		) AS postgres_syntax
		) AS another_postgres_syntax
		WHERE amount IS NOT NULL`
	} else {
		fillLinkEarned = `
		INSERT INTO link_earned
		SELECT * FROM
		(
		SELECT ROW_NUMBER() OVER (ORDER BY job_spec_id) AS id, job_spec_id, amount, finished_at
		FROM job_runs
		INNER JOIN run_results ON job_runs.overrides_id  = run_results.id
		WHERE amount IS NOT NULL
		)`
	}
	if err := tx.Exec(fillLinkEarned).Error; err != nil {
		return errors.Wrap(err, "failed to fill existing run rewards to link_earned table")
	}
	return nil
}
