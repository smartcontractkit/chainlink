package migration1587580235

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the LogConsumption table and adds a uniqueness constraint on the
// combination of block_hash / log_index / job_id
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE TABLE "log_consumptions" (
		"id" bigserial primary key NOT NULL,
		"block_hash" bytea NOT NULL,
		"log_index" bigint NOT NULL,
		"job_id" uuid REFERENCES job_specs(id) ON DELETE CASCADE NOT NULL,
		"created_at" timestamp without time zone NOT NULL
	);

	CREATE UNIQUE INDEX log_consumptions_unique_idx ON log_consumptions ("job_id", "block_hash", "log_index");
	CREATE INDEX log_consumptions_created_at_idx ON log_consumptions USING brin (created_at);
	`).Error
}
