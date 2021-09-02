package config

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type chainScopedConfigORM struct {
	id  *big.Int
	orm types.ChainConfigORM
}

func (o *chainScopedConfigORM) storeString(name, val string) error {
	return o.orm.StoreString(o.id, name, val)
}

func (o *chainScopedConfigORM) clear(name string) error {
	return o.orm.Clear(o.id, name)
}
