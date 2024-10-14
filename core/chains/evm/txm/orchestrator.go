package txm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	txmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// TODO: use this after the migration
//type TxmOrchestrator interface {
//	services.Service
//	Trigger(addr common.Address)
//	CreateTransaction(ctx context.Context, txRequest *types.Transaction) (id int64, err error)
//	GetForwarderForEOA(ctx context.Context, eoa common.Address) (forwarder common.Address, err error)
//	GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID common.Address) (forwarder common.Address, err error)
//	RegisterResumeCallback(fn ResumeCallback)
//	SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value *big.Int, gasLimit uint64) (tx *types.Transaction, err error)
//	CountTransactionsByState(ctx context.Context, state types.TxState) (count int, err error)
//	GetTransactionStatus(ctx context.Context, idempotencyKey string) (state commontypes.TransactionStatus, err error)
//	//Reset(addr ADDR, abandon bool) error // Potentially will be replaced by Abandon
//
//	// Testing methods(?)
//	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []types.TxState, chainID *big.Int) (txs []*types.Transaction, err error)
//	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []types.TxState, chainID *big.Int) (txs []*types.Transaction, err error)
//	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txs []*types.Transaction, err error)
//	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []types.TxState, chainID *big.Int) (txes []*types.Transaction, err error)
//	FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (nullv4.Time, error)
//	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (nullv4.Int, error)
//}

type OrchestratorTxStore interface {
	FetchUnconfirmedTransactionAtNonceWithCount(context.Context, uint64, common.Address) (*txmtypes.Transaction, int, error)
	FindTxWithIdempotencyKey(context.Context, *string) (*txmtypes.Transaction, error)
}

type Orchestrator struct {
	services.StateMachine
	lggr           logger.SugaredLogger
	chainID        *big.Int
	txm            *Txm
	txStore        OrchestratorTxStore
	fwdMgr         *forwarders.FwdMgr
	resumeCallback txmgr.ResumeCallback
}

func NewTxmOrchestrator(
	lggr logger.Logger,
	chainID *big.Int,
	txm *Txm,
	txStore OrchestratorTxStore,
	fwdMgr *forwarders.FwdMgr,
) txmgr.TxManager[*big.Int, types.Head[common.Hash], common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee] {
	return &Orchestrator{
		lggr:    logger.Sugared(logger.Named(lggr, "Orchestrator")),
		chainID: chainID,
		txm:     txm,
		txStore: txStore,
		fwdMgr:  fwdMgr,
	}
}

func (o *Orchestrator) Start(ctx context.Context) error {
	return o.StartOnce("Orchestrator", func() error {
		var ms services.MultiStart
		if err := ms.Start(ctx, o.txm); err != nil {
			return fmt.Errorf("Orchestrator: Txm failed to start: %w", err)
		}
		if o.fwdMgr != nil {
			if err := ms.Start(ctx, o.fwdMgr); err != nil {
				return fmt.Errorf("Orchestrator: ForwarderManager failed to start: %w", err)
			}
		}
		return nil
	})
}

func (o *Orchestrator) Close() (merr error) {
	return o.StopOnce("Orchestrator", func() error {
		if o.fwdMgr != nil {
			if err := o.fwdMgr.Close(); err != nil {
				merr = errors.Join(merr, fmt.Errorf("Orchestrator failed to stop ForwarderManager: %w", err))
			}
		}
		if err := o.txm.Close(); err != nil {
			merr = errors.Join(merr, fmt.Errorf("Orchestrator failed to stop Txm: %w", err))
		}
		return merr
	})
}

func (o *Orchestrator) Trigger(addr common.Address) {
	if err := o.txm.Trigger(); err != nil {
		o.lggr.Error(err)
	}
}

func (o *Orchestrator) Name() string {
	return o.lggr.Name()
}

func (o *Orchestrator) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *Orchestrator) RegisterResumeCallback(fn txmgr.ResumeCallback) {
	o.resumeCallback = fn
}

func (o *Orchestrator) Reset(addr common.Address, abandon bool) error {
	ok := o.IfStarted(func() {
		if err := o.txm.Abandon(); err != nil {
			o.lggr.Error(err)
		}
	})
	if !ok {
		return fmt.Errorf("Orchestrator not started yet")
	}
	return nil
}

func (o *Orchestrator) OnNewLongestChain(ctx context.Context, head types.Head[common.Hash]) {
}

