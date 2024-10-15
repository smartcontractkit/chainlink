package presenters

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// TronChainResource is an Tron chain JSONAPI resource.
type TronChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r TronChainResource) GetName() string {
	return "tron_chain"
}

// NewTronChainResource returns a new TronChainResource for chain.
func NewTronChainResource(chain types.ChainStatus) TronChainResource {
	return TronChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Config,
		Enabled: chain.Enabled,
	}}
}

// TronNodeResource is a Tron node JSONAPI resource.
type TronNodeResource struct {
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r TronNodeResource) GetName() string {
	return "tron_node"
}

// NewTronNodeResource returns a new TronNodeResource for node.
func NewTronNodeResource(node types.NodeStatus) TronNodeResource {
	return TronNodeResource{NodeResource{
		JAID:    NewPrefixedJAID(node.Name, node.ChainID),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
