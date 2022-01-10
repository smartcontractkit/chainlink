package terratxm

import (
	"encoding/hex"
	"fmt"
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
	_                   services.Service = (*Txm)(nil)
	failedMsgIndexRe, _                  = regexp.Compile(`^.*failed to execute message; message index: (?P<Index>\d{1}):.*$`)
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
	txm.lggr.Debugw("building a batch", "batch", unstarted)
	var msgsByFrom = make(map[string][]TerraMsg)
	for _, m := range unstarted {
		var ms wasmtypes.MsgExecuteContract
		err := ms.Unmarshal(m.Msg)
		if err != nil {
			// Should be impossible given the check in in Enqueue
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
			continue
		}
		key, err := txm.ks.Get(sender.String())
		if err != nil {
			txm.lggr.Errorw("unable to find key for from address", "err", err, "from", sender.String())
			continue
		}
		privKey := NewPrivKey(key)
		txm.lggr.Infow("sending a tx", "from", sender, "msgs", msgs)
		simResults, err := txm.simulate(msgs, sn)
		if err != nil {
			txm.lggr.Errorw("unable to estimate gas", "err", err, "from", sender.String())
			continue
		}
		// Mark failed ones
		err = txm.orm.UpdateMsgsWithState(GetIDs(simResults.failed), Errored)
		if err != nil {
			txm.lggr.Errorw("unable to mark failed sim txes as errored", "err", err, "from", sender.String())
			continue
		}
		signedTx, err := txm.tc.CreateAndSign(GetMsgs(simResults.succeeded), an, sn, simResults.gasLimit, gp, privKey)
		if err != nil {
			txm.lggr.Errorw("unable to sign tx", "err", err, "from", sender.String())
			continue
		}
		resp, err := txm.tc.Broadcast(signedTx, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
		if err != nil || resp.TxResponse == nil {
			txm.lggr.Errorw("error sending tx", "err", err, "resp", resp)
			continue
		}
		// Block mode will ensure the tx gets committed, but we still need to poll a little bit
		// until the tx hash becomes queryable and we can be certain that the next sequence number
		// can be used.
		if err := txm.ConfirmTx(resp.TxResponse.TxHash); err != nil {
			txm.lggr.Errorw("error confirming tx", "err", err, "hash", resp.TxResponse.TxHash)
			continue
		}
		// If confirmed mark these as completed.
		err = txm.orm.UpdateMsgsWithState(GetIDs(simResults.succeeded), Completed)
		if err != nil {
			return
		}
		txm.lggr.Infow("successfully sent batch", "hash", resp.TxResponse.TxHash, "msgs", msgs)
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
			toSim = toSim[failureIndex+1:]
			txm.lggr.Errorw("simulation error found in a msg", "retrying", toSim, "failure", toSim[failureIndex])
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

func (txm *Txm) ConfirmTx(txHash string) error {
	pollPeriod := 1 * time.Second
	tries := 10
	for tries = 0; tries < 10; tries++ {
		time.Sleep(pollPeriod)
		// Confirm that this tx is onchain, ensuring the sequence number has incremented
		// so we can build a new batch
		txes, err := txm.tc.TxsEvents([]string{fmt.Sprintf("tx.hash='%s'", txHash)})
		if err != nil {
			txm.lggr.Errorw("error looking for hash of tx", "err", err, "resp", txes)
			continue
		}
		if txes == nil {
			return errors.New("unexpected nil txes")
		}
		if len(txes.Txs) != 1 {
			txm.lggr.Errorw("expected one tx to be found", "txes", txes, "num", len(txes.Txs))
			return errors.New("unexpected num confirmed txes != 1")
		}
		return nil
	}
	return errors.Errorf("unable to confirm tx in %d poll periods of %v", pollPeriod, tries)
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
