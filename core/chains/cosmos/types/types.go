package types

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

type Config interface {
	chains.ChainConfig[db.Node]
}

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	CosmosChainID string `json:"cosmosChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}
