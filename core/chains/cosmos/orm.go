package cosmos

import (
	cosmosdb "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
)

func NewORMImmut(cfgs chains.Configs[string, *cosmosdb.ChainCfg, cosmosdb.Node]) types.ORM {
	return chains.NewORMImmut(cfgs)
}
