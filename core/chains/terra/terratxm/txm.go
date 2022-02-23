package terratxm

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ services.Service = (*Txm)(nil)

// Txm manages transactions for the terra blockchain.
type Txm struct {
	starter    utils.StartStopOnce
	eb         pg.EventBroadcaster
	sub        pg.Subscription
	orm        *ORM
	lggr       logger.Logger
	tc         func() (terraclient.ReaderWriter, error)
	ks         keystore.Terra
	stop, done chan struct{}
	cfg        terra.Config
	gpe        terraclient.ComposedGasPriceEstimator
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(db *sqlx.DB, tc func() (terraclient.ReaderWriter, error), gpe terraclient.ComposedGasPriceEstimator, chainID string, cfg terra.Config, ks keystore.Terra, lggr logger.Logger, logCfg pg.LogConfig, eb pg.EventBroadcaster) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		eb:      eb,
		orm:     NewORM(chainID, db, lggr, logCfg),
		ks:      ks,
		tc:      tc,
		lggr:    lggr,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		cfg:     cfg,
		gpe:     gpe,
	}
}

// Start subscribes to pg notifications about terra msg inserts and processes them.
func (txm *Txm) Start() error {
	return txm.starter.StartOnce("terratxm", func() error {
		sub, err := txm.eb.Subscribe(pg.ChannelInsertOnTerraMsg, "")
		if err != nil {
			return err
		}
		txm.sub = sub
		go txm.run()
		return nil
	})
}

func (txm *Txm) confirmAnyUnconfirmed(ctx context.Context) {
	// Confirm any broadcasted but not confirmed txes.
	// This is an edge case if we crash after having broadcasted but before we confirm.
	for {
		broadcasted, err := txm.orm.SelectMsgsWithState(db.Broadcasted, txm.cfg.MaxMsgsPerBatch())
		if err != nil {
			// Should never happen but if so, theoretically can retry with a reboot
			txm.lggr.Criticalw("unable to look for broadcasted but unconfirmed txes", "err", err)
			return
		}
		if len(broadcasted) == 0 {
			return
		}
		tc, err := txm.tc()
		if err != nil {
			txm.lggr.Criticalw("unable to get client for handling broadcasted but unconfirmed txes", "count", len(broadcasted), "err", err)
			return
		}
		msgsByTxHash := make(map[string]terra.Msgs)
		for _, msg := range broadcasted {
			msgsByTxHash[*msg.TxHash] = append(msgsByTxHash[*msg.TxHash], msg)
		}
		for txHash, msgs := range msgsByTxHash {
			maxPolls, pollPeriod := txm.confirmPollConfig()
			err := txm.confirmTx(ctx, tc, txHash, msgs.GetIDs(), maxPolls, pollPeriod)
			if err != nil {
				txm.lggr.Errorw("unable to confirm broadcasted but unconfirmed txes", "err", err, "txhash", txHash)
				if ctx.Err() != nil {
					return
				}
			}
		}
	}
}

func (txm *Txm) run() {
	defer close(txm.done)
	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()
	txm.confirmAnyUnconfirmed(ctx)
	// Jitter in case we have multiple terra chains each with their own client.
	tick := time.After(utils.WithJitter(txm.cfg.BlockRate()))
	for {
		select {
		case <-txm.sub.Events():
			txm.sendMsgBatch(ctx)
		case <-tick:
			txm.sendMsgBatch(ctx)
			tick = time.After(utils.WithJitter(txm.cfg.BlockRate()))
		case <-txm.stop:
			return
		}
	}
}

func (txm *Txm) sendMsgBatch(ctx context.Context) {
	unstarted, err := txm.orm.SelectMsgsWithState(db.Unstarted, txm.cfg.MaxMsgsPerBatch())
	if err != nil {
		txm.lggr.Errorw("unable to read unstarted msgs", "err", err)
		return
	}
	if len(unstarted) == 0 {
		return
	}
	txm.lggr.Debugw("building a batch", "batch", unstarted)
	var msgsByFrom = make(map[string]terra.Msgs)
	for _, m := range unstarted {
		var ms wasmtypes.MsgExecuteContract
		err := ms.Unmarshal(m.Raw)
		if err != nil {
			// Should be impossible given the check in Enqueue
			txm.lggr.Criticalw("failed to unmarshal msg, skipping", "err", err, "msg", m)
			continue
		}
		m.ExecuteContract = &ms
		_, err = sdk.AccAddressFromBech32(ms.Sender)
		if err != nil {
			// Should never happen, we parse sender on Enqueue
			txm.lggr.Criticalw("unable to parse sender", "err", err, "sender", ms.Sender)
			continue
		}
		msgsByFrom[ms.Sender] = append(msgsByFrom[ms.Sender], m)
	}

	txm.lggr.Debugw("msgsByFrom", "msgsByFrom", msgsByFrom)
	prices := txm.gpe.GasPrices()
	gasPrice, ok := prices["uluna"]
	if !ok {
		// Should be impossible
		txm.lggr.Criticalw("unexpected empty uluna price")
		return
	}
	for s, msgs := range msgsByFrom {
		sender, _ := sdk.AccAddressFromBech32(s) // Already checked validity above
		key, err := txm.ks.Get(sender.String())
		if err != nil {
			txm.lggr.Errorw("unable to find key for from address", "err", err, "from", sender.String())
			// We check the transmitter key exists when the job is added. So it would have to be deleted
			// after it was added for this to happen. Retry on next poll should the key be re-added.
			continue
		}
		txm.sendMsgBatchFromAddress(ctx, gasPrice, sender, key, msgs)
		if ctx.Err() != nil {
			return
		}
	}

}

