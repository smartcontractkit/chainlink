package migration1591204744

import (
	"github.com/jinzhu/gorm"
)

// Migrate ensures that heads are unique and adds parent hash for use in reorg detection
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE TABLE "job_spec_errors" (
		"id" bigserial primary key NOT NULL,
		"job_id" uuid REFERENCES job_specs(id) ON DELETE CASCADE NOT NULL,
		"description" text NOT NULL,
		"created_at" timestamp without time zone NOT NULL
	);
	`).Error
}
