package types

import "context"

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

// PriorAttempt provides a generic interface for reading tx data to be used in the fee esimators
type PriorAttempt[FEE any, HASH any] interface {
	Fee() FEE
	GetChainSpecificGasLimit() uint32
	GetBroadcastBeforeBlockNum() *int64
	GetHash() HASH
	GetTxType() int
}

// FeeEstimator provides a generic interface for fee estimation
//
//go:generate mockery --quiet --name FeeEstimator --output ./mocks/ --case=underscore
type FeeEstimator[HEAD any, FEE any, MAXPRICE any, HASH any] interface {
	HeadTrackable[HEAD]
	Start(context.Context) error
	Close() error

	OnNewLongestChain(ctx context.Context, head HeadView[HEAD])
	GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice MAXPRICE, opts ...Opt) (fee FEE, chainSpecificFeeLimit uint32, err error)
	BumpFee(ctx context.Context, originalFee FEE, feeLimit uint32, maxFeePrice MAXPRICE, attempts []PriorAttempt[FEE, HASH]) (bumpedFee FEE, chainSpecificFeeLimit uint32, err error)
}
