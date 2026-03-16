package types

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
)

type CommitStoreReader interface {
	ccip.CommitStoreReader
	SetGasEstimator(ctx context.Context, gpe gas.EvmFeeEstimator) error
	SetSourceMaxGasPrice(ctx context.Context, sourceMaxGasPrice *big.Int) error
}

type OffRampReader interface {
	ccip.OffRampReader
}

type OnRampReader interface {
	ccip.OnRampReader
}

type PriceRegistryReader interface {
	ccip.PriceRegistryReader
}