func (o *Orchestrator) CreateTransaction(ctx context.Context, request txmgrtypes.TxRequest[common.Address, common.Hash]) (tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	// TODO: Idempotency
	var wrappedTx *txmtypes.Transaction
	wrappedTx, err = o.txStore.FindTxWithIdempotencyKey(context.TODO(), request.IdempotencyKey)
	if err != nil {
		return
	}

	if wrappedTx != nil {
		o.lggr.Infof("Found Tx with IdempotencyKey: %v. Returning existing Tx without creating a new one.", *wrappedTx.IdempotencyKey)
	} else {
		var pipelineTaskRunID uuid.NullUUID
		if request.PipelineTaskRunID != nil {
			pipelineTaskRunID.UUID = *request.PipelineTaskRunID
			pipelineTaskRunID.Valid = true
		}

		if o.fwdMgr != nil && (!utils.IsZero(request.ForwarderAddress)) {
			fwdPayload, fwdErr := o.fwdMgr.ConvertPayload(request.ToAddress, request.EncodedPayload)
			if fwdErr == nil {
				// Handling meta not set at caller.
				if request.Meta != nil {
					request.Meta.FwdrDestAddress = &request.ToAddress
				} else {
					request.Meta = &txmgrtypes.TxMeta[common.Address, common.Hash]{
						FwdrDestAddress: &request.ToAddress,
					}
				}
				request.ToAddress = request.ForwarderAddress
				request.EncodedPayload = fwdPayload
			} else {
				o.lggr.Errorf("Failed to use forwarder set upstream: %v", fwdErr.Error())
			}
		}

		var meta *sqlutil.JSON
		if request.Meta != nil {
			raw, err := json.Marshal(request.Meta)
			if err != nil {
				return tx, err
			}
			m := sqlutil.JSON(raw)
			meta = &m
		}

		wrappedTxRequest := &txmtypes.TxRequest{
			IdempotencyKey:    request.IdempotencyKey,
			ChainID:           o.chainID,
			FromAddress:       request.FromAddress,
			ToAddress:         request.ToAddress,
			Value:             &request.Value,
			Data:              request.EncodedPayload,
			SpecifiedGasLimit: request.FeeLimit,
			Meta:              meta,
			ForwarderAddress:  request.ForwarderAddress,

			PipelineTaskRunID: pipelineTaskRunID,
			MinConfirmations:  request.MinConfirmations,
			SignalCallback:    request.SignalCallback,
		}

		wrappedTx, err = o.txm.CreateTransaction(ctx, wrappedTxRequest)
		if err != nil {
			return
		}
		o.txm.Trigger()
	}

	sequence := evmtypes.Nonce(wrappedTx.Nonce)
	tx = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]{
		ID:             int64(wrappedTx.ID),
		IdempotencyKey: wrappedTx.IdempotencyKey,
		Sequence:       &sequence,
		FromAddress:    wrappedTx.FromAddress,
		ToAddress:      wrappedTx.ToAddress,
		EncodedPayload: wrappedTx.Data,
		Value:          *wrappedTx.Value,
		FeeLimit:       wrappedTx.SpecifiedGasLimit,
		CreatedAt:      wrappedTx.CreatedAt,
		Meta:           wrappedTx.Meta,
		//Subject: wrappedTx.Subject,

		//TransmitChecker: wrappedTx.TransmitChecker,
		ChainID: wrappedTx.ChainID,

		PipelineTaskRunID: wrappedTx.PipelineTaskRunID,
		MinConfirmations:  wrappedTx.MinConfirmations,
		SignalCallback:    wrappedTx.SignalCallback,
		CallbackCompleted: wrappedTx.CallbackCompleted,
	}
	return
}

func (o *Orchestrator) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (uint32, error) {
	_, count, err := o.txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, common.Address{})
	return uint32(count), err
}

func (o *Orchestrator) FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (time nullv4.Time, err error) {
	return
}

func (o *Orchestrator) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (time nullv4.Int, err error) {
	return
}

func (o *Orchestrator) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

func (o *Orchestrator) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}
func (o *Orchestrator) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}
func (o *Orchestrator) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

func (o *Orchestrator) GetForwarderForEOA(ctx context.Context, eoa common.Address) (forwarder common.Address, err error) {
	if o.fwdMgr != nil {
		forwarder, err = o.fwdMgr.ForwarderFor(ctx, eoa)
	}
	return
}

func (o *Orchestrator) GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID common.Address) (forwarder common.Address, err error) {
	if o.fwdMgr != nil {
		forwarder, err = o.fwdMgr.ForwarderForOCR2Feeds(ctx, eoa, ocr2AggregatorID)
	}
	return
}

func (o *Orchestrator) GetTransactionStatus(ctx context.Context, transactionID string) (status commontypes.TransactionStatus, err error) {
	// Loads attempts and receipts in the transaction
	tx, err := o.txStore.FindTxWithIdempotencyKey(ctx, &transactionID)
	if err != nil || tx == nil {
		return status, fmt.Errorf("failed to find transaction with IdempotencyKey %s: %w", transactionID, err)
	}

	switch tx.State {
	case txmtypes.TxUnconfirmed:
		return commontypes.Pending, nil
	case txmtypes.TxConfirmed:
		// Return unconfirmed for confirmed transactions because they are not yet finalized
		return commontypes.Unconfirmed, nil
	case txmtypes.TxFinalized:
		return commontypes.Finalized, nil
	case txmtypes.TxFatalError:
		return commontypes.Fatal, nil
	default:
		return commontypes.Unknown, nil
	}
}

func (o *Orchestrator) SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value big.Int, gasLimit uint64) (tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	txRequest := txmgrtypes.TxRequest[common.Address, common.Hash]{
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: []byte{},
		Value:          value,
		FeeLimit:       gasLimit,
		//Strategy:       NewSendEveryStrategy(),
	}

	tx, err = o.CreateTransaction(ctx, txRequest)
	if err != nil {
		return
	}

	// Trigger the Txm to check for new transaction
	err = o.txm.Trigger()
	return tx, err
}
