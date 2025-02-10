package txm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	nullv4 "gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-framework/chains"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	txmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type OrchestratorTxStore interface {
	Add(addresses ...common.Address) error
	FetchUnconfirmedTransactionAtNonceWithCount(context.Context, uint64, common.Address) (*txmtypes.Transaction, int, error)
	FindTxWithIdempotencyKey(context.Context, string) (*txmtypes.Transaction, error)
}

type OrchestratorKeystore interface {
	CheckEnabled(ctx context.Context, address common.Address, chainID *big.Int) error
	EnabledAddressesForChain(ctx context.Context, chainID *big.Int) (addresses []common.Address, err error)
}

type OrchestratorAttemptBuilder[
	BLOCK_HASH chains.Hashable,
	HEAD chains.Head[BLOCK_HASH],
] interface {
	services.Service
	OnNewLongestChain(ctx context.Context, head HEAD)
}

// Generics are necessary to keep TXMv2 backwards compatible
type Orchestrator[
	BLOCK_HASH chains.Hashable,
	HEAD chains.Head[BLOCK_HASH],
] struct {
	services.StateMachine
	lggr           logger.SugaredLogger
	chainID        *big.Int
	txm            *Txm
	txStore        OrchestratorTxStore
	fwdMgr         *forwarders.FwdMgr
	keystore       OrchestratorKeystore
	attemptBuilder OrchestratorAttemptBuilder[BLOCK_HASH, HEAD]
	resumeCallback txmgr.ResumeCallback
}

