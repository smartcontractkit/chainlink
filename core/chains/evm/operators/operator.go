package operators

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Operator is the struct for Operator Addresses
type Operator struct {
	ID        int64
	Address   common.Address
	ChainId   utils.Big
	CreatedAt time.Time
	UpdatedAt time.Time
}
