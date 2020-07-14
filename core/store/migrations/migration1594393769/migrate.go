package migration1594393769

import (
	"github.com/jinzhu/gorm"
)

// Migrate ensures that heads are unique and adds parent hash for use in reorg detection
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE eth_txes DROP CONSTRAINT eth_txes_from_address_fkey;
		SELECT setval('keys_id_seq', (SELECT MAX(id)+1 FROM keys), false);
		UPDATE keys SET id = DEFAULT WHERE id=0;
		ALTER TABLE keys DROP CONSTRAINT keys_pkey, ADD PRIMARY KEY (id);
		ALTER TABLE eth_txes ADD CONSTRAINT eth_txes_from_address_fkey FOREIGN KEY (from_address) REFERENCES keys(address);
	`).Error
}
