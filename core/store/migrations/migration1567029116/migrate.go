package migration1567029116

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

// Migrate optimizes the JobRuns table to reduce the cost of IDs
func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		return tx.Exec(`
ALTER TABLE run_requests ADD COLUMN "payment" numeric(78, 0);
ALTER TABLE job_runs ADD COLUMN "payment" numeric(78, 0);
UPDATE job_runs
SET "payment" = CAST("earned" AS numeric)
FROM link_earned
WHERE job_run_id = job_runs.id;
DROP TABLE "link_earned";
ALTER TABLE run_results DROP COLUMN "amount";
`).Error
	}

	return tx.Exec(`
ALTER TABLE run_requests ADD COLUMN "payment" varchar(255);
ALTER TABLE job_runs ADD COLUMN "payment" varchar(255);
UPDATE job_runs
SET "payment" = (
	SELECT earned
	FROM link_earned
	WHERE job_run_id = job_runs.id
);
DROP TABLE "link_earned";
CREATE TABLE "run_results_new" (
	"id" integer primary key autoincrement,
	"data" text,
	"status" varchar(255),
	"error_message" varchar(255)
);
INSERT INTO "run_results_new" SELECT "id", "data", "status", "error_message" FROM "run_results";
DROP TABLE "run_results";
ALTER TABLE "run_results_new" RENAME TO "run_results";
	`).Error
}
