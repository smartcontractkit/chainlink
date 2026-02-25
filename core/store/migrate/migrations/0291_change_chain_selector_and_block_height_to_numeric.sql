-- +goose Up
-- Change chain_selector and finalized_block_height from TEXT to NUMERIC(20, 0)
-- These columns store uint64 values which fit within NUMERIC(20, 0)

-- Drop the primary key constraint temporarily since chain_selector is part of it
ALTER TABLE ccv_chain_statuses DROP CONSTRAINT ccv_chain_statuses_pkey;

-- Change column types
-- The USING clause will cast existing TEXT values to NUMERIC
-- Since the code stores numeric strings (e.g., "123456"), the cast will work correctly
ALTER TABLE ccv_chain_statuses ALTER COLUMN chain_selector TYPE NUMERIC(20, 0) USING chain_selector::NUMERIC(20, 0);
ALTER TABLE ccv_chain_statuses ALTER COLUMN finalized_block_height TYPE NUMERIC(20, 0) USING finalized_block_height::NUMERIC(20, 0);

-- Recreate the primary key constraint
ALTER TABLE ccv_chain_statuses ADD PRIMARY KEY (chain_selector, verifier_id);

-- +goose Down
-- Revert back to TEXT type
ALTER TABLE ccv_chain_statuses DROP CONSTRAINT ccv_chain_statuses_pkey;
ALTER TABLE ccv_chain_statuses ALTER COLUMN chain_selector TYPE TEXT USING chain_selector::TEXT;
ALTER TABLE ccv_chain_statuses ALTER COLUMN finalized_block_height TYPE TEXT USING finalized_block_height::TEXT;
ALTER TABLE ccv_chain_statuses ADD PRIMARY KEY (chain_selector, verifier_id);

