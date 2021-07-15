package vrf

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

const (
	// Gas costs associated with executing the fulfillment
	// aside from the seed-dependent variable costs in
	// getRandomnessFromProof (for which we use estimateGas).
	StaticFulfillExecuteGasCost = 5000 + // subID balance update
		2100 + // cold subscription balance read
		20000 + // first time oracle balance update, note first time will be 20k, but 5k subsequently
		2*2100 - // cold read oracle address and oracle balance
		4800 // request delete refund, note pre-london fork was 15k
	// Buffer to ensure that the gas received after the call to the consumer contract is
	// at least the amount they requested. Same argument as
	// https://github.com/cholladay0816/chainlink/blob/08b6fd1b910b5e9b5d20f834a09204d159a56142/contracts/src/v0.6/VRFCoordinator.sol#L201
	BufferForConsumerCallback = 6000
	// Gas costs associated with making the fulfillment call itself.
	CallFulfillGasCost = 21000 + // Base tx cost
		8005 // Static costs of argument encoding etc.
		// note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.
	// An sanity upper bound on the gas costs
	// We error the job if this is reached
	// We do not check this at request time to keep request costs low.
	TotalGasLimit = 2000000
)

type pendingRequest struct {
	confirmedAtBlock uint64
	req              *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	lb               log.Broadcast
}

type listenerV2 struct {
	utils.StartStopOnce
	cfg             Config
	l               logger.Logger
	abi             abi.ABI
	ethClient       eth.Client
	logBroadcaster  log.Broadcaster
	txm             bulletprooftxmanager.TxManager
	headBroadcaster httypes.HeadBroadcasterRegistry
	coordinator     *vrf_coordinator_v2.VRFCoordinatorV2
	pipelineRunner  pipeline.Runner
	pipelineORM     pipeline.ORM
	vorm            keystore.VRFORM
	job             job.Job
	db              *gorm.DB
	vrfks           *keystore.VRF
	gethks          *keystore.Eth
	mbLogs          *utils.Mailbox
	chStop          chan struct{}
	waitOnStop      chan struct{}
	latestHead      uint64
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	pendingLogs []pendingRequest
}

func (lsn *listenerV2) Start() error {
	return lsn.StartOnce("VRFListenerV2", func() error {
		// Take the larger of the global vs specific.
		// Note that the v2 vrf requests specify their own confirmation requirements.
		// We wait for max(minConfs, request required confs) to be safe.
		minConfs := lsn.cfg.MinIncomingConfirmations()
		if lsn.job.VRFSpec.Confirmations > lsn.cfg.MinIncomingConfirmations() {
			minConfs = lsn.job.VRFSpec.Confirmations
		}
		unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinator.Address(),
			LogsWithTopics: map[common.Hash][][]log.Topic{
				vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
					{
						log.Topic(lsn.job.VRFSpec.PublicKey.MustHash()),
					},
				},
			},
			// Do not specify min confirmations, as it varies from request to request.
		})

		// Subscribe to the head broadcaster for handling
		// per request conf requirements.
		unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)

		go gracefulpanic.WrapRecover(func() {
			lsn.run([]func(){unsubscribeLogs, unsubscribeHeadBroadcaster}, minConfs)
		})
		return nil
	})
}

func (lsn *listenerV2) Connect(head *models.Head) error {
	lsn.latestHead = uint64(head.Number)
	return nil
}

func (lsn *listenerV2) OnNewLongestChain(ctx context.Context, head models.Head) {
	// Check if any v2 logs are ready for processing.
	lsn.latestHead = uint64(head.Number)
	var remainingLogs []pendingRequest
	for _, pl := range lsn.pendingLogs {
		if pl.confirmedAtBlock <= lsn.latestHead {
			// Note below makes API calls and opens a database transaction
			// TODO: Batch these requests in a follow up.
			lsn.ProcessV2VRFRequest(pl.req, pl.lb)
		} else {
			remainingLogs = append(remainingLogs, pl)
		}
	}
	lsn.pendingLogs = remainingLogs
}

