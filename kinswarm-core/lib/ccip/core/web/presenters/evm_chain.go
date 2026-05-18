package presenters

import "github.com/smartcontractkit/chainlink-common/pkg/types"

// EVMChainResource is an EVM chain JSONAPI resource.
type EVMChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r EVMChainResource) GetName() string {
	return "evm_chain"
}

// NewEVMChainResource returns a new EVMChainResource for chain.
func NewEVMChainResource(chain types.ChainStatus) EVMChainResource {
	return EVMChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Config,
		Enabled: chain.Enabled,
	}}
}

// EVMNodeResource is an EVM node JSONAPI resource.
type EVMNodeResource struct {
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r EVMNodeResource) GetName() string {
	return "evm_node"
}

// NewEVMNodeResource returns a new EVMNodeResource for node.
func NewEVMNodeResource(node types.NodeStatus) EVMNodeResource {
	return EVMNodeResource{NodeResource{
		JAID:    NewPrefixedJAID(node.Name, node.ChainID),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
