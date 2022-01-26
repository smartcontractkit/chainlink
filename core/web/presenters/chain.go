package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
)

type ChainResource struct {
	JAID
	Enabled   bool           `json:"enabled"`
	Config    types.ChainCfg `json:"config"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r ChainResource) GetName() string {
	return "chain"
}

func NewChainResource(chain types.Chain) ChainResource {
	return ChainResource{
		JAID:      NewJAIDInt64(chain.ID.ToInt().Int64()),
		Config:    chain.Cfg,
		Enabled:   chain.Enabled,
		CreatedAt: chain.CreatedAt,
		UpdatedAt: chain.UpdatedAt,
	}
}

type NodeResource struct {
	JAID
	Name       string      `json:"name"`
	EVMChainID utils.Big   `json:"evmChainID"`
	WSURL      null.String `json:"wsURL"`
	HTTPURL    null.String `json:"httpURL"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r NodeResource) GetName() string {
	return "node"
}

func NewNodeResource(node types.Node) NodeResource {
	return NodeResource{
		JAID:       NewJAIDInt32(node.ID),
		Name:       node.Name,
		EVMChainID: node.EVMChainID,
		WSURL:      node.WSURL,
		HTTPURL:    node.HTTPURL,
		CreatedAt:  node.CreatedAt,
		UpdatedAt:  node.UpdatedAt,
	}
}
