package types

import "context"

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

type Fee interface {
	String() string
}

// PriorAttempt provides a generic interface for reading tx data to be used in the fee esimators
type PriorAttempt[F Fee, HASH any] interface {
	Fee() F
	GetChainSpecificGasLimit() uint32
	GetBroadcastBeforeBlockNum() *int64
	GetHash() HASH
	GetTxType() int
}

// FeeEstimator provides a generic interface for fee estimation
//
//go:generate mockery --quiet --name FeeEstimator --output ./mocks/ --case=underscore
type FeeEstimator[HEAD any, F Fee, MAXPRICE any, HASH any] interface {
	OnNewLongestChain(context.Context, HEAD)
	Start(context.Context) error
	Close() error

	GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice MAXPRICE, opts ...Opt) (fee F, chainSpecificFeeLimit uint32, err error)
	BumpFee(ctx context.Context, originalFee F, feeLimit uint32, maxFeePrice MAXPRICE, attempts []PriorAttempt[F, HASH]) (bumpedFee F, chainSpecificFeeLimit uint32, err error)
}
