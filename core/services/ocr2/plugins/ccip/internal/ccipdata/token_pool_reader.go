package ccipdata

import (
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery --quiet --name TokenPoolReader --filename token_pool_reader_mock.go --case=underscore
type TokenPoolReader interface {
	Address() common.Address
	Type() string
}
