package config

import (
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type chainScopedConfigORM struct {
	id  utils.Big
	orm types.ChainConfigORM
}

func (o *chainScopedConfigORM) storeString(name, val string) error {
	return o.orm.StoreString(o.id, name, val)
}

func (o *chainScopedConfigORM) clear(name string) error {
	return o.orm.Clear(o.id, name)
}