func (lsn *listenerV2) run(unsubscribeLogs []func(), minConfs uint32) {
	lsn.l.Infow("VRFListenerV2: listening for run requests",
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.chStop:
			for _, us := range unsubscribeLogs {
				us()
			}
			lsn.waitOnStop <- struct{}{}
			return
		case <-lsn.mbLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				i, exists := lsn.mbLogs.Retrieve()
				if !exists {
					break
				}
				lb, ok := i.(log.Broadcast)
				if !ok {
					panic(fmt.Sprintf("VRFListenerV2: invariant violated, expected log.Broadcast got %T", i))
				}
				alreadyConsumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lsn.db, lb)
				if err != nil {
					lsn.l.Errorw("VRFListenerV2: could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
					continue
				} else if alreadyConsumed {
					continue
				}
				req, err := lsn.coordinator.ParseRandomWordsRequested(lb.RawLog())
				if err != nil {
					lsn.l.Errorw("VRFListenerV2: failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
					lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
					return
				}
				lsn.pendingLogs = append(lsn.pendingLogs, pendingRequest{
					confirmedAtBlock: req.Raw.BlockNumber + uint64(req.MinimumRequestConfirmations),
					req:              req,
					lb:               lb,
				})
			}
		}
	}
}

func (lsn *listenerV2) ProcessV2VRFRequest(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, lb log.Broadcast) {
	// Check if the vrf req has already been fulfilled
	callback, err := lsn.coordinator.GetCallbackHash(nil, req.PreSeedAndRequestId)
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: unable to check if already fulfilled, processing anyways", "err", err, "txHash", req.Raw.TxHash)
	} else if utils.IsEmpty(callback[:]) {
		// If seedAndBlockNumber is zero then the response has been fulfilled
		// and we should skip it
		lsn.l.Infow("VRFListenerV2: request already fulfilled", "txHash", req.Raw.TxHash, "subID", req.SubId, "callback", callback)
		lsn.l.ErrorIf(lsn.logBroadcaster.MarkConsumed(lsn.db, lb), "failed to mark consumed")
		return
	}

	s := time.Now()
	proof, err1 := lsn.LogToProof(req, lb)
	gasLimit, err2 := lsn.computeTxGasLimit(uint64(req.CallbackGasLimit), proof)
	vrfCoordinatorPayload, _, err3 := lsn.ProcessLogV2(proof)
	err = multierr.Combine(err1, err2, err3)
	if err != nil {
		logger.Errorw("VRFListenerV2: error processing random request", "err", err, "txHash", req.Raw.TxHash)
	}
	f := time.Now()
	err = postgres.GormTransactionWithDefaultContext(lsn.db, func(tx *gorm.DB) error {
		if err == nil {
			// No errors processing the log, submit a transaction
			var etx models.EthTx
			var from common.Address
			from, err = lsn.gethks.GetRoundRobinAddress()
			if err != nil {
				return err
			}
			etx, err = lsn.txm.CreateEthTransaction(tx,
				from,
				lsn.coordinator.Address(),
				vrfCoordinatorPayload,
				gasLimit,
				&models.EthTxMetaV2{
					JobID:         lsn.job.ID,
					RequestTxHash: lb.RawLog().TxHash,
				},
				bulletprooftxmanager.SendEveryStrategy{},
			)
			if err != nil {
				return err
			}
			// TODO: Once we have eth tasks supported, we can use the pipeline directly
			// and be able to save errored proof generations. Until then only save
			// successful runs and log errors.
			_, err = lsn.pipelineRunner.InsertFinishedRun(tx, pipeline.Run{
				PipelineSpecID: lsn.job.PipelineSpecID,
				Errors:         []null.String{{}},
				Outputs: pipeline.JSONSerializable{
					Val: []interface{}{fmt.Sprintf("queued tx from %v to %v txdata %v",
						etx.FromAddress,
						etx.ToAddress,
						hex.EncodeToString(etx.EncodedPayload))},
				},
				Meta: pipeline.JSONSerializable{
					Val: map[string]interface{}{"eth_tx_id": etx.ID},
				},
				CreatedAt:  s,
				FinishedAt: &f,
			}, nil, false)
			if err != nil {
				return errors.Wrap(err, "VRFListenerV2: failed to insert finished run")
			}
		}
		// Always mark consumed regardless of whether the proof failed or not.
		err = lsn.logBroadcaster.MarkConsumed(tx, lb)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		lsn.l.Errorw("VRFListenerV2 failed to save run", "err", err)
	}
}

