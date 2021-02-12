package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up3 = `
    CREATE TABLE "eth_logs" (
        "id" BIGSERIAL PRIMARY KEY,
        "block_hash" bytea NOT NULL,
        "block_number" bigint NOT NULL,
        "index" bigint NOT NULL,
        "address" bytea NOT NULL,
        "topics" bytea[] NOT NULL,
        "data" bytea NOT NULL,
        "order_received" serial NOT NULL,
        "created_at" timestamp without time zone NOT NULL
    );

    CREATE UNIQUE INDEX idx_eth_logs_unique ON eth_logs (block_hash, index) INCLUDE (id);
    CREATE INDEX IF NOT EXISTS idx_eth_logs_block_number ON eth_logs (block_number);
    CREATE INDEX IF NOT EXISTS idx_eth_logs_address_block_number ON eth_logs (address, block_number);

    -- TODO: IS THIS SAFE????????????????????????????????????????????
    DELETE FROM log_consumptions;

	ALTER TABLE log_consumptions RENAME CONSTRAINT chk_log_consumptions_exactly_one_job_id TO chk_log_broadcasts_exactly_one_job_id;
	ALTER TABLE log_consumptions RENAME CONSTRAINT log_consumptions_job_id_fkey TO log_broadcasts_job_id_fkey;
	ALTER TABLE log_consumptions RENAME TO log_broadcasts;

    ALTER TABLE log_broadcasts
        ADD COLUMN "consumed" BOOL NOT NULL DEFAULT FALSE,
        ADD COLUMN "eth_log_id" BIGINT NOT NULL,
		ADD CONSTRAINT log_broadcasts_eth_log_id_fkey FOREIGN KEY (eth_log_id) REFERENCES eth_logs (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

    CREATE INDEX idx_log_broadcasts_unconsumed_eth_log_id ON log_broadcasts (eth_log_id) WHERE consumed = false;
    CREATE INDEX idx_log_broadcasts_unconsumed_job_id ON log_broadcasts (job_id) WHERE consumed = false AND job_id IS NOT NULL;
    CREATE INDEX idx_log_broadcasts_unconsumed_job_id_v2 ON log_broadcasts (job_id_v2) WHERE consumed = false AND job_id_v2 IS NOT NULL;

	DROP INDEX IF EXISTS log_consumptions_unique_v1_idx;
	DROP INDEX IF EXISTS log_consumptions_unique_v2_idx;

	CREATE UNIQUE INDEX log_consumptions_unique_v1_idx ON log_broadcasts(job_id, block_hash, log_index) INCLUDE (consumed) WHERE job_id IS NOT NULL;
	CREATE UNIQUE INDEX log_consumptions_unique_v2_idx ON log_broadcasts(job_id_v2, block_hash, log_index) INCLUDE (consumed) WHERE job_id_v2 IS NOT NULL;
`

// TODO: Finalise down3
const down3 = `
    -- ALTER TABLE log_broadcasts DROP COLUMN "consumed";
    ALTER TABLE log_broadcasts DROP CONSTRAINT "log_broadcasts_eth_logs_fkey";
    -- ALTER TABLE log_broadcasts RENAME CONSTRAINT log_broadcasts_job_id_fkey TO log_consumptions_job_id_fkey;
    -- ALTER TABLE log_broadcasts RENAME CONSTRAINT chk_log_broadcasts_exactly_one_job_id TO chk_log_consumptions_exactly_one_job_id;
    ALTER TABLE log_broadcasts RENAME TO log_consumptions;

    DROP TABLE eth_logs;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0003_eth_logs_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up3).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down3).Error
		},
	})
}
