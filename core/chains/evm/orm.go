package evm

import (
	"math/big"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type orm struct{ db *gorm.DB }

var _ types.ORM = &orm{}

func NewORM(db *gorm.DB) types.ORM {
	return &orm{db}
}

func (o *orm) LoadChains() ([]types.Chain, error) {
	var dbchains []types.Chain
	var nodes []types.Node
	// TODO: Can we use sqlx instead of gorm here
	// See: https://app.clubhouse.io/chainlinklabs/story/8781/remove-dependency-on-gorm
	if err := o.db.Where("enabled = true").Find(&dbchains).Error; err != nil {
		return nil, err
	}
	if err := o.db.Find(&nodes).Error; err != nil {
		return nil, err
	}
	// HACK: gorm can't handle non-comparable foreign keys (utils.Big cannot be
	// used with ==), so preloading is not possible. Just manually assign here
	// instead
	for i, c := range dbchains {
		for _, n := range nodes {
			if n.EVMChainID.ToInt().Cmp(c.ID.ToInt()) == 0 {
				// Performance note: quadratic
				dbchains[i].Nodes = append(dbchains[i].Nodes, n)
			}
		}
	}
	return dbchains, nil
}

// StoreString saves a string value into the config for the given chain and key
func (o *orm) StoreString(chainID *big.Int, name, val string) error {
	res := o.db.Exec(`UPDATE evm_chains SET cfg = cfg || jsonb_build_object(?::text, ?::text) WHERE id = ?`, name, val, utils.NewBig(chainID))
	if res.Error != nil {
		return errors.Wrapf(res.Error, "failed to store chain config for chain ID %s", chainID.String())
	}
	if res.RowsAffected == 0 {
		return errors.Errorf("no chain found with ID %s", chainID.String())
	}
	return nil
}

// Clear deletes a config value for the given chain and key
func (o *orm) Clear(chainID *big.Int, name string) error {
	res := o.db.Exec(`UPDATE evm_chains SET cfg = cfg - ? WHERE id = ?`, name, utils.NewBig(chainID))
	if res.Error != nil {
		return errors.Wrapf(res.Error, "failed to store chain config for chain ID %s", chainID.String())
	}
	if res.RowsAffected == 0 {
		return errors.Errorf("no chain found with ID %s", chainID.String())
	}
	return nil
}
