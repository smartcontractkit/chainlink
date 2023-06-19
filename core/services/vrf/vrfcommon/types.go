package vrfcommon

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type GethKeyStore interface {
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (common.Address, error)
}

//go:generate mockery --quiet --name Config --output ../mocks/ --case=underscore
type Config interface {
	EvmFinalityDepth() uint32
	EvmGasLimitDefault() uint32
	EvmGasLimitVRFJobType() *uint32
	KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei
	MinIncomingConfirmations() uint32
}
