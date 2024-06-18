package forwarders

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Forwarder is the struct for Forwarder Addresses
type Forwarder struct {
	ID      int64
	Address common.Address
	//	EVMChainID big.Big
	CreatedAt time.Time
	UpdatedAt time.Time
}
