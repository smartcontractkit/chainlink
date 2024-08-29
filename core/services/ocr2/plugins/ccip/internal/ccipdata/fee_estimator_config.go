package ccipdata

import (
	"context"
)

type FeeEstimatorConfigReader interface {
	GetDataAvailabilityConfig(ctx context.Context) (destDAOverheadGas, destGasPerDAByte, destDAMultiplierBps int64, err error)
}
