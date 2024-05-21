package txmgr

import (
	"context"
	"math/big"

	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/headtracker"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// ResumeCallback is assumed to be idempotent
type ResumeCallback func(ctx context.Context, id uuid.UUID, result interface{}, err error) error

// TxManager is the main component of the transaction manager.
// It is also the interface to external callers.
//
//go:generate mockery --quiet --recursive --name TxManager --output ./mocks/ --case=underscore --structname TxManager --filename tx_manager.go
type TxManager[
	CHAIN_ID types.ID,
	HEAD types.Head[BLOCK_HASH],
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	headtracker.HeadTrackable[HEAD, BLOCK_HASH]
	services.Service
	Trigger(addr ADDR)
	CreateTransaction(ctx context.Context, txRequest txmgrtypes.TxRequest[ADDR, TX_HASH]) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	GetForwarderForEOA(eoa ADDR) (forwarder ADDR, err error)
	GetForwarderForEOAOCR2Feeds(eoa, ocr2AggregatorID ADDR) (forwarder ADDR, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendNativeToken(ctx context.Context, chainID CHAIN_ID, from, to ADDR, value big.Int, gasLimit uint64) (etx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	Reset(addr ADDR, abandon bool) error
	// Find transactions by a field in the TxMeta blob and transaction states
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided by transaction states
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txes []*txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE], err error)
	FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error)
	CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (count uint32, err error)
}
