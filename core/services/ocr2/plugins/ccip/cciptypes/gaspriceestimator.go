package cciptypes

import (
	"context"
	"math/big"
	"time"
)

type GasPriceEstimator interface {
	GasPriceEstimatorCommit
	GasPriceEstimatorExec
}

type GasPriceEstimatorCommit interface {
	CommonGasPriceEstimator

	// Deviates checks if p1 gas price diffs from p2 by deviation options. Input prices should not be nil.
	Deviates(p1 *big.Int, p2 *big.Int) (bool, error)
}

// GasPriceEstimatorExec provides gasPriceEstimatorCommon + features needed in exec plugin, e.g. message cost estimation.
type GasPriceEstimatorExec interface {
	CommonGasPriceEstimator

	// EstimateMsgCostUSD estimates the costs for msg execution, and converts to USD value scaled by 1e18 (e.g. 5$ = 5e18).
	EstimateMsgCostUSD(p *big.Int, wrappedNativePrice *big.Int, msg EVM2EVMOnRampCCIPSendRequestedWithMeta) (*big.Int, error)
}

// CommonGasPriceEstimator is abstraction over multi-component gas prices.
type CommonGasPriceEstimator interface {
	// GetGasPrice fetches the current gas price.
	GetGasPrice(ctx context.Context) (*big.Int, error)
	// DenoteInUSD converts the gas price to be in units of USD. Input prices should not be nil.
	DenoteInUSD(p *big.Int, wrappedNativePrice *big.Int) (*big.Int, error)
	// Median finds the median gas price in slice. If gas price has multiple components, median of each individual component should be taken. Input prices should not contain nil.
	Median(gasPrices []*big.Int) (*big.Int, error)
}

// EVM2EVMOnRampCCIPSendRequestedWithMeta helper struct to hold the send request and some metadata
type EVM2EVMOnRampCCIPSendRequestedWithMeta struct {
	EVM2EVMMessage
	BlockTimestamp time.Time
	Executed       bool
	Finalized      bool
	LogIndex       uint
	TxHash         string
}
