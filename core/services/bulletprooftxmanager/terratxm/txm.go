package terratxm

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/terra.go/msg"
	wasmtypes "github.com/terra-money/core/x/wasm/types"

	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

var _ services.Service = (*Txm)(nil)

type Txm struct {
	starter    utils.StartStopOnce
	eb         pg.EventBroadcaster
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
	txm.lggr.Infow("building a batch", "batch", unstarted)
	// TODO: group by from address
	from := unstarted[0].From
	amino := codec.NewLegacyAmino()
	cdc := codec.NewAminoCodec(amino)
	var msgs []msg.Msg
	var ids []int64
	for _, m := range unstarted {
		var ms wasmtypes.MsgExecuteContract
		err := cdc.Unmarshal(m.Msg, &ms)
		if err != nil {
			// TODO
		}
		msgs = append(msgs, &ms)
		ids = append(ids, m.ID)
	}

	addr, _ := sdk.AccAddressFromBech32(from)
	a, _ := txm.tc.Account(addr)
	key, err := txm.ks.Get(from)
	if err != nil {
		txm.lggr.Errorw("unable to find key for from address", "err", err, "from", from)
		return
	}
	privKey := NewPrivKey(key)
	resp, err := txm.tc.SignAndBroadcast(msgs, a.GetAccountNumber(), a.GetSequence(), txm.tc.GasPrice(), privKey, txtypes.BroadcastMode_BROADCAST_MODE_BLOCK)
	if err != nil {
		// TODO
	}
	// Confirm that this tx is onchain, ensuring the sequence number has incremented
	// so we can build a new batch
	txes, err := txm.tc.TxSearch(fmt.Sprintf("tx.hash = %s", resp.TxHash))
	if err != nil {
		// TODO
	}
	if txes.TotalCount != 1 {
		// TODO
	}
	// Otherwise its definitely onchain, proceed to next batch
	err = txm.orm.UpdateMsgsWithState(ids, Completed)
	if err != nil {
		// TODO
	}
}

type PrivKey struct {
	key terrakey.Key
}

func NewPrivKey(key terrakey.Key) PrivKey {
	return PrivKey{key: key}
}

// protobuf methods (don't do anything)
func (k PrivKey) Reset()        {}
func (k PrivKey) ProtoMessage() {}
func (k PrivKey) String() string {
	return ""
}

func (k PrivKey) Bytes() []byte {
	return []byte{} // does not expose private key
}
func (k PrivKey) Sign(msg []byte) ([]byte, error) {
	return k.key.Sign(msg)
}
func (k PrivKey) PubKey() cryptotypes.PubKey {
	return k.key.PublicKey()
}
func (k PrivKey) Equals(a cryptotypes.LedgerPrivKey) bool {
	return k.PubKey().Address().String() == a.PubKey().Address().String()
}
func (k PrivKey) Type() string {
	return ""
}

func (txm *Txm) Enqueue(contractID string, msg []byte) error {
	_, err := txm.orm.InsertMsg(contractID, msg)
	return err
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
