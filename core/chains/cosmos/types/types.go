package types

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

// Configs manages cosmos chains and nodes.
type Configs interface {
	chains.ChainConfigs
	chains.NodeConfigs[string, db.Node]
}

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	CosmosChainID string `json:"cosmosChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}
