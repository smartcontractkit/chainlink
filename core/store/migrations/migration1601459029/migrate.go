package migration1601459029

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds an extra state 'confirmed_missing_receipt' to eth_txes_state.
// Due to gorm limitations we cannot use postgres' "add to enum" functionality
// and have to do an expensive type switcheroo instead
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`
		-- Drop constraints and indexes referencing old type
		ALTER TABLE eth_txes DROP CONSTRAINT chk_eth_txes_fsm;
		DROP INDEX idx_eth_txes_state;
		DROP INDEX idx_only_one_in_progress_tx_per_account;
		DROP INDEX idx_eth_txes_min_unconfirmed_nonce_for_key;

		CREATE TYPE eth_txes_state_new AS ENUM ('unstarted', 'in_progress', 'fatal_error', 'unconfirmed', 'confirmed_missing_receipt', 'confirmed');

		-- Convert to new type, casting via text representation
		ALTER TABLE eth_txes ALTER COLUMN state SET DEFAULT NULL;
		ALTER TABLE eth_txes ALTER COLUMN state TYPE eth_txes_state_new USING (state::text::eth_txes_state_new);
	
		-- and swap the types
		DROP TYPE eth_txes_state;
		ALTER TYPE eth_txes_state_new RENAME TO eth_txes_state;
	
		-- Add constraints and indexes back again
		ALTER TABLE eth_txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
			state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL
			OR state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL
			OR state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL AND broadcast_at IS NULL
			OR state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL
			OR state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL
			OR state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL
		);
	
		ALTER TABLE eth_txes ALTER COLUMN state SET DEFAULT 'unstarted'::eth_txes_state;
	
		CREATE INDEX idx_eth_txes_state ON eth_txes(state enum_ops) WHERE state <> 'confirmed'::eth_txes_state;
		CREATE UNIQUE INDEX idx_only_one_in_progress_tx_per_account ON eth_txes(from_address bytea_ops) WHERE state = 'in_progress'::eth_txes_state;
		CREATE INDEX idx_eth_txes_min_unconfirmed_nonce_for_key ON eth_txes(nonce int8_ops,from_address bytea_ops) WHERE state = 'unconfirmed'::eth_txes_state;

	`).Error
	return err
}
