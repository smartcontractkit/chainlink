package types

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/common/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// TxmClient is a superset of all the methods needed for the txm
type TxmClient[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	ChainClient[CHAIN_ID, ADDR, SEQ]
	TransactionClient[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]

	// receipt fetching used by confirmer
	BatchGetReceipts(
		ctx context.Context,
		attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
	) (txReceipt []R, txErr []error, err error)
}

// TransactionClient contains the methods for building, simulating, broadcasting transactions
type TransactionClient[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	ChainClient[CHAIN_ID, ADDR, SEQ]

	BatchSendTransactions(
		ctx context.Context,
		attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		bathSize int,
		lggr logger.SugaredLogger,
	) (
		txCodes []client.SendTxReturnCode,
		txErrs []error,
		broadcastTime time.Time,
		successfulTxIDs []int64,
		err error)
	SendTransactionReturnCode(
		ctx context.Context,
		tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		lggr logger.SugaredLogger,
	) (client.SendTxReturnCode, error)
	SendEmptyTransaction(
		ctx context.Context,
		newTxAttempt func(ctx context.Context, seq SEQ, feeLimit uint64, fee FEE, fromAddress ADDR) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error),
		seq SEQ,
		gasLimit uint64,
		fee FEE,
		fromAddress ADDR,
	) (txhash string, err error)
	CallContract(
		ctx context.Context,
		attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		blockNumber *big.Int,
	) (rpcErr fmt.Stringer, extractErr error)
}

// ChainClient contains the interfaces for reading chain parameters (chain id, sequences, etc)
type ChainClient[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	SEQ types.Sequence,
] interface {
	ConfiguredChainID() CHAIN_ID
	PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error)
	SequenceAt(ctx context.Context, addr ADDR, blockNum *big.Int) (SEQ, error)
}
