package forwarders

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// Forwarder is the struct for Forwarder Addresses
type Forwarder struct {
	ID         int64
	Address    common.Address
	EVMChainID big.Big
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
