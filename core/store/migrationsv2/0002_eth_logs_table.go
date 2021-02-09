package migrationsv2

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

const up2 = `
    CREATE TABLE "eth_logs" (
        "block_hash" bytea NOT NULL,
        "block_number" bigint NOT NULL,
        "index" bigint NOT NULL,
        "address" bytea NOT NULL,
        "topics" bytea[] NOT NULL,
        "data" bytea NOT NULL,
        "order_received" serial NOT NULL,
        "created_at" timestamp without time zone NOT NULL,
        PRIMARY KEY (block_hash, index)
    );

    CREATE INDEX IF NOT EXISTS idx_eth_logs_address_block_number ON eth_logs (address, block_number);

    ALTER TABLE log_consumptions RENAME CONSTRAINT chk_log_consumptions_exactly_one_job_id TO chk_log_broadcasts_exactly_one_job_id;
    ALTER TABLE log_consumptions RENAME CONSTRAINT log_consumptions_job_id_fkey TO log_broadcasts_job_id_fkey;
    ALTER TABLE log_consumptions RENAME TO log_broadcasts;
    CREATE INDEX IF NOT EXISTS idx_log_broadcasts_blockhash_logindex_jobid_jobidv2 ON log_broadcasts (block_hash, log_index, job_id, job_id_v2);

    ALTER TABLE log_broadcasts ADD CONSTRAINT "log_broadcasts_eth_logs_fkey"
        FOREIGN KEY (block_hash, log_index) REFERENCES eth_logs (block_hash, index)
        ON DELETE CASCADE;

    ALTER TABLE log_broadcasts ADD COLUMN "consumed" BOOL NOT NULL DEFAULT FALSE;
`

const down2 = `
    -- ALTER TABLE log_broadcasts DROP COLUMN "consumed";
    ALTER TABLE log_broadcasts DROP CONSTRAINT "log_broadcasts_eth_logs_fkey";
    -- ALTER TABLE log_broadcasts RENAME CONSTRAINT log_broadcasts_job_id_fkey TO log_consumptions_job_id_fkey;
    -- ALTER TABLE log_broadcasts RENAME CONSTRAINT chk_log_broadcasts_exactly_one_job_id TO chk_log_consumptions_exactly_one_job_id;
    ALTER TABLE log_broadcasts RENAME TO log_consumptions;

    DROP TABLE eth_logs;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0002_eth_logs_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up2).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down2).Error
		},
	})
}
