package presenters

import (
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains"
)

// CosmosChainResource is an Cosmos chain JSONAPI resource.
type CosmosChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r CosmosChainResource) GetName() string {
	return "cosmos_chain"
}

// NewCosmosChainResource returns a new CosmosChainResource for chain.
func NewCosmosChainResource(chain chains.ChainConfig) CosmosChainResource {
	return CosmosChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Cfg,
		Enabled: chain.Enabled,
	}}
}

// CosmosNodeResource is a Cosmos node JSONAPI resource.
type CosmosNodeResource struct {
	JAID
	Name          string `json:"name"`
	CosmosChainID string `json:"cosmosChainID"`
	TendermintURL string `json:"tendermintURL"`
}

// GetName implements the api2go EntityNamer interface
func (r CosmosNodeResource) GetName() string {
	return "cosmos_node"
}

// NewCosmosNodeResource returns a new CosmosNodeResource for node.
func NewCosmosNodeResource(node db.Node) CosmosNodeResource {
	return CosmosNodeResource{
		JAID:          NewJAID(node.Name),
		Name:          node.Name,
		CosmosChainID: node.CosmosChainID,
		TendermintURL: node.TendermintURL,
	}
}
