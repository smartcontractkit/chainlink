package types

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// Opt is an option for a gas estimator
type Opt int

const (
	// OptForceRefetch forces the estimator to bust a cache if necessary
	OptForceRefetch Opt = iota
)

type Fee fmt.Stringer

// PriorAttempt provides a generic interface for reading tx data to be used in the fee esimators
//
//go:generate mockery --quiet --name PriorAttempt --output ./mocks/ --case=underscore
type PriorAttempt[F Fee, TX_HASH types.Hashable] interface {
	Fee() F
	GetChainSpecificGasLimit() uint32
	GetBroadcastBeforeBlockNum() *int64
	GetHash() TX_HASH
	GetTxType() int
}

// FeeEstimator provides a generic interface for fee estimation
//
//go:generate mockery --quiet --name FeeEstimator --output ./mocks/ --case=underscore
type FeeEstimator[H Head, F Fee, MAXPRICE any, TX_HASH types.Hashable] interface {
	services.ServiceCtx
	HeadTrackable[H]

	GetFee(ctx context.Context, calldata []byte, feeLimit uint32, maxFeePrice MAXPRICE, opts ...Opt) (fee F, chainSpecificFeeLimit uint32, err error)
	BumpFee(ctx context.Context, originalFee F, feeLimit uint32, maxFeePrice MAXPRICE, attempts []PriorAttempt[F, TX_HASH]) (bumpedFee F, chainSpecificFeeLimit uint32, err error)
}
