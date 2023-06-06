package cosmostxm

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"golang.org/x/exp/slices"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	cosmosclient "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	_ services.ServiceCtx = (*Txm)(nil)
	_ adapters.TxManager  = (*Txm)(nil)
)

// Txm manages transactions for the cosmos blockchain.
type Txm struct {
	starter    utils.StartStopOnce
	eb         pg.EventBroadcaster
	sub        pg.Subscription
	orm        *ORM
	lggr       logger.Logger
	tc         func() (cosmosclient.ReaderWriter, error)
	ks         keystore.Cosmos
	stop, done chan struct{}
	cfg        coscfg.Config
	gpe        cosmosclient.ComposedGasPriceEstimator
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(db *sqlx.DB, tc func() (cosmosclient.ReaderWriter, error), gpe cosmosclient.ComposedGasPriceEstimator, chainID string, cfg coscfg.Config, ks keystore.Cosmos, lggr logger.Logger, logCfg pg.QConfig, eb pg.EventBroadcaster) *Txm {
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

// Start subscribes to pg notifications about cosmos msg inserts and processes them.
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("cosmostxm", func() error {
		sub, err := txm.eb.Subscribe(pg.ChannelInsertOnCosmosMsg, "")
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
		broadcasted, err := txm.orm.GetMsgsState(db.Broadcasted, txm.cfg.MaxMsgsPerBatch())
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
		msgsByTxHash := make(map[string]adapters.Msgs)
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
	ctx, cancel := utils.StopChan(txm.stop).NewCtx()
	defer cancel()
	txm.confirmAnyUnconfirmed(ctx)
	// Jitter in case we have multiple cosmos chains each with their own client.
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

var (
	typeMsgSend            = sdk.MsgTypeURL(&types.MsgSend{})
	typeMsgExecuteContract = sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{})
)

func unmarshalMsg(msgType string, raw []byte) (sdk.Msg, string, error) {
	switch msgType {
	case typeMsgSend:
		var ms types.MsgSend
		err := ms.Unmarshal(raw)
		if err != nil {
			return nil, "", err
		}
		return &ms, ms.FromAddress, nil
	case typeMsgExecuteContract:
		var ms wasmtypes.MsgExecuteContract
		err := ms.Unmarshal(raw)
		if err != nil {
			return nil, "", err
		}
		return &ms, ms.Sender, nil
	}
	return nil, "", errors.Errorf("unrecognized message type: %s", msgType)
}

type msgValidator struct {
	cutoff         time.Time
	expired, valid adapters.Msgs
}

func (e *msgValidator) add(msg adapters.Msg) {
	if msg.CreatedAt.Before(e.cutoff) {
		e.expired = append(e.expired, msg)
	} else {
		e.valid = append(e.valid, msg)
	}
}

func (e *msgValidator) sortValid() {
	slices.SortFunc(e.valid, func(a, b adapters.Msg) bool {
		ac, bc := a.CreatedAt, b.CreatedAt
		if ac.Equal(bc) {
			return a.ID < b.ID
		}
		return ac.Before(bc)
	})
}

func (txm *Txm) sendMsgBatch(ctx context.Context) {
	msgs := msgValidator{cutoff: time.Now().Add(-txm.cfg.TxMsgTimeout())}
	err := txm.orm.q.Transaction(func(tx pg.Queryer) error {
		// There may be leftover Started messages after a crash or failed send attempt.
		started, err := txm.orm.GetMsgsState(db.Started, txm.cfg.MaxMsgsPerBatch(), pg.WithQueryer(tx))
		if err != nil {
			txm.lggr.Errorw("unable to read unstarted msgs", "err", err)
			return err
		}
		if limit := txm.cfg.MaxMsgsPerBatch() - int64(len(started)); limit > 0 {
			// Use the remaining batch budget for Unstarted
			unstarted, err := txm.orm.GetMsgsState(db.Unstarted, limit, pg.WithQueryer(tx)) //nolint
			if err != nil {
				txm.lggr.Errorw("unable to read unstarted msgs", "err", err)
				return err
			}
			for _, msg := range unstarted {
				msgs.add(msg)
			}
			// Update valid, Unstarted messages to Started
			err = txm.orm.UpdateMsgs(msgs.valid.GetIDs(), db.Started, nil, pg.WithQueryer(tx))
			if err != nil {
				// Assume transient db error retry
				txm.lggr.Errorw("unable to mark unstarted txes as started", "err", err)
				return err
			}
		}
		for _, msg := range started {
			msgs.add(msg)
		}
		// Update expired messages (Unstarted or Started) to Errored
		err = txm.orm.UpdateMsgs(msgs.expired.GetIDs(), db.Errored, nil, pg.WithQueryer(tx))
		if err != nil {
			// Assume transient db error retry
			txm.lggr.Errorw("unable to mark expired txes as errored", "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	if len(msgs.valid) == 0 {
		return
	}
	msgs.sortValid()
	txm.lggr.Debugw("building a batch", "not expired", msgs.valid, "marked expired", msgs.expired)
	var msgsByFrom = make(map[string]adapters.Msgs)
	for _, m := range msgs.valid {
		msg, sender, err2 := unmarshalMsg(m.Type, m.Raw)
		if err2 != nil {
			// Should be impossible given the check in Enqueue
			txm.lggr.Criticalw("Failed to unmarshal msg, skipping", "err", err2, "msg", m)
			continue
		}
		m.DecodedMsg = msg
		_, err2 = sdk.AccAddressFromBech32(sender)
		if err2 != nil {
			// Should never happen, we parse sender on Enqueue
			txm.lggr.Criticalw("Unable to parse sender", "err", err2, "sender", sender)
			continue
		}
		msgsByFrom[sender] = append(msgsByFrom[sender], m)
	}

	txm.lggr.Debugw("msgsByFrom", "msgsByFrom", msgsByFrom)
	gasPrice, err := txm.GasPrice()
	if err != nil {
		// Should be impossible
		txm.lggr.Criticalw("Failed to get gas price", "err", err)
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

func (txm *Txm) sendMsgBatchFromAddress(ctx context.Context, gasPrice sdk.DecCoin, sender sdk.AccAddress, key cosmoskey.Key, msgs adapters.Msgs) {
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
		// Note one rare scenario in which this can happen: the cosmos node misbehaves
		// in that it confirms a txhash is present but still gives an old seq num.
		// This is benign as the next retry will succeeds.
		return
	}
	txm.lggr.Debugw("simulation results", "from", sender, "succeeded", simResults.Succeeded, "failed", simResults.Failed)
	err = txm.orm.UpdateMsgs(simResults.Failed.GetSimMsgsIDs(), db.Errored, nil)
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
	// There is still a small chance of network failure or node/db crash after broadcasting but before committing the tx,
	// in which case the msgs would be picked up again and re-broadcast, ensuring at-least once delivery.
	var resp *txtypes.BroadcastTxResponse
	err = txm.orm.q.Transaction(func(tx pg.Queryer) error {
		txHash := strings.ToUpper(hex.EncodeToString(tmhash.Sum(signedTx)))
		err = txm.orm.UpdateMsgs(simResults.Succeeded.GetSimMsgsIDs(), db.Broadcasted, &txHash, pg.WithQueryer(tx))
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

func (txm *Txm) confirmTx(ctx context.Context, tc cosmosclient.Reader, txHash string, broadcasted []int64, maxPolls int, pollPeriod time.Duration) error {
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
		err = txm.orm.UpdateMsgs(broadcasted, db.Confirmed, nil)
		if err != nil {
			return err
		}
		return nil
	}
	txm.lggr.Errorw("unable to confirm tx after timeout period, marking errored", "hash", txHash)
	// If we are unable to confirm the tx after the timeout period
	// mark these msgs as errored
	err := txm.orm.UpdateMsgs(broadcasted, db.Errored, nil)
	if err != nil {
		txm.lggr.Errorw("unable to mark timed out txes as errored", "err", err, "txes", broadcasted, "num", len(broadcasted))
		return err
	}
	return nil
}

// Enqueue enqueue a msg destined for the cosmos chain.
func (txm *Txm) Enqueue(contractID string, msg sdk.Msg) (int64, error) {
	typeURL, raw, err := txm.marshalMsg(msg)
	if err != nil {
		return 0, err
	}

	// We could consider simulating here too, but that would
	// introduce another network call and essentially double
	// the enqueue time. Enqueue is used in the context of OCRs Transmit
	// and must be fast, so we do the minimum.

	var id int64
	err = txm.orm.q.Transaction(func(tx pg.Queryer) (err error) {
		// cancel any unstarted msgs (normally just one)
		err = txm.orm.UpdateMsgsContract(contractID, db.Unstarted, db.Errored, pg.WithQueryer(tx))
		if err != nil {
			return err
		}
		id, err = txm.orm.InsertMsg(contractID, typeURL, raw, pg.WithQueryer(tx))
		return err
	})
	return id, err
}

func (txm *Txm) marshalMsg(msg sdk.Msg) (string, []byte, error) {
	switch ms := msg.(type) {
	case *wasmtypes.MsgExecuteContract:
		_, err := sdk.AccAddressFromBech32(ms.Sender)
		if err != nil {
			txm.lggr.Errorw("failed to parse sender, skipping", "err", err, "sender", ms.Sender)
			return "", nil, err
		}

	case *types.MsgSend:
		_, err := sdk.AccAddressFromBech32(ms.FromAddress)
		if err != nil {
			txm.lggr.Errorw("failed to parse sender, skipping", "err", err, "sender", ms.FromAddress)
			return "", nil, err
		}

	default:
		return "", nil, &cosmos.ErrMsgUnsupported{Msg: msg}
	}
	typeURL := sdk.MsgTypeURL(msg)
	raw, err := proto.Marshal(msg)
	if err != nil {
		txm.lggr.Errorw("failed to marshal msg, skipping", "err", err, "msg", msg)
		return "", nil, err
	}
	return typeURL, raw, nil
}

// GetMsgs returns any messages matching ids.
func (txm *Txm) GetMsgs(ids ...int64) (adapters.Msgs, error) {
	return txm.orm.GetMsgs(ids...)
}

// GasPrice returns the gas price from the estimator in uatom.
func (txm *Txm) GasPrice() (sdk.DecCoin, error) {
	prices := txm.gpe.GasPrices()
	gasPrice, ok := prices["uatom"]
	if !ok {
		return sdk.DecCoin{}, errors.New("unexpected empty uatom price")
	}
	return gasPrice, nil
}

// Close close service
func (txm *Txm) Close() error {
	txm.sub.Close()
	close(txm.stop)
	<-txm.done
	return nil
}

func (txm *Txm) Name() string { return "cosmostxm" }

// Healthy service is healthy
func (txm *Txm) Healthy() error {
	return nil
}

// Ready service is ready
func (txm *Txm) Ready() error {
	return nil
}

func (txm *Txm) HealthReport() map[string]error { return map[string]error{txm.Name(): txm.Healthy()} }
