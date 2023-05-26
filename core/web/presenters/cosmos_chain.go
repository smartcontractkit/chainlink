package presenters

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
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
func NewCosmosChainResource(chain types.ChainStatus) CosmosChainResource {
	return CosmosChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Config,
		Enabled: chain.Enabled,
	}}
}

// CosmosNodeResource is a Cosmos node JSONAPI resource.
type CosmosNodeResource struct {
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r CosmosNodeResource) GetName() string {
	return "cosmos_node"
}

// NewCosmosNodeResource returns a new CosmosNodeResource for node.
func NewCosmosNodeResource(node types.NodeStatus) CosmosNodeResource {
	return CosmosNodeResource{NodeResource{
		JAID:    NewJAID(node.Name),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