func NewTxmOrchestrator[BLOCK_HASH chains.Hashable, HEAD chains.Head[BLOCK_HASH]](
	lggr logger.Logger,
	chainID *big.Int,
	txm *Txm,
	txStore OrchestratorTxStore,
	fwdMgr *forwarders.FwdMgr,
	keystore OrchestratorKeystore,
	attemptBuilder OrchestratorAttemptBuilder[BLOCK_HASH, HEAD],
) *Orchestrator[BLOCK_HASH, HEAD] {
	return &Orchestrator[BLOCK_HASH, HEAD]{
		lggr:           logger.Sugared(logger.Named(lggr, "Orchestrator")),
		chainID:        chainID,
		txm:            txm,
		txStore:        txStore,
		keystore:       keystore,
		attemptBuilder: attemptBuilder,
		fwdMgr:         fwdMgr,
	}
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) Start(ctx context.Context) error {
	return o.StartOnce("Orchestrator", func() error {
		var ms services.MultiStart
		if err := ms.Start(ctx, o.attemptBuilder); err != nil {
			// TODO: hacky fix for DualBroadcast
			if !strings.Contains(err.Error(), "already been started once") {
				return fmt.Errorf("Orchestrator: AttemptBuilder failed to start: %w", err)
			}
		}
		addresses, err := o.keystore.EnabledAddressesForChain(ctx, o.chainID)
		if err != nil {
			return err
		}
		for _, address := range addresses {
			err := o.txStore.Add(address)
			if err != nil {
				return err
			}
		}
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

func (o *Orchestrator[BLOCK_HASH, HEAD]) Close() (merr error) {
	return o.StopOnce("Orchestrator", func() error {
		if o.fwdMgr != nil {
			if err := o.fwdMgr.Close(); err != nil {
				merr = errors.Join(merr, fmt.Errorf("Orchestrator failed to stop ForwarderManager: %w", err))
			}
		}
		if err := o.txm.Close(); err != nil {
			merr = errors.Join(merr, fmt.Errorf("Orchestrator failed to stop Txm: %w", err))
		}
		if err := o.attemptBuilder.Close(); err != nil {
			// TODO: hacky fix for DualBroadcast
			if !strings.Contains(err.Error(), "already been stopped") {
				merr = errors.Join(merr, fmt.Errorf("Orchestrator failed to stop AttemptBuilder: %w", err))
			}
		}
		return merr
	})
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) Trigger(addr common.Address) {
	o.txm.Trigger(addr)
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) Name() string {
	return o.lggr.Name()
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) RegisterResumeCallback(fn txmgr.ResumeCallback) {
	o.resumeCallback = fn
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) Reset(addr common.Address, abandon bool) error {
	ok := o.IfStarted(func() {
		if err := o.txm.Abandon(addr); err != nil {
			o.lggr.Error(err)
		}
	})
	if !ok {
		return errors.New("Orchestrator not started yet")
	}
	return nil
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) OnNewLongestChain(ctx context.Context, head HEAD) {
	ok := o.IfStarted(func() {
		o.attemptBuilder.OnNewLongestChain(ctx, head)
	})
	if !ok {
		o.lggr.Debugw("Not started; ignoring head", "head", head, "state", o.State())
	}
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) CreateTransaction(ctx context.Context, request txmgrtypes.TxRequest[common.Address, common.Hash]) (tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	var wrappedTx *txmtypes.Transaction
	if request.IdempotencyKey != nil {
		wrappedTx, err = o.txStore.FindTxWithIdempotencyKey(ctx, *request.IdempotencyKey)
		if err != nil {
			return
		}
	}

	if wrappedTx != nil {
		o.lggr.Infof("Found Tx with IdempotencyKey: %v. Returning existing Tx without creating a new one.", *wrappedTx.IdempotencyKey)
	} else {
		if kErr := o.keystore.CheckEnabled(ctx, request.FromAddress, o.chainID); kErr != nil {
			return tx, fmt.Errorf("cannot send transaction from %s on chain ID %s: %w", request.FromAddress, o.chainID.String(), kErr)
		}

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
			raw, mErr := json.Marshal(request.Meta)
			if mErr != nil {
				return tx, mErr
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
		o.txm.Trigger(request.FromAddress)
	}

	if wrappedTx.ID > math.MaxInt64 {
		return tx, fmt.Errorf("overflow for int64: %d", wrappedTx.ID)
	}

	tx = txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee]{
		ID:             int64(wrappedTx.ID),
		IdempotencyKey: wrappedTx.IdempotencyKey,
		FromAddress:    wrappedTx.FromAddress,
		ToAddress:      wrappedTx.ToAddress,
		EncodedPayload: wrappedTx.Data,
		Value:          *wrappedTx.Value,
		FeeLimit:       wrappedTx.SpecifiedGasLimit,
		CreatedAt:      wrappedTx.CreatedAt,
		Meta:           wrappedTx.Meta,
		Subject:        wrappedTx.Subject,
		ChainID:        wrappedTx.ChainID,

		PipelineTaskRunID: wrappedTx.PipelineTaskRunID,
		MinConfirmations:  wrappedTx.MinConfirmations,
		SignalCallback:    wrappedTx.SignalCallback,
		CallbackCompleted: wrappedTx.CallbackCompleted,
	}
	return
}

// CountTransactionsByState was required for backwards compatibility and it's used only for unconfirmed transactions.
func (o *Orchestrator[BLOCK_HASH, HEAD]) CountTransactionsByState(ctx context.Context, state txmgrtypes.TxState) (uint32, error) {
	addresses, err := o.keystore.EnabledAddressesForChain(ctx, o.chainID)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, address := range addresses {
		_, count, err := o.txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, address)
		if err != nil {
			return 0, err
		}
		total += count
	}

	//nolint:gosec // disable G115
	return uint32(total), nil
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) FindEarliestUnconfirmedBroadcastTime(ctx context.Context) (time nullv4.Time, err error) {
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context) (time nullv4.Int, err error) {
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

//nolint:revive // keep API backwards compatible
func (o *Orchestrator[BLOCK_HASH, HEAD]) FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []txmgrtypes.TxState, chainID *big.Int) (txs []*txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) GetForwarderForEOA(ctx context.Context, eoa common.Address) (forwarder common.Address, err error) {
	if o.fwdMgr != nil {
		forwarder, err = o.fwdMgr.ForwarderFor(ctx, eoa)
	}
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID common.Address) (forwarder common.Address, err error) {
	if o.fwdMgr != nil {
		forwarder, err = o.fwdMgr.ForwarderForOCR2Feeds(ctx, eoa, ocr2AggregatorID)
	}
	return
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) GetTransactionStatus(ctx context.Context, transactionID string) (status commontypes.TransactionStatus, err error) {
	// Loads attempts and receipts in the transaction
	tx, err := o.txStore.FindTxWithIdempotencyKey(ctx, transactionID)
	if err != nil || tx == nil {
		return status, fmt.Errorf("failed to find transaction with IdempotencyKey %s: %w", transactionID, err)
	}

	switch tx.State {
	case txmgr.TxUnconfirmed:
		return commontypes.Pending, nil
	case txmgr.TxConfirmed:
		// Return unconfirmed for confirmed transactions because they are not yet finalized
		return commontypes.Unconfirmed, nil
	case txmgr.TxFinalized:
		return commontypes.Finalized, nil
	case txmgr.TxFatalError:
		return commontypes.Fatal, nil
	default:
		return commontypes.Unknown, nil
	}
}

func (o *Orchestrator[BLOCK_HASH, HEAD]) SendNativeToken(ctx context.Context, chainID *big.Int, from, to common.Address, value big.Int, gasLimit uint64) (tx txmgrtypes.Tx[*big.Int, common.Address, common.Hash, common.Hash, evmtypes.Nonce, gas.EvmFee], err error) {
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
	o.txm.Trigger(from)
	return tx, err
}
