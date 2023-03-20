package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/common/types"
	"github.com/smartcontractkit/chainlink/core/services"
)

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

// PriorAttempt provides a generic interface for reading tx data to be used in the fee esimators
type PriorAttempt[FEE any, TX_HASH types.Hashable] interface {
	Fee() FEE
	GetChainSpecificGasLimit() uint32
	GetBroadcastBeforeBlockNum() *int64
	GetHash() TX_HASH
	GetTxType() int
}

// FeeEstimator provides a generic interface for fee estimation
//
//go:generate mockery --quiet --name FeeEstimator --output ./mocks/ --case=underscore
type FeeEstimator[H Head, FEE any, MAXPRICE any, TX_HASH types.Hashable] interface {
	services.ServiceCtx
	HeadTrackable[H]

	GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice MAXPRICE, opts ...Opt) (fee FEE, chainSpecificFeeLimit uint32, err error)
	BumpFee(ctx context.Context, originalFee FEE, feeLimit uint32, maxFeePrice MAXPRICE, attempts []PriorAttempt[FEE, TX_HASH]) (bumpedFee FEE, chainSpecificFeeLimit uint32, err error)
}
