package ccipdata

import (
	"github.com/ethereum/go-ethereum/common"
)

type TokenPoolReader interface {
	Address() common.Address
	Type() string
}
