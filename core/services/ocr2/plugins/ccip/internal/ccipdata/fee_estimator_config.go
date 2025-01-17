package ccipdata

import (
	"context"
	"math/big"
)

type FeeEstimatorConfigReader interface {
	GetDataAvailabilityConfig(ctx context.Context) (destDAOverheadGas, destGasPerDAByte, destDAMultiplierBps int64, err error)
	ModifyGasPriceComponents(ctx context.Context, execGasPrice, daGasPrice *big.Int) (modExecGasPrice, modDAGasPrice *big.Int, err error)
}
