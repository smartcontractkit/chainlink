package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
)

// AttemptBuilder takes the base unsigned transaction + optional parameters (tx type, gas parameters)
// and returns a signed TxAttempt
// it is able to estimate fees and sign transactions
// H - chain head type
// F - chain fee type
// A - chain address type
// HA - chain hash type
// GU - chain gas unit type
// T - tx type (will be replaced in future)
// TA - tx attempt type (will be replaced  in future)
//
//go:generate mockery --quiet --name AttemptBuilder --output ./mocks/ --case=underscore
type AttemptBuilder[H Head, F Fee, A any, HA any, GU any, T any, TA any] interface {
	// interfaces for running the underlying estimator
	services.ServiceCtx
	HeadTrackable[H]

	// NewAttempt builds a transaction using the configured transaction type and fee estimator (new estimation)
	NewAttempt(ctx context.Context, etx T, lggr logger.Logger, opts ...Opt) (attempt TA, fee F, feeLimit uint32, retryable bool, err error)

	// NewAttemptWithType builds a transaction using the configured fee estimator (new estimation) + passed in tx type
	NewAttemptWithType(ctx context.Context, etx T, lggr logger.Logger, txType int, opts ...Opt) (attempt TA, fee F, feeLimit uint32, retryable bool, err error)

	// NewBumpAttempt builds a transaction using the configured fee estimator (bumping) + passed in tx type
	// this should only be used after an initial attempt has been broadcast and the underlying gas estimator only needs to bump the fee
	NewBumpAttempt(ctx context.Context, etx T, previousAttempt TA, txType int, priorAttempts []PriorAttempt[F, HA], lggr logger.Logger) (attempt TA, bumpedFee F, bumpedFeeLimit uint32, retryable bool, err error)

	// NewCustomAttempt builds a transaction using the passed in fee + tx type
	NewCustomAttempt(etx T, fee F, gasLimit uint32, txType int, lggr logger.Logger) (attempt TA, retryable bool, err error)

	// FeeEstimator returns the underlying gas estimator
	FeeEstimator() FeeEstimator[H, F, GU, HA]

	// NewEmptyTransaction is used in ForceRebroadcast to create a signed tx with zero value sent to the zero address
	NewEmptyTransaction(nonce uint64, feeLimit uint32, fee F, fromAddress A) (attempt TA, err error)
}
