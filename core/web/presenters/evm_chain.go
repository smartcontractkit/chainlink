package presenters

import (
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// EVMChainResource is an EVM chain JSONAPI resource.
type EVMChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r EVMChainResource) GetName() string {
	return "evm_chain"
}

// NewEVMChainResource returns a new EVMChainResource for chain.
func NewEVMChainResource(chain chains.ChainConfig) EVMChainResource {
	return EVMChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Cfg,
		Enabled: chain.Enabled,
	}}
}

// EVMNodeResource is an EVM node JSONAPI resource.
type EVMNodeResource struct {
	JAID
	Name       string      `json:"name"`
	EVMChainID utils.Big   `json:"evmChainID"`
	WSURL      null.String `json:"wsURL"`
	HTTPURL    null.String `json:"httpURL"`
	State      string      `json:"state"`
}

// GetName implements the api2go EntityNamer interface
func (r EVMNodeResource) GetName() string {
	return "evm_node"
}

// NewEVMNodeResource returns a new EVMNodeResource for node.
func NewEVMNodeResource(node evmtypes.Node) EVMNodeResource {
	return EVMNodeResource{
		JAID:       NewJAID(node.Name),
		Name:       node.Name,
		EVMChainID: node.EVMChainID,
		WSURL:      node.WSURL,
		HTTPURL:    node.HTTPURL,
		State:      node.State,
	}
}
