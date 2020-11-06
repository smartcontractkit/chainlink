package migration1604437959

import "github.com/jinzhu/gorm"

const up = `
	CREATE TABLE "pipeline_spec_errors" (
		"id" bigserial primary key,
		"pipeline_spec_id" int REFERENCES pipeline_specs(id) ON DELETE CASCADE,
		"description" text NOT NULL,
		"occurrences" integer DEFAULT 1 NOT NULL,
		"created_at" timestamptz NOT NULL,
		"updated_at" timestamptz NOT NULL
	);
	CREATE UNIQUE INDEX pipeline_spec_errors_unique_idx ON pipeline_spec_errors ("pipeline_spec_id", "description");
`

const down = `
	DROP TABLE "pipeline_spec_errors"
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
