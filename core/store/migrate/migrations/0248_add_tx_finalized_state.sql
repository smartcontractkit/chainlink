-- +goose Up
-- Creating new column and enum instead of just adding new value to the existing enum so the migration changes match the rollback logic
-- Otherwise, migration will complain about mismatching column order

-- +goose StatementBegin
-- Rename the existing enum with finalized state to mark it as old
ALTER TYPE evm.txes_state RENAME TO txes_state_old;

-- Create new enum without finalized state
CREATE TYPE evm.txes_state AS ENUM (
    'unstarted',
    'in_progress',
    'fatal_error',
    'unconfirmed',
    'confirmed_missing_receipt',
    'confirmed',
    'finalized'
);

-- Add a new state column with the new enum type to the txes table
ALTER TABLE evm.txes ADD COLUMN state_new evm.txes_state;

-- Copy data from the old column to the new
UPDATE evm.txes SET state_new = state::text::evm.txes_state;

-- Drop constraints referring to old enum type on the old state column
ALTER TABLE evm.txes ALTER COLUMN state DROP DEFAULT;
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
DROP INDEX IF EXISTS idx_eth_txes_state_from_address_evm_chain_id;
DROP INDEX IF EXISTS idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id;
DROP INDEX IF EXISTS idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id;
DROP INDEX IF EXISTS idx_eth_txes_unstarted_subject_id_evm_chain_id;

-- Drop the old state column
ALTER TABLE evm.txes DROP state;

-- Drop the old enum type
DROP TYPE evm.txes_state_old;

-- Rename the new column name state to replace the old column
ALTER TABLE evm.txes RENAME state_new TO state;

-- Reset the state column's default
ALTER TABLE evm.txes ALTER COLUMN state SET DEFAULT 'unstarted'::evm.txes_state, ALTER COLUMN state SET NOT NULL;

-- Recreate constraint with finalized state
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
    state = 'unstarted'::evm.txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'in_progress'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'fatal_error'::evm.txes_state AND error IS NOT NULL
    OR
    state = 'unconfirmed'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed_missing_receipt'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'finalized'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
) NOT VALID;

-- Recreate index to include finalized state
CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON evm.txes(evm_chain_id, from_address, state) WHERE state <> 'confirmed'::evm.txes_state AND state <> 'finalized'::evm.txes_state;
CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON evm.txes(evm_chain_id, from_address, nonce) WHERE state = 'unconfirmed'::evm.txes_state;
CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON evm.txes(evm_chain_id, from_address) WHERE state = 'in_progress'::evm.txes_state;
CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON evm.txes(evm_chain_id, subject, id) WHERE subject IS NOT NULL AND state = 'unstarted'::evm.txes_state;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rename the existing enum with finalized state to mark it as old
ALTER TYPE evm.txes_state RENAME TO txes_state_old;

-- Create new enum without finalized state
CREATE TYPE evm.txes_state AS ENUM (
    'unstarted',
    'in_progress',
    'fatal_error',
    'unconfirmed',
    'confirmed_missing_receipt',
    'confirmed'
);

-- Add a new state column with the new enum type to the txes table
ALTER TABLE evm.txes ADD COLUMN state_new evm.txes_state;

-- Update all transactions with finalized state to confirmed in the old state column
UPDATE evm.txes SET state = 'confirmed'::evm.txes_state_old WHERE state = 'finalized'::evm.txes_state_old;

-- Copy data from the old column to the new
UPDATE evm.txes SET state_new = state::text::evm.txes_state;

-- Drop constraints referring to old enum type on the old state column
ALTER TABLE evm.txes ALTER COLUMN state DROP DEFAULT;
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
DROP INDEX IF EXISTS idx_eth_txes_state_from_address_evm_chain_id;
DROP INDEX IF EXISTS idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id;
DROP INDEX IF EXISTS idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id;
DROP INDEX IF EXISTS idx_eth_txes_unstarted_subject_id_evm_chain_id;

-- Drop the old state column
ALTER TABLE evm.txes DROP state;

-- Drop the old enum type
DROP TYPE evm.txes_state_old;

-- Rename the new column name state to replace the old column
ALTER TABLE evm.txes RENAME state_new TO state;

-- Reset the state column's default
ALTER TABLE evm.txes ALTER COLUMN state SET DEFAULT 'unstarted'::evm.txes_state, ALTER COLUMN state SET NOT NULL;

-- Recereate constraint without finalized state
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
    state = 'unstarted'::evm.txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'in_progress'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'fatal_error'::evm.txes_state AND error IS NOT NULL
    OR
    state = 'unconfirmed'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed_missing_receipt'::evm.txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
) NOT VALID;

-- Recreate index with new enum type
CREATE INDEX idx_eth_txes_state_from_address_evm_chain_id ON evm.txes(evm_chain_id, from_address, state) WHERE state <> 'confirmed'::evm.txes_state;
CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key_evm_chain_id ON evm.txes(evm_chain_id, from_address, nonce) WHERE state = 'unconfirmed'::evm.txes_state;
CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account_id_per_evm_chain_id ON evm.txes(evm_chain_id, from_address) WHERE state = 'in_progress'::evm.txes_state;
CREATE INDEX idx_eth_txes_unstarted_subject_id_evm_chain_id ON evm.txes(evm_chain_id, subject, id) WHERE subject IS NOT NULL AND state = 'unstarted'::evm.txes_state;
-- +goose StatementEnd
