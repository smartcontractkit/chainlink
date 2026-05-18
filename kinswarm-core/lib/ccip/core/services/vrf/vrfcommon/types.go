package vrfcommon

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

type GethKeyStore interface {
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addresses ...common.Address) (common.Address, error)
}

type Config interface {
	FinalityDepth() uint32
	MinIncomingConfirmations() uint32
}

type FeeConfig interface {
	LimitDefault() uint64
	LimitJobType() config.LimitJobType
	PriceMaxKey(addr common.Address) *assets.Wei
}
