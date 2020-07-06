package migration1592355365

import (
	"github.com/jinzhu/gorm"
)

// Migrate ensures that heads are unique and adds parent hash for use in reorg detection
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE "job_spec_errors" (
			"id" bigserial primary key,
			"job_spec_id" uuid REFERENCES job_specs(id) ON DELETE CASCADE NOT NULL,
			"description" text NOT NULL,
			"occurrences" integer DEFAULT 1 NOT NULL,
			"created_at" timestamptz NOT NULL,
			"updated_at" timestamptz NOT NULL
		);

		CREATE UNIQUE INDEX job_spec_errors_unique_idx ON job_spec_errors ("job_spec_id", "description");
		CREATE INDEX job_spec_errors_created_at_idx ON job_spec_errors USING brin (created_at);
		CREATE INDEX job_spec_errors_updated_at_idx ON job_spec_errors USING brin (updated_at);
		CREATE INDEX job_spec_errors_occurrences_idx ON job_spec_errors (occurrences);
	`).Error
}
