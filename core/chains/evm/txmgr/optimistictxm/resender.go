package txm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// timeout value for batchSendTransactions
const batchSendTransactionTimeout = 30 * time.Second

type ResenderTxAttemptBuilder interface {
	NewAttempt(context.Context, txmgr.Tx, logger.Logger) (txmgr.TxAttempt, error)
	NewBumpTxAttempt(context.Context, txmgr.Tx, txmgr.TxAttempt, []txmgr.TxAttempt, logger.Logger) (txmgr.TxAttempt, gas.EvmFee, uint32, bool, error)
}

type ResenderTxStore interface {
	FindUnconfirmedTxsRequiringBumping(context.Context, time.Time, uint32, *big.Int, common.Address, evmtypes.Nonce) ([]txmgr.Tx, error)
	MarkTxsConfirmed(context.Context, *big.Int, common.Address, evmtypes.Nonce) error
	UpdateBroadcastAtsForUnconfirmed(context.Context, time.Time, []int64) error
}

type ResenderClient interface {
	BatchCallContextAll(context.Context, []rpc.BatchElem) error
	ConfiguredChainID() *big.Int
	SequenceAt(context.Context, common.Address, *big.Int) (evmtypes.Nonce, error)
}

type ResenderConfig struct {
	BumpAfterThreshold  time.Duration // Block inclusion time. i.e. 3 * block time
	MaxBumpCycles       int           // max cycles to apply bump - cycles * bump percent = max bumped market price
	MaxInFlight         uint32
	ResendInterval      time.Duration // block time or lower for fast chains
	RPCDefaultBatchSize uint32
}

type Resender struct {
	txAttemptBuilder ResenderTxAttemptBuilder
	lggr             logger.SugaredLogger
	txStore          ResenderTxStore
	client           ResenderClient
	ks               KeyStore
	chainID          *big.Int
	config           ResenderConfig

	ctx    context.Context
	cancel context.CancelFunc
	chDone chan struct{}
}

func NewResender(
	txAttemptBuilder ResenderTxAttemptBuilder,
	lggr logger.Logger,
	txStore ResenderTxStore,
	client ResenderClient,
	ks KeyStore,
	config ResenderConfig,
) *Resender {
	ctx, cancel := context.WithCancel(context.Background())
	return &Resender{
		txAttemptBuilder: txAttemptBuilder,
		lggr:             logger.Sugared(logger.Named(lggr, "Resender")),
		txStore:          txStore,
		client:           client,
		ks:               ks,
		chainID:          client.ConfiguredChainID(),
		config:           config,
		ctx:              ctx,
		cancel:           cancel,
		chDone:           make(chan struct{}),
	}
}

func (r *Resender) Start() {
	r.lggr.Debugf("Enabled with resend interval of %s and age threshold of %s", r.config.ResendInterval, r.config.BumpAfterThreshold)
	go r.runLoop()
}

func (r *Resender) Stop() {
	r.cancel()
	<-r.chDone
}

func (r *Resender) runLoop() {
	defer close(r.chDone)

	ticker := time.NewTicker(utils.WithJitter(r.config.ResendInterval))
	defer ticker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			if err := r.resendUnconfirmed(); err != nil {
				r.lggr.Warnw("Failed to resend unconfirmed transactions", "err", err)
			}
			r.lggr.Debug("resendUnconfirmed duration: ", time.Since(start))
		}
	}
}

func (r *Resender) resendUnconfirmed() error {
	resenderAddresses, err := r.ks.EnabledAddressesForChain(r.chainID)
	if err != nil {
		return fmt.Errorf("Resender failed getting enabled keys for chain %s: %w", r.chainID.String(), err)
	}

	olderThan := time.Now().Add(-r.config.BumpAfterThreshold)
	var allAttempts []txmgr.TxAttempt

	for _, address := range resenderAddresses {
		// Each tx equal or higher than the mined nonce is considered unconfirmed.
		sequenceAt, err := r.client.SequenceAt(r.ctx, address, nil)
		if err != nil {
			r.lggr.Errorw("Error occurred while fetching sequence for address. Skipping resend.", "address", address, "err", err)
			continue
		}

		err = r.txStore.MarkTxsConfirmed(r.ctx, r.chainID, address, sequenceAt)
		if err != nil {
			return fmt.Errorf("failed to MarkTxsConfirmed: %w", err)
		}

		// Finds transactions that are considered unconfirmed and marks them as that.
		txs, err := r.txStore.FindUnconfirmedTxsRequiringBumping(r.ctx, olderThan, r.config.MaxInFlight, r.chainID, address, sequenceAt)
		if err != nil {
			return fmt.Errorf("failed to FindUnconfirmedTxsRequiringBumping: %w", err)
		}

		for _, tx := range txs {
			marketAttempt, err := r.txAttemptBuilder.NewAttempt(r.ctx, tx, r.lggr)
			if err != nil {
				return fmt.Errorf("failed on NewTxAttempt: %w", err)
			}

			marketAttempt.Tx = tx
			r.lggr.Debug("Created market priced attempt for tx. ", marketAttempt.Tx.PrettyPrint(), marketAttempt.PrettyPrint())

			// If bumping fails, bumpedAttempt is the same as marketAttempt.
			bumpedAttempt := r.bumpAttempt(r.ctx, tx, marketAttempt)
			bumpedAttempt.Tx = tx
			allAttempts = append(allAttempts, bumpedAttempt)
		}
	}

	ctx, cancel := context.WithTimeout(r.ctx, batchSendTransactionTimeout)
	defer cancel()
	broadcastTime, successfulBroadcastIDs, err := batchSendTransactions(ctx, r.lggr, r.client, allAttempts, int(r.config.RPCDefaultBatchSize))

	if len(successfulBroadcastIDs) > 0 {
		if updateErr := r.txStore.UpdateBroadcastAtsForUnconfirmed(r.ctx, broadcastTime, successfulBroadcastIDs); updateErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to update broadcast time: %w", updateErr))
		}
	}
	if err != nil {
		return fmt.Errorf("failed to re-send transactions: %w", err)
	}

	return nil
}

