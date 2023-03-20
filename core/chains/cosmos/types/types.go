package types

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages cosmos chains and nodes.
type ORM interface {
	chains.ChainsORM[string, *db.ChainCfg, ChainConfig]
	chains.NodesORM[string, db.Node]

	EnsureChains([]string, ...pg.QOpt) error
}

type ChainConfig = chains.ChainConfig[string, *db.ChainCfg]

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	CosmosChainID string `json:"cosmosChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}
