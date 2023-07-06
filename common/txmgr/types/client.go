package types

import (
	"context"
	"fmt"
	"math/big"
	"time"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
		updateBroadcastTime func(now time.Time, txIDs []int64) error,
		attempts []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		bathSize int,
		lggr logger.Logger,
	) ([]clienttypes.SendTxReturnCode, []error, error)
	SendTransactionReturnCode(
		ctx context.Context,
		tx Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE],
		lggr logger.Logger,
	) (clienttypes.SendTxReturnCode, error)
	SendEmptyTransaction(
		ctx context.Context,
		newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error),
		seq SEQ,
		gasLimit uint32,
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