func (txm *Txm) sendMsgBatchFromAddress(ctx context.Context, gasPrice sdk.DecCoin, sender sdk.AccAddress, key terrakey.Key, msgs terra.Msgs) {
	tc, err := txm.tc()
	if err != nil {
		txm.lggr.Criticalw("unable to get client", "err", err)
		return
	}
	an, sn, err := tc.Account(sender)
	if err != nil {
		txm.lggr.Warnw("unable to read account", "err", err, "from", sender.String())
		// If we can't read the account, assume transient api issues and leave msgs unstarted
		// to retry on next poll.
		return
	}

	txm.lggr.Debugw("simulating batch", "from", sender, "msgs", msgs, "seqnum", sn)
	simResults, err := tc.BatchSimulateUnsigned(msgs.GetSimMsgs(), sn)
	if err != nil {
		txm.lggr.Warnw("unable to simulate", "err", err, "from", sender.String())
		// If we can't simulate assume transient api issue and retry on next poll.
		// Note one rare scenario in which this can happen: the terra node misbehaves
		// in that it confirms a txhash is present but still gives an old seq num.
		// This is benign as the next retry will succeeds.
		return
	}
	txm.lggr.Debugw("simulation results", "from", sender, "succeeded", simResults.Succeeded, "failed", simResults.Failed)
	err = txm.orm.UpdateMsgsWithState(simResults.Failed.GetSimMsgsIDs(), db.Errored, nil)
	if err != nil {
		txm.lggr.Errorw("unable to mark failed sim txes as errored", "err", err, "from", sender.String())
		// If we can't mark them as failed retry on next poll. Presumably same ones will fail.
		return
	}

	// Continue if there are no successful txes
	if len(simResults.Succeeded) == 0 {
		txm.lggr.Warnw("all sim msgs errored, not sending tx", "from", sender.String())
		return
	}
	// Get the gas limit for the successful batch
	s, err := tc.SimulateUnsigned(simResults.Succeeded.GetMsgs(), sn)
	if err != nil {
		// In the OCR context this should only happen upon stale report
		txm.lggr.Warnw("unexpected failure after successful simulation", "err", err)
		return
	}
	gasLimit := s.GasInfo.GasUsed

	lb, err := tc.LatestBlock()
	if err != nil {
		txm.lggr.Warnw("unable to get latest block", "err", err, "from", sender.String())
		// Assume transient api issue and retry.
		return
	}
	timeoutHeight := uint64(lb.Block.Header.Height) + uint64(txm.cfg.BlocksUntilTxTimeout())
	signedTx, err := tc.CreateAndSign(simResults.Succeeded.GetMsgs(), an, sn, gasLimit, txm.cfg.GasLimitMultiplier(),
		gasPrice, NewKeyWrapper(key), timeoutHeight)
	if err != nil {
		txm.lggr.Errorw("unable to sign tx", "err", err, "from", sender.String())
		return
	}

	// We need to ensure that we either broadcast successfully and mark the tx as
	// broadcasted OR we do not broadcast successfully and we do not mark it as broadcasted.
	// We do this by first marking it broadcasted then rolling back if the broadcast api call fails.
	var resp *txtypes.BroadcastTxResponse
	err = txm.orm.q.Transaction(func(tx pg.Queryer) error {
		txHash := strings.ToUpper(hex.EncodeToString(tmhash.Sum(signedTx)))
		err = txm.orm.UpdateMsgsWithState(simResults.Succeeded.GetSimMsgsIDs(), db.Broadcasted, &txHash, pg.WithQueryer(tx))
		if err != nil {
			return err
		}

		txm.lggr.Infow("broadcasting tx", "from", sender, "msgs", simResults.Succeeded, "gasLimit", gasLimit, "gasPrice", gasPrice.String(), "timeoutHeight", timeoutHeight, "hash", txHash)
		resp, err = tc.Broadcast(signedTx, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
		if err != nil {
			// Rollback marking as broadcasted
			// Note can happen if the node's mempool is full, where we expect errCode 20.
			return err
		}
		if resp.TxResponse == nil {
			// Rollback marking as broadcasted
			return errors.New("unexpected nil tx response")
		}
		if resp.TxResponse.TxHash != txHash {
			// Should never happen
			txm.lggr.Criticalw("txhash mismatch", "got", resp.TxResponse.TxHash, "want", txHash)
		}
		return nil
	})
	if err != nil {
		txm.lggr.Errorw("error broadcasting tx", "err", err, "from", sender.String())
		// Was unable to broadcast, retry on next poll
		return
	}

	maxPolls, pollPeriod := txm.confirmPollConfig()
	if err := txm.confirmTx(ctx, tc, resp.TxResponse.TxHash, simResults.Succeeded.GetSimMsgsIDs(), maxPolls, pollPeriod); err != nil {
		txm.lggr.Errorw("error confirming tx", "err", err, "hash", resp.TxResponse.TxHash)
		return
	}
}

func (txm *Txm) confirmPollConfig() (maxPolls int, pollPeriod time.Duration) {
	blocks := txm.cfg.BlocksUntilTxTimeout()
	blockPeriod := txm.cfg.BlockRate()
	pollPeriod = txm.cfg.ConfirmPollPeriod()
	if pollPeriod == 0 {
		// don't divide by zero
		maxPolls = 1
	} else {
		maxPolls = int((time.Duration(blocks) * blockPeriod) / pollPeriod)
	}
	return
}

func (txm *Txm) confirmTx(ctx context.Context, tc terraclient.Reader, txHash string, broadcasted []int64, maxPolls int, pollPeriod time.Duration) error {
	// We either mark these broadcasted txes as confirmed or errored.
	// Confirmed: we see the txhash onchain. There are no reorgs in cosmos chains.
	// Errored: we do not see the txhash onchain after waiting for N blocks worth
	// of time (plus a small buffer to account for block time variance) where N
	// is TimeoutHeight - HeightAtBroadcast. In other words, if we wait for that long
	// and the tx is not confirmed, we know it has timed out.
	for tries := 0; tries < maxPolls; tries++ {
		// Jitter in-case we're confirming multiple txes in parallel for different keys
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(utils.WithJitter(pollPeriod)):
		}
		// Confirm that this tx is onchain, ensuring the sequence number has incremented
		// so we can build a new batch
		tx, err := tc.Tx(txHash)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				txm.lggr.Infow("txhash not found yet, still confirming", "hash", txHash)
			} else {
				txm.lggr.Errorw("error looking for hash of tx", "err", err, "hash", txHash)
			}
			continue
		}
		// Sanity check
		if tx.TxResponse == nil || tx.TxResponse.TxHash != txHash {
			txm.lggr.Errorw("error looking for hash of tx, unexpected response", "tx", tx, "hash", txHash)
			continue
		}

		txm.lggr.Infow("successfully sent batch", "hash", txHash, "msgs", broadcasted)
		// If confirmed mark these as completed.
		err = txm.orm.UpdateMsgsWithState(broadcasted, db.Confirmed, nil)
		if err != nil {
			return err
		}
		return nil
	}
	txm.lggr.Errorw("unable to confirm tx after timeout period, marking errored", "hash", txHash)
	// If we are unable to confirm the tx after the timeout period
	// mark these msgs as errored
	err := txm.orm.UpdateMsgsWithState(broadcasted, db.Errored, nil)
	if err != nil {
		txm.lggr.Errorw("unable to mark timed out txes as errored", "err", err, "txes", broadcasted, "num", len(broadcasted))
		return err
	}
	return nil
}

