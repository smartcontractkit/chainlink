package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// TxAttemptBuilder takes the base unsigned transaction + optional parameters (tx type, gas parameters)
// and returns a signed TxAttempt
// it is able to estimate fees and sign transactions
//
//go:generate mockery --quiet --name TxAttemptBuilder --output ./mocks/ --case=underscore
type TxAttemptBuilder[
	HEAD Head, // HEAD - chain head type
	FEE Fee, // FEE - chain fee type
	ADDR types.Hashable, // ADDR - chain address type
	TX_HASH types.Hashable, // TX_HASH - chain tx hash type
	TX any, // TX - tx type (will be replaced in future)
	TXATTEMPT any, // TXATTEMPT - tx attempt type (will be replaced  in future)
	SEQ Sequence,
] interface {
	// interfaces for running the underlying estimator
	services.ServiceCtx
	HeadTrackable[HEAD]

	// NewTxAttempt builds a transaction using the configured transaction type and fee estimator (new estimation)
	NewTxAttempt(ctx context.Context, tx TX, lggr logger.Logger, opts ...Opt) (attempt TXATTEMPT, fee FEE, feeLimit uint32, retryable bool, err error)

	// NewTxAttemptWithType builds a transaction using the configured fee estimator (new estimation) + passed in tx type
	NewTxAttemptWithType(ctx context.Context, tx TX, lggr logger.Logger, txType int, opts ...Opt) (attempt TXATTEMPT, fee FEE, feeLimit uint32, retryable bool, err error)

	// NewBumpTxAttempt builds a transaction using the configured fee estimator (bumping) + tx type from previous attempt
	// this should only be used after an initial attempt has been broadcast and the underlying gas estimator only needs to bump the fee
	NewBumpTxAttempt(ctx context.Context, tx TX, previousAttempt TXATTEMPT, priorAttempts []PriorAttempt[FEE, TX_HASH], lggr logger.Logger) (attempt TXATTEMPT, bumpedFee FEE, bumpedFeeLimit uint32, retryable bool, err error)

	// NewCustomTxAttempt builds a transaction using the passed in fee + tx type
	NewCustomTxAttempt(tx TX, fee FEE, gasLimit uint32, txType int, lggr logger.Logger) (attempt TXATTEMPT, retryable bool, err error)

	// NewEmptyTxAttempt is used in ForceRebroadcast to create a signed tx with zero value sent to the zero address
	NewEmptyTxAttempt(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt TXATTEMPT, err error)
}
