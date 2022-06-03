package presenters

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// EVMForwarderResource is an EVM forwarder JSONAPI resource.
type EVMForwarderResource struct {
	JAID
	Address    common.Address `json:"address"`
	EOA        common.Address `json:"eoa_address"`
	Dest       common.Address `json:"dest_address"`
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
		EOA:        fwd.EOA,
		Dest:       fwd.Dest,
		EVMChainID: fwd.EVMChainID,
		CreatedAt:  fwd.CreatedAt,
		UpdatedAt:  fwd.UpdatedAt,
	}
}
