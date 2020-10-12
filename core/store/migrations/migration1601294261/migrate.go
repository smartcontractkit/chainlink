package migration1601294261

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds a trigger that notifies listeners when a new eth_tx is inserted
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE OR REPLACE FUNCTION notifyEthTxInsertion() RETURNS TRIGGER AS $_$
        BEGIN
		PERFORM pg_notify('insert_on_eth_txes'::text, NOW()::text);
		RETURN NULL;
        END
        $_$ LANGUAGE 'plpgsql';
        CREATE TRIGGER notify_eth_tx_insertion
        AFTER INSERT ON eth_txes
        FOR EACH STATEMENT EXECUTE PROCEDURE notifyEthTxInsertion();
	`).Error
}
