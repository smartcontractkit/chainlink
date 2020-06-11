package migration1589206996

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the requisite tables for the BulletproofTxManager
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TYPE eth_txes_state AS ENUM ('unstarted', 'in_progress', 'fatal_error', 'unconfirmed', 'confirmed');

	  	CREATE TABLE eth_txes (
			id BIGSERIAL PRIMARY KEY,
			nonce bigint, 
			from_address bytea NOT NULL REFERENCES keys (address),
			to_address bytea NOT NULL,
			encoded_payload bytea NOT NULL,
			value numeric(78, 0) NOT NULL,
			gas_limit bigint NOT NULL,
			error text,
			broadcast_at timestamptz,
			created_at timestamptz NOT NULL,
			state eth_txes_state NOT NULL DEFAULT 'unstarted'::eth_txes_state
		);
		  
		ALTER TABLE eth_txes ADD CONSTRAINT chk_from_address_length CHECK (
			octet_length(from_address) = 20
		);
		ALTER TABLE eth_txes ADD CONSTRAINT chk_to_address_length CHECK (
			octet_length(to_address) = 20
		);

		ALTER TABLE eth_txes ADD CONSTRAINT chk_broadcast_at_is_sane CHECK (
			broadcast_at > '2019-01-01'
		);

		ALTER TABLE eth_txes ADD CONSTRAINT chk_error_cannot_be_empty CHECK (
			error IS NULL OR length(error) > 0
		);

		CREATE UNIQUE INDEX idx_eth_txes_nonce_from_address ON eth_txes (nonce, from_address);
		CREATE INDEX idx_eth_txes_state ON eth_txes (state) WHERE state != 'confirmed'::eth_txes_state;
		CREATE INDEX idx_eth_txes_broadcast_at ON eth_txes USING BRIN (broadcast_at);
		CREATE INDEX idx_eth_txes_created_at ON eth_txes USING BRIN (created_at);

		-- Only one in progress transaction allowed per account
		CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account ON eth_txes (from_address) WHERE state = 'in_progress';

		ALTER TABLE eth_txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
			state = 'unstarted' AND (nonce is NULL AND error IS NULL AND broadcast_at IS NULL) 
				OR
			state = 'in_progress' AND (nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL)
				OR
			state = 'fatal_error' AND (nonce is NULL AND error IS NOT NULL AND broadcast_at IS NULL)
				OR
			state = 'unconfirmed' AND (nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL)
				OR
			state = 'confirmed' AND (nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL)
		);

		CREATE TYPE eth_tx_attempts_state AS ENUM ('in_progress', 'broadcast');

		CREATE TABLE eth_tx_attempts (
			id BIGSERIAL PRIMARY KEY,
		 	eth_tx_id bigint NOT NULL REFERENCES eth_txes (id) ON DELETE CASCADE,
		 	gas_price numeric(78,0) NOT NULL,
		 	signed_raw_tx bytea NOT NULL,
		 	hash bytea NOT NULL,
			broadcast_before_block_num bigint,
			state eth_tx_attempts_state NOT NULL,
		 	created_at timestamptz NOT NULL
		);

		ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_eth_tx_attempts_fsm CHECK (
			(state = 'in_progress' AND broadcast_before_block_num IS NULL) OR state = 'broadcast'
		);

		ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_signed_raw_tx_present CHECK (
			octet_length(signed_raw_tx) > 0
		);

		ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_hash_length CHECK (
			octet_length(hash) = 32
		);

		ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_cannot_broadcast_before_block_zero CHECK (
			broadcast_before_block_num IS NULL OR broadcast_before_block_num > 0
		);
		
		-- Should never have more than one in_progress attempt per nonce
		CREATE UNIQUE INDEX idx_only_one_in_progress_attempt_per_eth_tx ON eth_tx_attempts (eth_tx_id) WHERE state = 'in_progress';

		-- NOTE: We could have a unique index on signed_raw_tx here but it would be large and expensive, enforcing the index on hash is cheaper
		CREATE UNIQUE INDEX idx_eth_tx_attempts_hash ON eth_tx_attempts (hash);
		CREATE UNIQUE INDEX idx_eth_tx_attempts_unique_gas_prices ON eth_tx_attempts (eth_tx_id, gas_price);
		CREATE INDEX idx_eth_tx_attempts_broadcast_before_block_num ON eth_tx_attempts (broadcast_before_block_num);
		CREATE INDEX idx_eth_tx_attempts_in_progress ON eth_tx_attempts (state) WHERE state = 'in_progress'::eth_tx_attempts_state;
		CREATE INDEX idx_eth_tx_attempts_created_at ON eth_tx_attempts USING BRIN (created_at);

		CREATE TABLE eth_task_run_txes (
			task_run_id uuid NOT NULL REFERENCES task_runs (id) ON DELETE CASCADE,
			eth_tx_id bigint NOT NULL REFERENCES eth_txes (id) ON DELETE CASCADE
		);

		CREATE UNIQUE INDEX idx_eth_task_run_txes_task_run_id ON eth_task_run_txes (task_run_id);
		CREATE UNIQUE INDEX idx_eth_task_run_txes_eth_tx_id ON eth_task_run_txes (eth_tx_id);

		CREATE TABLE eth_receipts (
			id BIGSERIAL PRIMARY KEY,
			tx_hash bytea NOT NULL REFERENCES eth_tx_attempts (hash) ON DELETE CASCADE,
			block_hash bytea NOT NULL,
			block_number bigint NOT NULL,
			transaction_index bigint NOT NULL,
			receipt jsonb NOT NULL,
			created_at timestamptz NOT NULL
		);

		ALTER TABLE eth_receipts ADD CONSTRAINT chk_hash_length CHECK (
			octet_length(tx_hash) = 32 AND octet_length(block_hash) = 32
		);

		CREATE INDEX idx_eth_receipts_block_number ON eth_receipts (block_number);
		CREATE UNIQUE INDEX idx_eth_receipts_unique ON eth_receipts (tx_hash, block_hash);
		CREATE INDEX idx_eth_receipts_created_at ON eth_receipts USING BRIN (created_at);
	`).Error
}
