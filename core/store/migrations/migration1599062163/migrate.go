package migration1599062163

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		-- Drop constraints and indexes referencing old type
		ALTER TABLE eth_tx_attempts DROP CONSTRAINT chk_eth_tx_attempts_fsm;
		DROP INDEX idx_eth_tx_attempts_in_progress;
		DROP INDEX idx_only_one_in_progress_attempt_per_eth_tx;

		CREATE TYPE eth_tx_attempts_state_new AS ENUM ('in_progress', 'insufficient_eth', 'broadcast');

		-- Convert to new type, casting via text representation
		ALTER TABLE eth_tx_attempts 
		ALTER COLUMN state TYPE eth_tx_attempts_state_new 
			USING (state::text::eth_tx_attempts_state_new);

		-- and swap the types
		DROP TYPE eth_tx_attempts_state;

		ALTER TYPE eth_tx_attempts_state_new RENAME TO eth_tx_attempts_state;

		-- Add constraints and indexes back again
		ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_eth_tx_attempts_fsm CHECK (
			(state IN ('in_progress', 'insufficient_eth') AND broadcast_before_block_num IS NULL) OR state = 'broadcast'
		);
		CREATE INDEX idx_eth_tx_attempts_in_progress ON eth_tx_attempts(state enum_ops) WHERE state = 'in_progress'::eth_tx_attempts_state;
		CREATE UNIQUE INDEX idx_only_one_in_progress_attempt_per_eth_tx ON eth_tx_attempts(eth_tx_id int8_ops) WHERE state = 'in_progress'::eth_tx_attempts_state;
	`).Error
}
