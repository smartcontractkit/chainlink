package terratxm

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

var (
	_ services.Service = (*Txm)(nil)
)

const (
	// MaxMsgsPerBatch The max gas limit per block is 1_000_000_000
	// https://github.com/terra-money/core/blob/d6037b9a12c8bf6b09fe861c8ad93456aac5eebb/app/legacy/migrate.go#L69.
	// The max msg size is 10KB https://github.com/terra-money/core/blob/d6037b9a12c8bf6b09fe861c8ad93456aac5eebb/x/wasm/types/params.go#L15.
	// Our msgs are only OCR reports for now, which will not exceed that size.
	// There appears to be no gas limit per tx, only per block, so theoretically
	// we could include 1000 msgs which use up to 1M gas.
	// To be conservative and since the number of messages we'd
	// have in a batch on average roughly correponds to the number of terra ocr jobs we're running (do not expect more than 100),
	// we can set a max msgs per batch of 100.
	MaxMsgsPerBatch = 100

	// BlocksUntilTxTimeout ~8s per block, so ~80s until we give up on the tx getting confirmed
	// Anecdotally it appears anything more than 4 blocks would be an extremely long wait.
	BlocksUntilTxTimeout = 10
)

// Txm manages transactions for the terra blockchain.
type Txm struct {
	starter            utils.StartStopOnce
	eb                 pg.EventBroadcaster
	sub                pg.Subscription
	ticker             *time.Ticker
	orm                *ORM
	lggr               logger.Logger
	tc                 terraclient.ReaderWriter
	ks                 keystore.Terra
	stop, done         chan struct{}
	confirmPollPeriod  time.Duration
	confirmMaxPolls    int
	fallbackGasPrice   sdk.DecCoin
	gasLimitMultiplier float64
}

// NewTxm creates a txm
func NewTxm(db *sqlx.DB, tc terraclient.ReaderWriter, fallbackGasPrice string, gasLimitMultiplier float64, ks keystore.Terra, lggr logger.Logger, cfg pg.LogConfig, eb pg.EventBroadcaster, pollPeriod time.Duration) (*Txm, error) {
	ticker := time.NewTicker(pollPeriod)
	fgp, err := sdk.NewDecFromStr(fallbackGasPrice)
	if err != nil {
		return nil, err
	}
	return &Txm{
		starter:            utils.StartStopOnce{},
		eb:                 eb,
		orm:                NewORM(db, lggr, cfg),
		ks:                 ks,
		ticker:             ticker,
		tc:                 tc,
		lggr:               lggr,
		stop:               make(chan struct{}),
		done:               make(chan struct{}),
		confirmPollPeriod:  1 * time.Second,
		confirmMaxPolls:    100,
		fallbackGasPrice:   sdk.NewDecCoinFromDec("uluna", fgp),
		gasLimitMultiplier: gasLimitMultiplier,
	}, nil
}

// Start subscribes to pg notifications about terra msg inserts and processes them.
func (txm *Txm) Start() error {
	return txm.starter.StartOnce("terratxm", func() error {
		sub, err := txm.eb.Subscribe(pg.ChannelInsertOnTerraMsg, "")
		if err != nil {
			return err
		}
		txm.sub = sub
		go txm.run(sub)
		return nil
	})
}

func (txm *Txm) confirmAnyUnconfirmed() {
	// Confirm any broadcasted but not confirmed txes.
	// This is an edge case if we crash after having broadcasted but before we confirm.
	broadcasted, err := txm.orm.SelectMsgsWithState(Broadcasted)
	if err != nil {
		// Should never happen but if so, theoretically can retry with a reboot
		txm.lggr.Errorw("unable to look for broadcasted but unconfirmed txes", "err", err)
		return
	}
	if len(broadcasted) == 0 {
		return
	}
	msgsByTxHash := make(map[string][]TerraMsg)
	for _, msg := range broadcasted {
		msgsByTxHash[*msg.TxHash] = append(msgsByTxHash[*msg.TxHash], msg)
	}
	for txHash, msgs := range msgsByTxHash {
		err := txm.confirmTx(txHash, getIDs(msgs))
		if err != nil {
			txm.lggr.Errorw("unable to confirm broadcasted but unconfirmed txes", "err", err, "txhash", txHash)
		}
	}
}

func (txm *Txm) run(sub pg.Subscription) {
	defer close(txm.done)
	txm.confirmAnyUnconfirmed()
	for {
		select {
		case <-sub.Events():
			txm.sendMsgBatch()
		case <-txm.ticker.C:
			txm.sendMsgBatch()
		case <-txm.stop:
			return
		}
	}
}

func (txm *Txm) sendMsgBatch() {
	unstarted, err := txm.orm.SelectMsgsWithState(Unstarted)
	if err != nil {
		txm.lggr.Errorw("unable to read unstarted msgs", "err", err)
		return
	}
	if len(unstarted) == 0 {
		return
	}
	if len(unstarted) > MaxMsgsPerBatch {
		unstarted = unstarted[:MaxMsgsPerBatch+1]
	}
	txm.lggr.Debugw("building a batch", "batch", unstarted)
	var msgsByFrom = make(map[string][]TerraMsg)
	for _, m := range unstarted {
		var ms wasmtypes.MsgExecuteContract
		err := ms.Unmarshal(m.Msg)
		if err != nil {
			// Should be impossible given the check in Enqueue
			txm.lggr.Errorw("failed to unmarshal msg, skipping", "err", err, "msg", m)
			continue
		}
		m.ExecuteContract = &ms
		_, err = sdk.AccAddressFromBech32(ms.Sender)
		if err != nil {
			// Should never happen, we parse sender on Enqueue
			txm.lggr.Errorw("unable to parse sender", "err", err, "sender", ms.Sender)
			continue
		}
		msgsByFrom[ms.Sender] = append(msgsByFrom[ms.Sender], m)
	}

	txm.lggr.Debugw("msgsByFrom", "msgsByFrom", msgsByFrom)
	gp := txm.tc.GasPrice(txm.fallbackGasPrice)
	for s, msgs := range msgsByFrom {
		sender, _ := sdk.AccAddressFromBech32(s) // Already checked validity above
		key, err := txm.ks.Get(sender.String())
		if err != nil {
			txm.lggr.Errorw("unable to find key for from address", "err", err, "from", sender.String())
			// We check the transmitter key exists when the job is added. So it would have to be deleted
			// after it was added for this to happen. Retry on next poll should the key be re-added.
			return
		}
		txm.sendMsgBatchFromAddress(gp, sender, key, msgs)
	}
}

