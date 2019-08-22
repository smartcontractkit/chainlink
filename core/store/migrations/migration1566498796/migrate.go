package migration1566498796

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

// Migrate optimizes the JobRuns table to reduce the cost of IDs
func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		return tx.Exec(`
ALTER TABLE run_results DROP COLUMN IF EXISTS cached_job_run_id;
ALTER TABLE run_results DROP COLUMN IF EXISTS cached_task_run_id;
		`).Error
	}

	return tx.Exec(`
CREATE TABLE "run_results_new" (
	"id" integer primary key autoincrement,
	"data" text,
	"status" varchar(255),
	"error_message" varchar(255),
	"amount" varchar(255)
);
INSERT INTO "run_results_new" SELECT "id", "data", "status", "error_message", "amount" FROM "run_results";
DROP TABLE "run_results";
ALTER TABLE "run_results_new" RENAME TO "run_results";
	`).Error
}
