package migrations

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const up59 = `
ALTER TABLE vrf_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE direct_request_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE keeper_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE offchainreporting_oracle_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE flux_monitor_specs ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0059_specs_define_chains",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up59).Error
		},
		Rollback: func(db *gorm.DB) error {
			return errors.New("irreversible")
		},
	})
}
