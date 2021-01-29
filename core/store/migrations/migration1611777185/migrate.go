package migration1611777185

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE "logs" (
            "block_hash" bytea NOT NULL,
            "block_number" bigint NOT NULL,
            "index" bigint NOT NULL,
            "address" bytea NOT NULL,
            "topics" bytea[] NOT NULL,
            "data" bytea NOT NULL,
            "created_at" timestamp without time zone NOT NULL,
            PRIMARY KEY (block_hash, index)
        );

        ALTER TABLE log_consumptions RENAME CONSTRAINT chk_log_consumptions_exactly_one_job_id TO chk_log_broadcasts_exactly_one_job_id;
        ALTER TABLE log_consumptions RENAME CONSTRAINT log_consumptions_job_id_fkey TO log_broadcasts_job_id_fkey;
        ALTER TABLE log_consumptions RENAME TO log_broadcasts;


        ALTER TABLE log_broadcasts ADD CONSTRAINT "log_broadcasts_logs_fkey"
            FOREIGN KEY (block_hash, log_index) REFERENCES logs (block_hash, index)
            ON DELETE CASCADE;

        ALTER TABLE log_broadcasts ADD COLUMN "consumed" BOOL NOT NULL DEFAULT FALSE;


    `).Error
}
