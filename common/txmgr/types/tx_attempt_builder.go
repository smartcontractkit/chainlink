package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// TxAttemptBuilder takes the base unsigned transaction + optional parameters (tx type, gas parameters)
// and returns a signed TxAttempt
// it is able to estimate fees and sign transactions
//
//go:generate mockery --quiet --name TxAttemptBuilder --output ./mocks/ --case=underscore
type TxAttemptBuilder[
	CHAIN_ID types.ID, // CHAIN_ID - chain id type
	HEAD types.Head[BLOCK_HASH], // HEAD - chain head type
	ADDR types.Hashable, // ADDR - chain address type
	TX_HASH, BLOCK_HASH types.Hashable, // various chain hash types
	SEQ types.Sequence, // SEQ - chain sequence type (nonce, utxo, etc)
	FEE feetypes.Fee, // FEE - chain fee type
] interface {
	// interfaces for running the underlying estimator
	services.Service
	headtracker.HeadTrackable[HEAD, BLOCK_HASH]

	// NewTxAttempt builds a transaction using the configured transaction type and fee estimator (new estimation)
	NewTxAttempt(ctx context.Context, tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], lggr logger.Logger, opts ...feetypes.Opt) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fee FEE, feeLimit uint64, retryable bool, err error)

	// NewTxAttemptWithType builds a transaction using the configured fee estimator (new estimation) + passed in tx type
	NewTxAttemptWithType(ctx context.Context, tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], lggr logger.Logger, txType int, opts ...feetypes.Opt) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fee FEE, feeLimit uint64, retryable bool, err error)

	// NewBumpTxAttempt builds a transaction using the configured fee estimator (bumping) + tx type from previous attempt
	// this should only be used after an initial attempt has been broadcast and the underlying gas estimator only needs to bump the fee
	NewBumpTxAttempt(ctx context.Context, tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], previousAttempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], priorAttempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], lggr logger.Logger) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], bumpedFee FEE, bumpedFeeLimit uint64, retryable bool, err error)

	// NewCustomTxAttempt builds a transaction using the passed in fee + tx type
	NewCustomTxAttempt(ctx context.Context, tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], fee FEE, gasLimit uint64, txType int, lggr logger.Logger) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], retryable bool, err error)

	// NewEmptyTxAttempt is used in ForceRebroadcast to create a signed tx with zero value sent to the zero address
	NewEmptyTxAttempt(ctx context.Context, seq SEQ, feeLimit uint64, fee FEE, fromAddress ADDR) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)

	// NewPurgeTxAttempt is used to create empty transaction attempts with higher gas than the previous attempt to purge stuck transactions
	NewPurgeTxAttempt(ctx context.Context, etx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], lggr logger.Logger) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
}
