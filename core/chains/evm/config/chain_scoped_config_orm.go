package config

import (
	"math/big"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type chainScopedConfigORM struct {
	id *big.Int
	db *gorm.DB
}

func (o *chainScopedConfigORM) storeString(name, val string) error {
	res := o.db.Exec(`UPDATE evm_chains SET cfg = cfg || jsonb_build_object(?::text, ?::text) WHERE id = ?`, name, val, o.id.String())
	if res.Error != nil {
		return errors.Wrapf(res.Error, "failed to store chain config for chain ID %d", o.id)
	}
	if res.RowsAffected == 0 {
		return errors.Errorf("no chain found with ID %d", o.id)
	}
	return nil
}

func (o *chainScopedConfigORM) clear(name string) error {
	res := o.db.Exec(`UPDATE evm_chains SET cfg = cfg - ? WHERE id = ?`, name, o.id.String())
	if res.Error != nil {
		return errors.Wrapf(res.Error, "failed to store chain config for chain ID %d", o.id)
	}
	if res.RowsAffected == 0 {
		return errors.Errorf("no chain found with ID %d", o.id)
	}
	return nil
}
