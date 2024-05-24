-- +goose Up
    ALTER TABLE log_broadcasts DROP COLUMN "eth_log_id";
    DROP TABLE "eth_logs";

-- +goose Down
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

    -- NOTE: one-time deletion is necessary to maintain FK constraints, this probably won't hurt
    DELETE FROM log_broadcasts;

    ALTER TABLE log_broadcasts
        ADD COLUMN "eth_log_id" BIGINT NOT NULL,
        ADD CONSTRAINT log_broadcasts_eth_log_id_fkey FOREIGN KEY (eth_log_id) REFERENCES eth_logs (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

    CREATE INDEX idx_log_broadcasts_unconsumed_eth_log_id ON log_broadcasts (eth_log_id) WHERE consumed = false;
