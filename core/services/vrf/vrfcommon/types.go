package vrfcommon

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

type GethKeyStore interface {
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (common.Address, error)
}

//go:generate mockery --quiet --name Config --output ../mocks/ --case=underscore
type Config interface {
	FinalityDepth() uint32
	MinIncomingConfirmations() uint32
}

//go:generate mockery --quiet --name FeeConfig --output ../mocks/ --case=underscore
type FeeConfig interface {
	LimitDefault() uint32
	LimitJobType() config.LimitJobType
	PriceMaxKey(addr common.Address) *assets.Wei
}