func (r *Resender) bumpAttempt(ctx context.Context, tx txmgr.Tx, marketAttempt txmgr.TxAttempt) txmgr.TxAttempt {

	var bumpedFee gas.EvmFee
	var bumpedFeeLimit uint32
	var err error
	bumpedAttempt := marketAttempt

	bumpingCycles := int(time.Since(*tx.BroadcastAt) / r.config.BumpAfterThreshold / time.Nanosecond)
	bumpingCycles = min(bumpingCycles, r.config.MaxBumpCycles) // Don't bump more than MaxBumpCycles

	var i int
	for i = 0; i < bumpingCycles; i++ {
		preBumpedAttempt := bumpedAttempt
		bumpedAttempt, bumpedFee, bumpedFeeLimit, _, err = r.txAttemptBuilder.NewBumpTxAttempt(ctx, tx, bumpedAttempt, nil, r.lggr)
		if err != nil {
			r.lggr.Errorw("Failed to bump gas. Returning latest bumped attempt.", tx, "err", err)
			return preBumpedAttempt
		}
	}

	if err == nil {
		r.lggr.Debugw("Bumped market priced attempt.", "txID", tx.ID, "bumpedFee", bumpedFee.String(), "bumpedFeeLimit", bumpedFeeLimit, "hash", bumpedAttempt.Hash, "cycles", i)
	}

	return bumpedAttempt
}

func batchSendTransactions(
	ctx context.Context,
	lggr logger.SugaredLogger,
	client ResenderClient,
	attempts []txmgr.TxAttempt,
	batchSize int,
) (
	broadcastTime time.Time,
	successfulBroadcastIDs []int64,
	err error,
) {
	broadcastTime = time.Now()
	if len(attempts) == 0 {
		return broadcastTime, nil, nil
	}

	reqs := make([]rpc.BatchElem, len(attempts))
	ethTxIDs := make([]int64, len(attempts))
	hashes := make([]string, len(attempts))
	for i, attempt := range attempts {
		ethTxIDs[i] = attempt.TxID
		hashes[i] = attempt.Hash.String()
		// Decode the signed raw tx back into a Transaction object
		signedTx, decodeErr := txmgr.GetGethSignedTx(attempt.SignedRawTx)
		if decodeErr != nil {
			return broadcastTime, successfulBroadcastIDs, fmt.Errorf("failed to decode signed raw tx into Transaction object: %w", decodeErr)
		}
		// Get the canonical encoding of the Transaction object needed for the eth_sendRawTransaction request
		// The signed raw tx cannot be used directly because it uses a different encoding
		txBytes, marshalErr := signedTx.MarshalBinary()
		if marshalErr != nil {
			return broadcastTime, successfulBroadcastIDs, fmt.Errorf("failed to marshal tx into canonical encoding: %w", marshalErr)
		}
		req := rpc.BatchElem{
			Method: "eth_sendRawTransaction",
			Args:   []interface{}{hexutil.Encode(txBytes)},
			Result: &common.Hash{},
		}
		reqs[i] = req
	}

	lggr.Debugw(fmt.Sprintf("Batch sending %d unconfirmed transactions.", len(attempts)), "n", len(attempts), "ethTxIDs", ethTxIDs, "hashes", hashes)

	if batchSize == 0 {
		batchSize = len(reqs)
	}

	for i := 0; i < len(reqs); i += batchSize {
		j := i + batchSize
		if j > len(reqs) {
			j = len(reqs)
		}

		lggr.Debugw(fmt.Sprintf("Batch sending transactions %v through %v", i, j))

		if err := client.BatchCallContextAll(ctx, reqs[i:j]); err != nil {
			return broadcastTime, successfulBroadcastIDs, fmt.Errorf("failed to batch send transactions: %w", err)
		}
		for k, req := range reqs[i:j] {
			if req.Result.(*common.Hash).String() == attempts[k+i].Hash.String() {
				lggr.Debugw("Sent transaction attempt.", "tx", attempts[k+i].Tx.PrettyPrint(), "attempt", attempts[i].PrettyPrint(), "err", req.Error)
			} else {
				return broadcastTime, successfulBroadcastIDs,
					fmt.Errorf("request response and attempt hash were different. reqHash: %s , attemptHash: %s", req.Result.(*common.Hash).String(), attempts[i].Hash.String())
			}
		}
		successfulBroadcastIDs = append(successfulBroadcastIDs, ethTxIDs[i:j]...)
	}

	return broadcastTime, successfulBroadcastIDs, nil
}
