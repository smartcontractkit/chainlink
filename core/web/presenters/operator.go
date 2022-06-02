package presenters

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/chains/evm/operators"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// OperatorResource is an EVM forwarder JSONAPI resource.
type OperatorResource struct {
	JAID
	Address   common.Address `json:"address"`
	ChainID   utils.Big      `json:"chainId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r OperatorResource) GetName() string {
	return "operator"
}

// NewOperatorResourcereturns a new OperatorResource for chain.
func NewOperatorResource(operator operators.Operator) OperatorResource {
	return OperatorResource{
		JAID:      NewJAIDInt64(operator.ID),
		Address:   operator.Address,
		ChainID:   operator.ChainId,
		CreatedAt: operator.CreatedAt,
		UpdatedAt: operator.UpdatedAt,
	}
}
