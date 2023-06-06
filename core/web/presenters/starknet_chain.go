package presenters

import (
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// StarkNetChainResource is an StarkNet chain JSONAPI resource.
type StarkNetChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetChainResource) GetName() string {
	return "starknet_chain"
}

// NewStarkNetChainResource returns a new StarkNetChainResource for chain.
func NewStarkNetChainResource(chain types.ChainStatus) StarkNetChainResource {
	return StarkNetChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Config,
		Enabled: chain.Enabled,
	}}
}

// StarkNetNodeResource is a StarkNet node JSONAPI resource.
type StarkNetNodeResource struct {
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetNodeResource) GetName() string {
	return "starknet_node"
}

// NewStarkNetNodeResource returns a new StarkNetNodeResource for node.
func NewStarkNetNodeResource(node types.NodeStatus) StarkNetNodeResource {
	return StarkNetNodeResource{NodeResource{
		JAID:    NewJAID(node.Name),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
