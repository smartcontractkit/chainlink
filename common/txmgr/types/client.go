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

// A generic client interface for communication with the RPC node
// Every native chain must implement independently
type Client[
	CHAINID ID,
	SEQ Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	TXRECEIPT any,
	EVENT any,
	EVENTOPS any, // event filter query options
] interface {
	// ChainID stored for quick access
	ConfiguredChainID() CHAINID
	// ChainID RPC call
	ChainID() (CHAINID, error)

	Accounts[ADDR, SEQ]
	Transactions[TX, TXHASH, TXRECEIPT]
	Events[EVENT, EVENTOPS]
	Blocks[BLOCK, BLOCKHASH]

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type Accounts[ADDR types.Hashable, SEQ Sequence] interface {
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (SEQ, error)
}

type Transactions[TX any, TXHASH types.Hashable, TXRECEIPT any] interface {
	SendTransaction(ctx context.Context, tx *TX) error
	SimulateTransaction(ctx context.Context, tx *TX) error
	TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error)
	TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error)
}

type Blocks[BLOCK any, BLOCKHASH types.Hashable] interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error)
	BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error)
	LatestBlockHeight(context.Context) (*big.Int, error)
}

type Events[EVENT any, EVENTOPS any] interface {
	FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error)
}

// TxmClient is a superset of all the methods needed for the txm
type TxmClient[
	CHAIN_ID ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ Sequence,
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
	CHAIN_ID ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ Sequence,
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
	CHAIN_ID ID,
	ADDR types.Hashable,
	SEQ Sequence,
] interface {
	ConfiguredChainID() CHAIN_ID
	PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error)
	SequenceAt(ctx context.Context, addr ADDR, blockNum *big.Int) (SEQ, error)
}
