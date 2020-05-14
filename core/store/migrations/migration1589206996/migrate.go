package migration1589206996

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the requisite tables for the BulletproofTxManager
// I have tried to make an intelligent guess at the required indexes and
// constraints but this will need revisiting after the system has been finished
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	  	CREATE TABLE eth_transactions (
			id BIGSERIAL PRIMARY KEY,
			nonce bigint, 
			from_address bytea REFERENCES keys (address),
			to_address bytea NOT NULL,
			encoded_payload bytea NOT NULL,
			value numeric(78, 0) NOT NULL,
			gas_limit bigint NOT NULL,
			error text,
			created_at timestamptz NOT NULL
	  	);

		CREATE UNIQUE INDEX idx_eth_transactions_nonce_from_address ON eth_transactions (nonce, from_address);
		CREATE INDEX idx_eth_transactions_created_at ON eth_transactions USING BRIN (created_at);

		ALTER TABLE eth_transactions ADD CONSTRAINT chk_nonce_requires_from_address CHECK (
			nonce IS NULL OR from_address IS NOT NULL
		);

		ALTER TABLE eth_transactions ADD CONSTRAINT chk_nonce_may_not_be_present_with_error CHECK (
			nonce IS NULL OR error IS NULL
		);

		CREATE TABLE eth_transaction_attempts (
			id BIGSERIAL PRIMARY KEY,
		 	eth_transaction_id bigint REFERENCES eth_transactions (id) NOT NULL,
		 	gas_price numeric(78,0) NOT NULL,
		 	signed_raw_tx bytea NOT NULL,
		 	hash bytea,
		 	error text,
		 	confirmed_in_block_num bigint,
		 	confirmed_in_block_hash bytea,
		 	confirmed_at timestamptz,
		 	created_at timestamptz NOT NULL
		);

		ALTER TABLE eth_transaction_attempts ADD CONSTRAINT chk_hash_must_have_associated_fields CHECK (
			(hash IS NULL AND confirmed_in_block_num IS NULL AND confirmed_in_block_hash IS NULL AND confirmed_at IS NULL)
			OR
			(hash IS NOT NULL AND confirmed_in_block_num IS NOT NULL AND confirmed_in_block_hash IS NOT NULL AND confirmed_at IS NOT NULL)
		);

		CREATE UNIQUE INDEX idx_eth_transaction_attempts_signed_raw_tx ON eth_transaction_attempts (signed_raw_tx);
		CREATE UNIQUE INDEX idx_eth_transaction_attempts_hash ON eth_transaction_attempts (hash);
		CREATE INDEX idx_eth_transaction_attempts ON eth_transaction_attempts (eth_transaction_id);
		CREATE INDEX idx_eth_transactions_confirmed_in_block_num ON eth_transaction_attempts (confirmed_in_block_num) WHERE confirmed_in_block_num IS NOT NULL;
		CREATE INDEX idx_eth_transactions_confirmed_in_block_hash ON eth_transaction_attempts (confirmed_in_block_hash) WHERE confirmed_in_block_hash IS NOT NULL;
		CREATE INDEX idx_eth_transactions_attempts_created_at ON eth_transaction_attempts USING BRIN (created_at);

		CREATE TABLE eth_task_run_transactions (
			task_run_id uuid NOT NULL REFERENCES task_runs (id) ON DELETE CASCADE,
			eth_transaction_id bigint NOT NULL REFERENCES eth_transactions (id) ON DELETE CASCADE
		);

		CREATE UNIQUE INDEX idx_eth_task_run_transactions_task_run_id ON eth_task_run_transactions (task_run_id);
		CREATE UNIQUE INDEX idx_eth_task_run_transactions_eth_transaction_id ON eth_task_run_transactions (eth_transaction_id);
	`).Error
}