// Enqueue enqueue a msg destined for the terra chain.
func (txm *Txm) Enqueue(contractID string, msg []byte) (int64, error) {
	// Double check this is an unmarshalable execute contract message.
	// Add more supported message types as needed.
	var ms wasmtypes.MsgExecuteContract
	err := ms.Unmarshal(msg)
	if err != nil {
		txm.lggr.Errorw("failed to unmarshal msg, skipping", "err", err, "msg", hex.EncodeToString(msg))
		return 0, err
	}
	_, err = sdk.AccAddressFromBech32(ms.Sender)
	if err != nil {
		txm.lggr.Errorw("failed to parse sender, skipping", "err", err, "sender", ms.Sender)
		return 0, err
	}
	// We could consider simulating here too, but that would
	// introduce another network call and essentially double
	// the enqueue time. Enqueue is used in the context of OCRs Transmit
	// and must be fast, so we do the minimum of a db write.
	return txm.orm.InsertMsg(contractID, msg)
}

// Close close service
func (txm *Txm) Close() error {
	txm.sub.Close()
	close(txm.stop)
	<-txm.done
	return nil
}

// Healthy service is healthy
func (txm *Txm) Healthy() error {
	return nil
}

// Ready service is ready
func (txm *Txm) Ready() error {
	return nil
}
