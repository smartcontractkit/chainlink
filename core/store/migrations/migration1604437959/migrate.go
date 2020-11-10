package migration1604437959

import "github.com/jinzhu/gorm"

const up = `
	CREATE TABLE "job_spec_errors_v2" (
		"id" bigserial primary key,
		"job_id" int REFERENCES jobs(id) ON DELETE CASCADE,
		"description" text NOT NULL,
		"occurrences" integer DEFAULT 1 NOT NULL,
		"created_at" timestamptz NOT NULL,
		"updated_at" timestamptz NOT NULL
	);
	CREATE UNIQUE INDEX job_spec_errors_v2_unique_idx ON job_spec_errors_v2 ("job_id", "description");
`

const down = `
	DROP TABLE "job_spec_errors_v2"
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
