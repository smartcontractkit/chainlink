package forwarders

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Forwarder is the external interface for ForwarderAddresses
type EVMForwarder struct {
	ID         int64
	Address    common.Address
	EVMChainID utils.Big
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
