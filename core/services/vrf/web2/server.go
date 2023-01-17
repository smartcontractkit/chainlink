package web2

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	evmlog "github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/lottery_consumer"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ evmlog.Listener = &vrfServer{}
	_ job.ServiceCtx  = &vrfServer{}
)

type vrfServer struct {
	utils.StartStopOnce

	j job.Job

	txManager      txmgr.TxManager    // to make request tx's on behalf of web2 clients
	logBroadcaster evmlog.Broadcaster // to get notified on fulfillments
	gethks         keystore.Eth

	chainID *big.Int

	lotteryConsumerABI     abi.ABI        // to create the lottery call
	lotteryConsumerAddress common.Address // lottery contract address
	lotteryConsumer        *lottery_consumer.LotteryConsumer

	orm  *orm
	lggr logger.Logger
	q    pg.Q

	unsubscribes []func()

	srv *http.Server
	wg  *sync.WaitGroup
}

type lotteryInput struct {
	ClientRequestID string `json:"clientRequestId"`
	LotteryType     uint8  `json:"lotteryType"`
}

// NewLottery is called by clients that want to request a new lottery outcome.
func (v *vrfServer) NewLottery(w http.ResponseWriter, r *http.Request) {
	fromAddress, err := v.gethks.GetRoundRobinAddress(v.chainID, v.fromAddresses()...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decode request body to form tx data
	var req lotteryInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: check if provided request id already seen in DB,
	// To disallow re-requests.
	clientReqID, err := hexutil.Decode(req.ClientRequestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var clientRequestID [32]byte
	copy(clientRequestID[:], clientReqID)

	// TODO: lottery type validation

	vrfExternalRequestID := new(big.Int).SetBytes(uuid.NewV4().Bytes())
	packed, err := v.lotteryConsumerABI.Methods["requestRandomness"].Inputs.Pack(
		clientRequestID,      // bytes32
		req.LotteryType,      // uint8
		vrfExternalRequestID, // uint128
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var etx txmgr.EthTx
	err = v.q.Transaction(func(tx pg.Queryer) error {
		etx, err = v.txManager.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      v.lotteryConsumerAddress,
			EncodedPayload: packed,
			Strategy:       txmgr.NewSendEveryStrategy(),
			GasLimit:       1_000_000,
		}, pg.WithQueryer(tx))
		return err
	})

	v.lggr.Infow("created lottery tx", "etx", etx.ID)
	w.WriteHeader(http.StatusOK)
}

// RedeemLottery is called by clients that want to get the outcome for their lottery.
func (v *vrfServer) RedeemLottery(w http.ResponseWriter, r *http.Request) {
	// Decode request body to form tx data
	var req lotteryInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	winningNumbers, err := v.orm.GetFulfillment([]byte(req.ClientRequestID), req.LotteryType)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusProcessing) // not ready yet
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]any{
		"winningNumbers": winningNumbers,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// HandleLog stores the winning numbers in the local DB for fast retrieval
// in RedeemLottery and updates the request transaction hash for the new
// lottery created in NewLottery.
func (v *vrfServer) HandleLog(b evmlog.Broadcast) {
	switch lg := b.DecodedLog().(type) {
	case lottery_consumer.LotteryConsumerLotteryStarted:
		err := v.orm.InsertRequest(
			lg.Request.ClientRequestId[:],
			lg.Request.LotteryType,
			lg.Request.VrfExternalRequestId.Bytes(),
			lg.Raw.Address,
			lg.Raw.TxHash,
		)
		if err != nil {
			v.lggr.Errorw("error inserting request into db", "request_log", lg)
		} else {
			v.lggr.Infow("successfully inserted request into db", "request_log", lg)
		}
	case lottery_consumer.LotteryConsumerLotterySettled:
		err := v.orm.InsertFulfillment(
			lg.Request.ClientRequestId[:],
			lg.Request.LotteryType,
			lg.Request.VrfExternalRequestId.Bytes(),
			lg.Raw.Address,
			lg.Outcome.WinningNumbers,
			lg.Raw.TxHash,
		)
		if err != nil {
			v.lggr.Errorw("error inserting fulfillment into db", "fulfill_log", lg)
		} else {
			v.lggr.Infow("successfully inserted fulfillment into db", "fulfill_log", lg)
		}
	default:
		v.lggr.Warn("unknown log type", "log", lg)
	}
}

// JobID complies with log.Listener
func (v *vrfServer) JobID() int32 {
	return v.j.ID
}

func (v *vrfServer) Start(ctx context.Context) error {
	return v.StartOnce("VRFWeb2Server", func() error {
		v.unsubscribes = append(v.unsubscribes, v.logBroadcaster.Register(v, evmlog.ListenerOpts{
			Contract: v.lotteryConsumerAddress,
			ParseLog: v.lotteryConsumer.ParseLog,
			LogsWithTopics: map[common.Hash][][]evmlog.Topic{
				lottery_consumer.LotteryConsumerLotterySettled{}.Topic(): {},
				lottery_consumer.LotteryConsumerLotteryStarted{}.Topic(): {},
			},
			MinIncomingConfirmations: 1,
		}))

		router := mux.NewRouter()

		router.HandleFunc("/lottery/new", v.NewLottery)
		router.HandleFunc("/lottery/redeem", v.RedeemLottery)

		v.srv = &http.Server{
			Handler: router,
			Addr:    "127.0.0.1:8548",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		v.wg.Add(1)
		go func() {
			defer v.wg.Done()
			log.Fatal(v.srv.ListenAndServe())
		}()

		return nil
	})
}

func (v *vrfServer) Close() error {
	return v.StopOnce("VRFWeb2Server", func() error {
		err := v.srv.Close()
		for _, unsub := range v.unsubscribes {
			unsub()
		}
		v.wg.Wait()
		return err
	})
}

func (v *vrfServer) fromAddresses() []common.Address {
	var addresses []common.Address
	for _, a := range v.j.VRFWeb2Spec.FromAddresses {
		addresses = append(addresses, a.Address())
	}
	return addresses
}
