-- +goose Up

CREATE TABLE IF NOT EXISTS ccv_chain_statuses (
    chain_selector TEXT PRIMARY KEY,
    finalized_block_height TEXT NOT NULL,
    disabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ccv_chain_statuses_updated_at ON ccv_chain_statuses(updated_at);

-- +goose Down

DROP INDEX IF EXISTS idx_ccv_chain_statuses_updated_at;
DROP TABLE IF EXISTS ccv_chain_statuses;