// Compute the gasLimit required for the fulfillment transaction
// such that the user gets their requested amount of gas.
// We only estimate the getRandomnessFromProof as opposed to the whole fulfillRandomWords
// to avoid including the users contract code as part of the estimate for security concerns.
func (lsn *listenerV2) computeTxGasLimit(requestedCallbackGas uint64, proof []byte) (uint64, error) {
	vrfCoordinatorArgs, err := lsn.abi.Methods["getRandomnessFromProof"].Inputs.PackValues(
		[]interface{}{
			proof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: error building fulfill args", "err", err)
		return 0, err
	}
	to := lsn.coordinator.Address()
	variableFulfillmentCost, err := lsn.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &to,
		Data: append(lsn.abi.Methods["getRandomnessFromProof"].ID, vrfCoordinatorArgs...),
	})
	if err != nil {
		return 0, err
	}
	/* The fulfillment can be summarized as follows:
		fulfill {
			   1. get and verify randomness (max cost = func(seed))
	           2. user callback (max cost = requested gas limit)
			   3. calculate payment and pay oracles (max cost = deterministic)
		}
	   For step 3 see definition of StaticFulfillExecuteGasCost.
	   Note we do not include CallFulfillGasCost, as that is included as
	   part of the EstimateGas. It also has the same argument as the outer
	   fulfillment contract method, so has the same encoding costs.
	*/
	staticVerifyGas := uint64(BufferForConsumerCallback + StaticFulfillExecuteGasCost)
	totalGas := variableFulfillmentCost + requestedCallbackGas + staticVerifyGas
	logger.Infow("VRFListenerV2: estimated gas limit for tx", "gasLimit", totalGas, "callbackLimit", requestedCallbackGas)

	// Sanity check total gas
	if totalGas > TotalGasLimit {
		return totalGas, errors.Errorf("total gas limit exceeded %v > %v", totalGas, TotalGasLimit)
	}
	return totalGas, nil
}

func (lsn *listenerV2) LogToProof(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, lb log.Broadcast) ([]byte, error) {
	lsn.l.Infow("VRFListenerV2: received log request",
		"log", lb.String(),
		"reqID", req.PreSeedAndRequestId.String(),
		"keyHash", hex.EncodeToString(req.KeyHash[:]),
		"txHash", req.Raw.TxHash,
		"blockNumber", req.Raw.BlockNumber,
		"seed", req.PreSeedAndRequestId.String())
	// Validate the key against the spec
	kh, err := lsn.job.VRFSpec.PublicKey.Hash()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(req.KeyHash[:], kh[:]) {
		return nil, fmt.Errorf("invalid key hash %v expected %v", hex.EncodeToString(req.KeyHash[:]), hex.EncodeToString(kh[:]))
	}

	// req.PreSeed is uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)))
	preSeed, err := BigToSeed(req.PreSeedAndRequestId)
	if err != nil {
		return nil, errors.New("unable to parse preseed")
	}
	seed := PreSeedDataV2{
		PreSeed:          preSeed,
		BlockHash:        req.Raw.BlockHash,
		BlockNum:         req.Raw.BlockNumber,
		SubId:            req.SubId,
		CallbackGasLimit: req.CallbackGasLimit,
		NumWords:         req.NumWords,
		Sender:           req.Sender,
	}
	solidityProof, err := GenerateProofResponseV2(lsn.vrfks, lsn.job.VRFSpec.PublicKey, seed)
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: error generating proof", "err", err)
		return nil, err
	}
	return solidityProof[:], nil
}

func (lsn *listenerV2) ProcessLogV2(solidityProof []byte) ([]byte, *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, error) {
	vrfCoordinatorArgs, err := lsn.abi.Methods["fulfillRandomWords"].Inputs.PackValues(
		[]interface{}{
			solidityProof[:], // geth expects slice, even if arg is constant-length
		})
	if err != nil {
		lsn.l.Errorw("VRFListenerV2: error building fulfill args", "err", err)
		return nil, nil, err
	}

	return append(lsn.abi.Methods["fulfillRandomWords"].ID, vrfCoordinatorArgs...), nil, nil
}

// Close complies with job.Service
func (lsn *listenerV2) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		<-lsn.waitOnStop
		return nil
	})
}

func (lsn *listenerV2) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.mbLogs.Deliver(lb)
	if wasOverCapacity {
		logger.Error("VRFListenerV2: log mailbox is over capacity - dropped the oldest log")
	}
}

// JobID complies with log.Listener
func (*listenerV2) JobID() models.JobID {
	return models.NilJobID
}

// Job complies with log.Listener
func (lsn *listenerV2) JobIDV2() int32 {
	return lsn.job.ID
}

// IsV2Job complies with log.Listener
func (*listenerV2) IsV2Job() bool {
	return true
}
