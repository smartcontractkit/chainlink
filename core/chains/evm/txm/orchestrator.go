package txm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type TxmOrchestrator interface {
	services.Service
	Trigger(addr common.Address)
	CreateTransaction(ctx context.Context, txRequest *types.Transaction) (id int64, err error)
	GetForwarderForEOA(ctx context.Context, eoa common.Address) (forwarder common.Address, err error)
	GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID common.Address) (forwarder common.Address, err error)
	RegisterResumeCallback(fn ResumeCallback)
	SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value *big.Int, gasLimit uint64) (tx *types.Transaction, err error)
	CountTransactionsByState(ctx context.Context, state types.TxState) (count int, err error)
	GetTransactionStatus(ctx context.Context, idempotencyKey string) (state commontypes.TransactionStatus, err error)
	//Reset(addr ADDR, abandon bool) error // Potentially will be replaced by Abandon

	// Testing methods(?)
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []types.TxState, chainID *big.Int) (txs []*types.Transaction, err error)
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []types.TxState, chainID *big.Int) (txs []*types.Transaction, err error)
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txs []*types.Transaction, err error)
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []types.TxState, chainID *big.Int) (txes []*types.Transaction, err error)
	FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error)
}

type ResumeCallback func(ctx context.Context, id uuid.UUID, result interface{}, err error) error
