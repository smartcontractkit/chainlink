package presenters

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// EVMForwarderResource is an EVM forwarder JSONAPI resource.
type EVMForwarderResource struct {
	JAID
	Address    common.Address `json:"address"`
	EVMChainID utils.Big      `json:"evmChainId"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r EVMForwarderResource) GetName() string {
	return "evm_forwarder"
}

// NewEVMForwarderResource returns a new EVMForwarderResource for chain.
func NewEVMForwarderResource(fwd forwarders.Forwarder) EVMForwarderResource {
	return EVMForwarderResource{
		JAID:       NewJAIDInt64(fwd.ID),
		Address:    fwd.Address,
		EVMChainID: fwd.EVMChainID,
		CreatedAt:  fwd.CreatedAt,
		UpdatedAt:  fwd.UpdatedAt,
	}
}
