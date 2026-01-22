-- +goose Up
-- Add verifier_id column as VARCHAR(255) to store verifier ID
ALTER TABLE ccv_chain_statuses ADD COLUMN verifier_id VARCHAR(255);

-- For existing rows, set a default empty string value (this should be temporary)
UPDATE ccv_chain_statuses SET verifier_id = '' WHERE verifier_id IS NULL;

-- Make verifier_id NOT NULL
ALTER TABLE ccv_chain_statuses ALTER COLUMN verifier_id SET NOT NULL;

-- Drop the old primary key
ALTER TABLE ccv_chain_statuses DROP CONSTRAINT ccv_chain_statuses_pkey;

-- Create new composite primary key
ALTER TABLE ccv_chain_statuses ADD PRIMARY KEY (chain_selector, verifier_id);

-- +goose Down
-- Restore old primary key (drop composite key first)
ALTER TABLE ccv_chain_statuses DROP CONSTRAINT ccv_chain_statuses_pkey;

-- Add back the old primary key
ALTER TABLE ccv_chain_statuses ADD PRIMARY KEY (chain_selector);

-- Drop verifier_id column
ALTER TABLE ccv_chain_statuses DROP COLUMN verifier_id;
