//
package terratxm

import (
	"encoding/hex"
	"regexp"
	"strconv"
	"time"

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
	failedMsgIndexRe, _ = regexp.Compile(`^.*failed to execute message; message index: (?P<Index>\d{1}):.*$`)
)

const (
	// TODO: reason out a valid upper bound
	MaxMsgsPerBatch = 50
	// ~8s per block, so ~80s
	BlocksUntilTxTimeout = 10
)

type Txm struct {
	starter    utils.StartStopOnce
	eb         pg.EventBroadcaster
	sub        pg.Subscription
	ticker     *time.Ticker
	orm        *ORM
	lggr       logger.Logger
	tc         terraclient.ReaderWriter
	ks         keystore.Terra
	stop, done chan struct{}
}

func NewTxm(db *sqlx.DB, tc terraclient.ReaderWriter, ks keystore.Terra, lggr logger.Logger, cfg pg.LogConfig, eb pg.EventBroadcaster, pollPeriod time.Duration) *Txm {
	ticker := time.NewTicker(pollPeriod)
	return &Txm{
		starter: utils.StartStopOnce{},
		eb:      eb,
		orm:     NewORM(db, lggr, cfg),
		ks:      ks,
		ticker:  ticker,
		tc:      tc,
		lggr:    lggr,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

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

func (txm *Txm) run(sub pg.Subscription) {
	defer func() { txm.done <- struct{}{} }()
	// TODO: Confirm or error any txes that are in broadcasted state,
	// i.e. node crashed while confirming them.
	for {
		select {
		case <-sub.Events():
			txm.sendMsgBatch()
		case <-txm.ticker.C:
			txm.sendMsgBatch()
		case <-txm.stop:
			txm.sub.Close()
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
		msgsByFrom[ms.Sender] = append(msgsByFrom[ms.Sender], m)
	}

	txm.lggr.Debugw("msgsByFrom", "msgsByFrom", msgsByFrom)
	gp := txm.tc.GasPrice()
	for s, msgs := range msgsByFrom {
		sender, _ := sdk.AccAddressFromBech32(s)
		an, sn, err := txm.tc.Account(sender)
		if err != nil {
			txm.lggr.Errorw("to read account", "err", err, "from", sender.String())
			// If we can't read the account, assume transient api issues and leave msgs unstarted
			// to retry on next poll.
			continue
		}

		key, err := txm.ks.Get(sender.String())
		if err != nil {
			txm.lggr.Errorw("unable to find key for from address", "err", err, "from", sender.String())
			// We check the transmitter key exists when the job is added. So it would have to be deleted
			// after it was added for this to happen. Retry on next poll should the key be re-added.
			continue
		}

		txm.lggr.Debugw("simulating batch", "from", sender, "msgs", msgs)
		simResults, err := txm.simulate(msgs, sn)
		if err != nil {
			txm.lggr.Errorw("unable to simulate", "err", err, "from", sender.String())
			// If we can't simulate assume transient api issue and retry on next poll.
			continue
		}
		txm.lggr.Debugw("simulation results", "from", sender, "succeeded", simResults.succeeded, "failed", simResults.failed)
		err = txm.orm.UpdateMsgsWithState(GetIDs(simResults.failed), Errored)
		if err != nil {
			txm.lggr.Errorw("unable to mark failed sim txes as errored", "err", err, "from", sender.String())
			// If we can't mark them as failed retry on next poll. Presumably same ones will fail.
			continue
		}

		lb, err := txm.tc.LatestBlock()
		if err != nil {
			txm.lggr.Errorw("unable to get latest block", "err", err, "from", sender.String())
			continue
		}
		signedTx, err := txm.tc.CreateAndSign(GetMsgs(simResults.succeeded), an, sn, simResults.gasLimit, gp, NewPrivKey(key), uint64(lb.Block.Header.Height)+uint64(BlocksUntilTxTimeout))
		if err != nil {
			txm.lggr.Errorw("unable to sign tx", "err", err, "from", sender.String())
			continue
		}

		// We need to ensure that we either broadcast successfully and mark the tx as
		// broadcasted OR we do not broadcast successfully and we do not mark it as broadcasted.
		var resp *txtypes.BroadcastTxResponse
		err = txm.orm.q.Transaction(func(tx pg.Queryer) error {
			err = txm.orm.UpdateMsgsWithState(GetIDs(simResults.succeeded), Broadcasted, pg.WithQueryer(tx))
			if err != nil {
				return err
			}
			txm.lggr.Infow("broadcasting tx", "from", sender, "msgs", simResults.succeeded)
			resp, err = txm.tc.Broadcast(signedTx, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
			if err != nil {
				return err
			}
			if resp.TxResponse == nil {
				return errors.New("unexpected nil tx response")
			}
			return nil
		})
		if err != nil {
			txm.lggr.Errorw("error broadcasting tx", "err", err, "from", sender.String())
			// Was unable to broadcast, retry on next poll
			continue
		}

		if err := txm.ConfirmTx(resp.TxResponse.TxHash, simResults.succeeded); err != nil {
			txm.lggr.Errorw("error confirming tx", "err", err, "hash", resp.TxResponse.TxHash)
			continue
		}
	}
}

type simResults struct {
	failed    []TerraMsg
	succeeded []TerraMsg
	gasLimit  uint64
}

func (txm *Txm) simulate(msgs []TerraMsg, sequence uint64) (*simResults, error) {
	// Assumes at least one msg is present.
	// If we fail to simulate the batch, remove the offending tx
	// and try again. Repeat until we have a successful batch.
	// Keep track of failures so we can mark them as errored.
	var succeeded []TerraMsg
	var failed []TerraMsg
	toSim := msgs
	for {
		txm.lggr.Infow("simulating", "toSim", toSim)
		_, err := txm.tc.SimulateUnsigned(GetMsgs(toSim), sequence)
		containsFailure, failureIndex := txm.failedMsgIndex(err)
		if err != nil && !containsFailure {
			return nil, err
		}
		if containsFailure {
			failed = append(failed, toSim[failureIndex])
			succeeded = append(succeeded, toSim[:failureIndex]...)
			// remove offending msg and retry
			if failureIndex == len(toSim)-1 {
				// we're done, last one failed
				break
			}
			// otherwise there may be more to sim
			txm.lggr.Errorw("simulation error found in a msg", "retrying", toSim[failureIndex+1:], "failure", toSim[failureIndex], "failureIndex", failureIndex)
			toSim = toSim[failureIndex+1:]
		} else {
			// we're done they all succeeded
			succeeded = append(succeeded, toSim...)
			break
		}
	}
	// Last simulation with all successful txes to get final gas limit
	s, err := txm.tc.SimulateUnsigned(GetMsgs(succeeded), sequence)
	containsFailure, _ := txm.failedMsgIndex(err)
	if err != nil && !containsFailure {
		return nil, err
	}
	if containsFailure {
		// should never happen
		return nil, errors.Errorf("unexpected failure after successful simulation err %v", err)
	}
	return &simResults{
		failed:    failed,
		succeeded: succeeded,
		gasLimit:  s.GasInfo.GasUsed,
	}, nil
}

func (txm *Txm) failedMsgIndex(err error) (bool, int) {
	if err == nil {
		return false, 0
	}
	m := failedMsgIndexRe.FindStringSubmatch(err.Error())
	if len(m) != 2 {
		return false, 0
	}
	index, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return false, 0
	}
	return true, int(index)
}

func (txm *Txm) ConfirmTx(txHash string, broadcasted []TerraMsg) error {
	// We either mark these broadcasted txes as confirmed or errored.
	// Confirmed: we see the txhash onchain. There are no reorgs in cosmos chains.
	// Errored: we do not see the txhash onchain after waiting for N blocks worth
	// of time (plus a small buffer to account for block time variance) where N
	// is TimeoutHeight - HeightAtBroadcast. In other words, if we wait for that long
	// and the tx is not confirmed, we know it has timed out.
	pollPeriod := 1 * time.Second
	tries := 100
	for tries = 0; tries < 10; tries++ {
		time.Sleep(pollPeriod)
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
		err = txm.orm.UpdateMsgsWithState(GetIDs(broadcasted), Confirmed)
		if err != nil {
			return err
		}
		return nil
	}
	// If we are unable to confirm the tx after the timeout period
	// mark these msgs as errored
	err := txm.orm.UpdateMsgsWithState(GetIDs(broadcasted), Errored)
	if err != nil {
		txm.lggr.Errorw("unable to mark timed out txes as errored", "err", err, "txes", broadcasted, "num", len(broadcasted))
		return err
	}
	return nil
}

func (txm *Txm) Enqueue(contractID string, msg []byte) (int64, error) {
	// Double check this is an unmarshalable execute contract message.
	// Add more supported message types as needed.
	var ms wasmtypes.MsgExecuteContract
	err := ms.Unmarshal(msg)
	if err != nil {
		txm.lggr.Errorw("failed to unmarshal msg, skipping", "err", err, "msg", hex.EncodeToString(msg))
		return 0, err
	}
	// We could consider simulating here too, but that would
	// introduce another network call and essentially double
	// the enqueue time. Enqueue is used in the context of OCRs Transmit
	// and must be fast, so we do the minimum of a db write.
	return txm.orm.InsertMsg(contractID, msg)
}

func (txm *Txm) Close() error {
	txm.stop <- struct{}{}
	<-txm.done
	return nil
}

func (txm *Txm) Healthy() error {
	return nil
}

func (txm *Txm) Ready() error {
	return nil
}