func (txm *Txm) sendMsgBatchFromAddress(gasPrice sdk.DecCoin, sender sdk.AccAddress, key terrakey.Key, msgs []TerraMsg) {
	an, sn, err := txm.tc.Account(sender)
	if err != nil {
		txm.lggr.Errorw("unable to read account", "err", err, "from", sender.String())
		// If we can't read the account, assume transient api issues and leave msgs unstarted
		// to retry on next poll.
		return
	}

	txm.lggr.Debugw("simulating batch", "from", sender, "msgs", msgs)
	simResults, err := txm.tc.BatchSimulateUnsigned(getSimMsgs(msgs), sn)
	if err != nil {
		txm.lggr.Errorw("unable to simulate", "err", err, "from", sender.String())
		// If we can't simulate assume transient api issue and retry on next poll.
		return
	}
	txm.lggr.Debugw("simulation results", "from", sender, "succeeded", simResults.Succeeded, "failed", simResults.Failed)
	err = txm.orm.UpdateMsgsWithState(getSimMsgsIDs(simResults.Failed), Errored, nil)
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
	s, err := txm.tc.SimulateUnsigned(getMsgs(simResults.Succeeded), sn)
	if err != nil {
		// Should never happen
		txm.lggr.Errorw("unexpected failure after successful simulation", "err", err)
		return
	}
	gasLimit := s.GasInfo.GasUsed

	lb, err := txm.tc.LatestBlock()
	if err != nil {
		txm.lggr.Errorw("unable to get latest block", "err", err, "from", sender.String())
		return
	}
	signedTx, err := txm.tc.CreateAndSign(getMsgs(simResults.Succeeded), an, sn, gasLimit, txm.gasLimitMultiplier,
		gasPrice, NewKeyWrapper(key), uint64(lb.Block.Header.Height)+uint64(BlocksUntilTxTimeout))
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
		err = txm.orm.UpdateMsgsWithState(getSimMsgsIDs(simResults.Succeeded), Broadcasted, &txHash, pg.WithQueryer(tx))
		if err != nil {
			return err
		}

		txm.lggr.Infow("broadcasting tx", "from", sender, "msgs", simResults.Succeeded)
		resp, err = txm.tc.Broadcast(signedTx, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
		if err != nil {
			// Rollback marking as broadcasted
			return err
		}
		if resp.TxResponse == nil {
			// Rollback marking as broadcasted
			return errors.New("unexpected nil tx response")
		}
		if resp.TxResponse.TxHash != txHash {
			// Should never happen
			txm.lggr.Errorw("txhash mismatch", "got", resp.TxResponse.TxHash, "want", txHash)
		}
		return nil
	})
	if err != nil {
		txm.lggr.Errorw("error broadcasting tx", "err", err, "from", sender.String())
		// Was unable to broadcast, retry on next poll
		return
	}

	if err := txm.confirmTx(resp.TxResponse.TxHash, getSimMsgsIDs(simResults.Succeeded)); err != nil {
		txm.lggr.Errorw("error confirming tx", "err", err, "hash", resp.TxResponse.TxHash)
		return
	}
}

func (txm *Txm) confirmTx(txHash string, broadcasted []int64) error {
	// We either mark these broadcasted txes as confirmed or errored.
	// Confirmed: we see the txhash onchain. There are no reorgs in cosmos chains.
	// Errored: we do not see the txhash onchain after waiting for N blocks worth
	// of time (plus a small buffer to account for block time variance) where N
	// is TimeoutHeight - HeightAtBroadcast. In other words, if we wait for that long
	// and the tx is not confirmed, we know it has timed out.
	for tries := 0; tries < txm.confirmMaxPolls; tries++ {
		time.Sleep(txm.confirmPollPeriod)
		// Confirm that this tx is onchain, ensuring the sequence number has incremented
		// so we can build a new batch
		tx, err := txm.tc.Tx(txHash)
		if err != nil {
			txm.lggr.Errorw("error looking for hash of tx", "err", err, "resp", txHash)
			continue
		}
		// Sanity check
		if tx.TxResponse == nil || tx.TxResponse.TxHash != txHash {
			txm.lggr.Errorw("error looking for hash of tx, unexpected response", "tx", tx, "hash", txHash)
			continue
		}
		txm.lggr.Infow("successfully sent batch", "hash", txHash, "msgs", broadcasted)
		// If confirmed mark these as completed.
		err = txm.orm.UpdateMsgsWithState(broadcasted, Confirmed, nil)
		if err != nil {
			return err
		}
		return nil
	}
	// If we are unable to confirm the tx after the timeout period
	// mark these msgs as errored
	err := txm.orm.UpdateMsgsWithState(broadcasted, Errored, nil)
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
	txm.stop <- struct{}{}
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
