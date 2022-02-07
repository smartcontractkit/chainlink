package vrf

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	heaps "github.com/theodesp/go-heaps"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/aggregator_v2v3_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ log.Listener = &listenerV2{}
	_ job.Service  = &listenerV2{}
)

const (
	// Gas used after computing the payment
	GasAfterPaymentCalculation = 21000 + // base cost of the transaction
		100 + 5000 + // warm subscription balance read and update. See https://eips.ethereum.org/EIPS/eip-2929
		2*2100 + 20000 - // cold read oracle address and oracle balance and first time oracle balance update, note first time will be 20k, but 5k subsequently
		4800 + // request delete refund (refunds happen after execution), note pre-london fork was 15k. See https://eips.ethereum.org/EIPS/eip-3529
		6685 // Positive static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.
)

type pendingRequest struct {
	confirmedAtBlock uint64
	req              *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	lb               log.Broadcast
	utcTimestamp     time.Time
}

type listenerV2 struct {
	utils.StartStopOnce
	cfg            Config
	l              logger.Logger
	ethClient      evmclient.Client
	logBroadcaster log.Broadcaster
	txm            bulletprooftxmanager.TxManager
	coordinator    *vrf_coordinator_v2.VRFCoordinatorV2
	pipelineRunner pipeline.Runner
	job            job.Job
	q              pg.Q
	gethks         keystore.Eth
	reqLogs        *utils.Mailbox
	chStop         chan struct{}
	// We can keep these pending logs in memory because we
	// only mark them confirmed once we send a corresponding fulfillment transaction.
	// So on node restart in the middle of processing, the lb will resend them.
	reqsMu   sync.Mutex // Both the log listener and the request handler write to reqs
	reqs     []pendingRequest
	reqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	respCountMu sync.Mutex
	respCount   map[string]uint64
	// This auxiliary heap is used when we need to purge the
	// respCount map - we repeatedly want to remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	blockNumberToReqID *pairing.PairHeap

	// head tracking data structures
	headBroadcaster  httypes.HeadBroadcasterRegistry
	latestHeadMu     sync.RWMutex
	latestHeadNumber uint64

	// Wait group to wait on all goroutines to shut down.
	wg *sync.WaitGroup

	// aggregator client to get link/eth feed prices from chain.
	aggregator *aggregator_v2v3_interface.AggregatorV2V3Interface
}

func (lsn *listenerV2) Start() error {
	return lsn.StartOnce("VRFListenerV2", func() error {
		spec := job.LoadEnvConfigVarsVRF(lsn.cfg, *lsn.job.VRFSpec)

		unsubscribeLogs := lsn.logBroadcaster.Register(lsn, log.ListenerOpts{
			Contract: lsn.coordinator.Address(),
			ParseLog: lsn.coordinator.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(): {
					{
						log.Topic(spec.PublicKey.MustHash()),
					},
				},
			},
			// Do not specify min confirmations, as it varies from request to request.
		})

		latestHead, unsubscribeHeadBroadcaster := lsn.headBroadcaster.Subscribe(lsn)
		if latestHead != nil {
			lsn.setLatestHead(latestHead)
		}

		// Log listener gathers request logs
		lsn.wg.Add(1)
		go func() {
			lsn.runLogListener([]func(){unsubscribeLogs, unsubscribeHeadBroadcaster}, spec.MinIncomingConfirmations, lsn.wg)
		}()

		// Request handler periodically computes a set of logs which can be fulfilled.
		lsn.wg.Add(1)
		go func() {
			lsn.runRequestHandler(spec.PollPeriod, lsn.wg)
		}()
		return nil
	})
}

func (lsn *listenerV2) setLatestHead(head *evmtypes.Head) {
	lsn.latestHeadMu.Lock()
	defer lsn.latestHeadMu.Unlock()
	num := uint64(head.Number)
	if num > lsn.latestHeadNumber {
		lsn.latestHeadNumber = num
	}
}

// OnNewLongestChain is called by the head broadcaster when a new head is available.
func (lsn *listenerV2) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	lsn.setLatestHead(head)
}

func (lsn *listenerV2) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return uint64(lsn.latestHeadNumber)
}

// Returns all the confirmed logs from
// the pending queue by subscription
func (lsn *listenerV2) getAndRemoveConfirmedLogsBySub(latestHead uint64) map[uint64][]pendingRequest {
	lsn.reqsMu.Lock()
	defer lsn.reqsMu.Unlock()
	var toProcess = make(map[uint64][]pendingRequest)
	var toKeep []pendingRequest
	for i := 0; i < len(lsn.reqs); i++ {
		if r := lsn.reqs[i]; r.confirmedAtBlock <= latestHead {
			toProcess[r.req.SubId] = append(toProcess[r.req.SubId], r)
		} else {
			toKeep = append(toKeep, lsn.reqs[i])
		}
	}
	lsn.reqs = toKeep
	return toProcess
}

// Remove all entries 10000 blocks or older
// to avoid a memory leak.
func (lsn *listenerV2) pruneConfirmedRequestCounts() {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	min := lsn.blockNumberToReqID.FindMin()
	for min != nil {
		m := min.(fulfilledReqV2)
		if m.blockNumber > (lsn.getLatestHead() - 10000) {
			break
		}
		delete(lsn.respCount, m.reqID)
		lsn.blockNumberToReqID.DeleteMin()
		min = lsn.blockNumberToReqID.FindMin()
	}
}

// Determine a set of logs that are confirmed
// and the subscription has sufficient balance to fulfill,
// given a eth call with the max gas price.
// Note we have to consider the pending reqs already in the bptxm as already "spent" link,
// using a max link consumed in their metadata.
// A user will need a minBalance capable of fulfilling a single req at the max gas price or nothing will happen.
// This is acceptable as users can choose different keyhashes which have different max gas prices.
// Other variables which can change the bill amount between our eth call simulation and tx execution:
// - Link/eth price fluctuation
// - Falling back to BHS
// However the likelihood is vanishingly small as
// 1) the window between simulation and tx execution is tiny.
// 2) the max gas price provides a very large buffer most of the time.
// Its easier to optimistically assume it will go though and in the rare case of a reversion
// we simply retry TODO: follow up where if we see a fulfillment revert, return log to the queue.
func (lsn *listenerV2) processPendingVRFRequests() {
	confirmed := lsn.getAndRemoveConfirmedLogsBySub(lsn.getLatestHead())
	keys, err := lsn.gethks.SendingKeys()
	if err != nil {
		lsn.l.Errorw("Unable to read sending keys", "err", err)
		return
	}
	fromAddress := keys[0].Address
	if lsn.job.VRFSpec.FromAddress != nil {
		fromAddress = *lsn.job.VRFSpec.FromAddress
	}
	maxGasPriceWei := lsn.cfg.KeySpecificMaxGasPriceWei(fromAddress.Address())
	// TODO: also probably want to order these by request time so we service oldest first
	// Get subscription balance. Note that outside of this request handler, this can only decrease while there
	// are no pending requests
	if len(confirmed) == 0 {
		lsn.l.Infow("No pending requests", "maxGasPrice", maxGasPriceWei, "fromAddress", fromAddress.Address())
		return
	}
	for subID, reqs := range confirmed {
		sub, err := lsn.coordinator.GetSubscription(nil, subID)
		if err != nil {
			lsn.l.Errorw("Unable to read subscription balance", "err", err)
			return
		}
		startBalance := sub.Balance
		lsn.processRequestsPerSub(subID, fromAddress.Address(), startBalance, maxGasPriceWei, reqs)
	}
	lsn.pruneConfirmedRequestCounts()
}

// MaybeSubtractReservedLink figures out how much LINK is reserved for other VRF requests that
// have not been fully confirmed yet on-chain, and subtracts that from the given startBalance,
// and returns that value if there are no errors.
func MaybeSubtractReservedLink(l logger.Logger, q pg.Q, fromAddress common.Address, startBalance *big.Int, chainID, subID uint64) (*big.Int, error) {
	var reservedLink string
	err := q.Get(&reservedLink, `SELECT SUM(CAST(meta->>'MaxLink' AS NUMERIC(78, 0)))
				   FROM eth_txes
				   WHERE meta->>'MaxLink' IS NOT NULL
				   AND evm_chain_id = $1
				   AND CAST(meta->>'SubId' AS NUMERIC) = $2
				   AND state IN ('unconfirmed', 'unstarted', 'in_progress')
				   GROUP BY meta->>'SubId'`, chainID, subID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		l.Errorw("Could not get reserved link", "err", err)
		return startBalance, err
	}

	if reservedLink != "" {
		reservedLinkInt, success := big.NewInt(0).SetString(reservedLink, 10)
		if !success {
			l.Errorw("Error converting reserved link", "reservedLink", reservedLink)
			return startBalance, errors.New("unable to convert returned link")
		}
		// Subtract the reserved link
		return startBalance.Sub(startBalance, reservedLinkInt), nil
	}
	return startBalance, nil
}

type fulfilledReqV2 struct {
	blockNumber uint64
	reqID       string
}

func (a fulfilledReqV2) Compare(b heaps.Item) int {
	a1 := a
	a2 := b.(fulfilledReqV2)
	switch {
	case a1.blockNumber > a2.blockNumber:
		return 1
	case a1.blockNumber < a2.blockNumber:
		return -1
	default:
		return 0
	}
}

func (lsn *listenerV2) processRequestsPerSub(
	subID uint64,
	fromAddress common.Address,
	startBalance *big.Int,
	maxGasPriceWei *big.Int,
	reqs []pendingRequest,
) {
	startBalanceNoReserveLink, err := MaybeSubtractReservedLink(
		lsn.l, lsn.q, fromAddress, startBalance, lsn.ethClient.ChainID().Uint64(), subID)
	if err != nil {
		lsn.l.Errorw("Couldn't get reserved LINK for subscription", "sub", reqs[0].req.SubId)
		return
	}
	lggr := lsn.l.With(
		"subID", reqs[0].req.SubId,
		"maxGasPrice", maxGasPriceWei.String(),
		"reqs", len(reqs),
		"startBalance", startBalance.String(),
		"startBalanceNoReservedLink", startBalanceNoReserveLink.String(),
	)
	lggr.Infow("Processing requests for subscription")

	// Attempt to process every request, break if we run out of balance
	var processed = make(map[string]struct{})
	for _, req := range reqs {
		vrfRequest := req.req
		rlog := lggr.With(
			"reqID", vrfRequest.RequestId.String(),
			"txHash", vrfRequest.Raw.TxHash,
		)

		// This check to see if the log was consumed needs to be in the same
		// goroutine as the mark consumed to avoid processing duplicates.
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(req.lb)
		if err != nil {
			// Do not process, let lb resend it as a retry mechanism.
			rlog.Errorw("Could not determine if log was already consumed", "error", err)
			continue
		} else if consumed {
			processed[vrfRequest.RequestId.String()] = struct{}{}
			continue
		}

		// Check if we can ignore the request due to it's age.
		if time.Now().UTC().Sub(req.utcTimestamp) >= lsn.job.VRFSpec.RequestTimeout {
			rlog.Infow("Request too old, dropping it")
			lsn.markLogAsConsumed(req.lb)
			processed[vrfRequest.RequestId.String()] = struct{}{}
			continue
		}

		// Check if the vrf req has already been fulfilled
		// If so we just mark it completed
		callback, err := lsn.coordinator.GetCommitment(nil, vrfRequest.RequestId)
		if err != nil {
			rlog.Errorw("Unable to check if already fulfilled, processing anyways", "err", err)
		} else if utils.IsEmpty(callback[:]) {
			// If seedAndBlockNumber is zero then the response has been fulfilled
			// and we should skip it
			rlog.Infow("Request already fulfilled", "callback", callback)
			lsn.markLogAsConsumed(req.lb)
			processed[vrfRequest.RequestId.String()] = struct{}{}
			continue
		}
		// Run the pipeline to determine the max link that could be billed at maxGasPrice.
		// The ethcall will error if there is currently insufficient balance onchain.
		maxLink, run, payload, gaslimit, err := lsn.getMaxLinkForFulfillment(maxGasPriceWei, req)
		if err != nil {
			rlog.Warnw("Unable to get max link for fulfillment, skipping request", "err", err)
			continue
		}
		if startBalance.Cmp(maxLink) < 0 {
			// Insufficient funds, have to wait for a user top up
			// leave it unprocessed for now
			rlog.Infow("Insufficient link balance to fulfill a request, breaking", "maxLink", maxLink)
			break
		}
		rlog.Infow("Enqueuing fulfillment")
		// We have enough balance to service it, lets enqueue for bptxm
		err = lsn.q.Transaction(func(tx pg.Queryer) error {
			if err = lsn.pipelineRunner.InsertFinishedRun(&run, true, pg.WithQueryer(tx)); err != nil {
				return err
			}
			if err = lsn.logBroadcaster.MarkConsumed(req.lb, pg.WithQueryer(tx)); err != nil {
				return err
			}
			_, err = lsn.txm.CreateEthTransaction(bulletprooftxmanager.NewTx{
				FromAddress:    fromAddress,
				ToAddress:      lsn.coordinator.Address(),
				EncodedPayload: hexutil.MustDecode(payload),
				GasLimit:       gaslimit,
				Meta: &bulletprooftxmanager.EthTxMeta{
					RequestID: common.BytesToHash(vrfRequest.RequestId.Bytes()),
					MaxLink:   maxLink.String(),
					SubID:     vrfRequest.SubId,
				},
				MinConfirmations: null.Uint32From(uint32(lsn.cfg.MinRequiredOutgoingConfirmations())),
				Strategy:         bulletprooftxmanager.NewSendEveryStrategy(),
				Checker: bulletprooftxmanager.TransmitCheckerSpec{
					CheckerType:           bulletprooftxmanager.TransmitCheckerTypeVRFV2,
					VRFCoordinatorAddress: lsn.coordinator.Address(),
				},
			}, pg.WithQueryer(tx))
			return err
		})
		if err != nil {
			rlog.Errorw("Error enqueuing fulfillment, requeuing request", "err", err)
			continue
		}
		// If we successfully enqueued for the bptxm, subtract that balance
		// And loop to attempt to enqueue another fulfillment
		startBalanceNoReserveLink = startBalanceNoReserveLink.Sub(startBalanceNoReserveLink, maxLink)
		processed[vrfRequest.RequestId.String()] = struct{}{}
	}
	// Remove all the confirmed logs
	var toKeep []pendingRequest
	for _, req := range reqs {
		if _, ok := processed[req.req.RequestId.String()]; !ok {
			toKeep = append(toKeep, req)
		}
	}
	lsn.reqsMu.Lock()
	// There could be logs accumulated to this slice while request processor is running,
	// so we merged the new ones with the ones that need to be requeued.
	lsn.reqs = append(lsn.reqs, toKeep...)
	lsn.reqsMu.Unlock()
	lggr.Infow("Finished processing for sub",
		"total reqs", len(reqs),
		"total processed", len(processed),
		"total remaining", len(toKeep),
		"total unique", len(toRequestSet(reqs)),
	)
}

func (lsn *listenerV2) estimateJuelsNeeded(
	req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested,
	maxGasPriceWei *big.Int,
) (*big.Int, error) {
	// Don't use up too much time to get this info, it's not critical for operating vrf.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	weiPerUnitLink, err := lsn.aggregator.LatestAnswer(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, errors.Wrap(err, "get aggregator latestAnswer")
	}
	// NOTE: no need to sanity check this as this is for logging purposes only
	// and should not be used to determine whether a user has enough funds in actuality,
	// we should always simulate for that.
	juelsNeeded := EstimateFeeJuels(
		req.CallbackGasLimit,
		maxGasPriceWei,
		weiPerUnitLink,
	)
	return juelsNeeded, nil
}

// Here we use the pipeline to parse the log, generate a vrf response
// then simulate the transaction at the max gas price to determine its maximum link cost.
func (lsn *listenerV2) getMaxLinkForFulfillment(maxGasPriceWei *big.Int, req pendingRequest) (*big.Int, pipeline.Run, string, uint64, error) {
	// estimate how much juels are needed so that we can log it if the simulation fails.
	juelsNeeded, err := lsn.estimateJuelsNeeded(req.req, maxGasPriceWei)
	if err != nil {
		// not critical, just log and continue
		lsn.l.Debugw("unable to estimate juels needed for request, continuing anyway", "reqID", req.req.RequestId)
		juelsNeeded = big.NewInt(0)
	}
	var (
		maxLink  *big.Int
		payload  string
		gaslimit uint64
	)
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"databaseID":    lsn.job.ID,
			"externalJobID": lsn.job.ExternalJobID,
			"name":          lsn.job.Name.ValueOrZero(),
			"publicKey":     lsn.job.VRFSpec.PublicKey[:],
			"maxGasPrice":   maxGasPriceWei.String(),
		},
		"jobRun": map[string]interface{}{
			"logBlockHash":   req.req.Raw.BlockHash[:],
			"logBlockNumber": req.req.Raw.BlockNumber,
			"logTxHash":      req.req.Raw.TxHash,
			"logTopics":      req.req.Raw.Topics,
			"logData":        req.req.Raw.Data,
		},
	})
	run, trrs, err := lsn.pipelineRunner.ExecuteRun(context.Background(), *lsn.job.PipelineSpec, vars, lsn.l)
	if err != nil {
		lsn.l.Errorw("Failed executing run", "err", err)
		return maxLink, run, payload, gaslimit, err
	}
	// The call task will fail if there are insufficient funds
	if run.AllErrors.HasError() {
		lsn.l.Warnw("Simulation errored, possibly insufficient funds. Request will remain unprocessed until funds are available",
			"err", run.AllErrors.ToError(), "max gas price", maxGasPriceWei, "reqID", req.req.RequestId, "juelsNeeded", juelsNeeded)
		return maxLink, run, payload, gaslimit, errors.Wrap(run.AllErrors.ToError(), "simulation errored")
	}
	if len(trrs.FinalResult(lsn.l).Values) != 1 {
		lsn.l.Errorw("Unexpected number of outputs", "expectedNumOutputs", 1, "actualNumOutputs", len(trrs.FinalResult(lsn.l).Values))
		return maxLink, run, payload, gaslimit, errors.New("unexpected number of outputs")
	}
	// Run succeeded, we expect a byte array representing the billing amount
	b, ok := trrs.FinalResult(lsn.l).Values[0].([]uint8)
	if !ok {
		lsn.l.Errorw("Unexpected type, expected []uint8 final result")
		return maxLink, run, payload, gaslimit, errors.New("expected []uint8 final result")
	}
	maxLink = utils.HexToBig(hexutil.Encode(b)[2:])
	for _, trr := range trrs {
		if trr.Task.Type() == pipeline.TaskTypeVRFV2 {
			m := trr.Result.Value.(map[string]interface{})
			payload = m["output"].(string)
		}
		if trr.Task.Type() == pipeline.TaskTypeEstimateGasLimit {
			gaslimit = trr.Result.Value.(uint64)
		}
	}
	return maxLink, run, payload, gaslimit, nil
}

func (lsn *listenerV2) runRequestHandler(pollPeriod time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	tick := time.NewTicker(pollPeriod)
	defer tick.Stop()
	for {
		select {
		case <-lsn.chStop:
			return
		case <-tick.C:
			lsn.processPendingVRFRequests()
		}
	}
}

func (lsn *listenerV2) runLogListener(unsubscribes []func(), minConfs uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	lsn.l.Infow("Listening for run requests",
		"minConfs", minConfs)
	for {
		select {
		case <-lsn.chStop:
			for _, f := range unsubscribes {
				f()
			}
			return
		case <-lsn.reqLogs.Notify():
			// Process all the logs in the queue if one is added
			for {
				i, exists := lsn.reqLogs.Retrieve()
				if !exists {
					break
				}
				lb, ok := i.(log.Broadcast)
				if !ok {
					panic(fmt.Sprintf("VRFListenerV2: invariant violated, expected log.Broadcast got %T", i))
				}
				lsn.handleLog(lb, minConfs)
			}
		}
	}
}

func (lsn *listenerV2) getConfirmedAt(req *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested, nodeMinConfs uint32) uint64 {
	lsn.respCountMu.Lock()
	defer lsn.respCountMu.Unlock()
	// Take the max(nodeMinConfs, requestedConfs + requestedConfsDelay).
	// Add the requested confs delay if provided in the jobspec so that we avoid an edge case
	// where the primary and backup VRF v2 nodes submit a proof at the same time.
	minConfs := nodeMinConfs
	if uint32(req.MinimumRequestConfirmations)+uint32(lsn.job.VRFSpec.RequestedConfsDelay) > nodeMinConfs {
		minConfs = uint32(req.MinimumRequestConfirmations) + uint32(lsn.job.VRFSpec.RequestedConfsDelay)
	}
	newConfs := uint64(minConfs) * (1 << lsn.respCount[req.RequestId.String()])
	// We cap this at 200 because solidity only supports the most recent 256 blocks
	// in the contract so if it was older than that, fulfillments would start failing
	// without the blockhash store feeder. We use 200 to give the node plenty of time
	// to fulfill even on fast chains.
	if newConfs > 200 {
		newConfs = 200
	}
	if lsn.respCount[req.RequestId.String()] > 0 {
		lsn.l.Warnw("Duplicate request found after fulfillment, doubling incoming confirmations",
			"txHash", req.Raw.TxHash,
			"blockNumber", req.Raw.BlockNumber,
			"blockHash", req.Raw.BlockHash,
			"reqID", req.RequestId.String(),
			"newConfs", newConfs)
	}
	return req.Raw.BlockNumber + newConfs
}

func (lsn *listenerV2) handleLog(lb log.Broadcast, minConfs uint32) {
	if v, ok := lb.DecodedLog().(*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled); ok {
		lsn.l.Infow("Received fulfilled log", "reqID", v.RequestId, "success", v.Success)
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lb)
		if err != nil {
			lsn.l.Errorw("Could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
			return
		} else if consumed {
			return
		}
		lsn.respCountMu.Lock()
		lsn.respCount[v.RequestId.String()]++
		lsn.respCountMu.Unlock()
		lsn.blockNumberToReqID.Insert(fulfilledReqV2{
			blockNumber: v.Raw.BlockNumber,
			reqID:       v.RequestId.String(),
		})
		lsn.markLogAsConsumed(lb)
		return
	}

	req, err := lsn.coordinator.ParseRandomWordsRequested(lb.RawLog())
	if err != nil {
		lsn.l.Errorw("Failed to parse log", "err", err, "txHash", lb.RawLog().TxHash)
		consumed, err := lsn.logBroadcaster.WasAlreadyConsumed(lb)
		if err != nil {
			lsn.l.Errorw("Could not determine if log was already consumed", "error", err, "txHash", lb.RawLog().TxHash)
			return
		} else if consumed {
			return
		}
		lsn.markLogAsConsumed(lb)
		return
	}

	confirmedAt := lsn.getConfirmedAt(req, minConfs)
	lsn.l.Infow("VRFListenerV2: Received log request", "reqID", req.RequestId, "confirmedAt", confirmedAt, "subID", req.SubId, "sender", req.Sender)
	lsn.reqsMu.Lock()
	lsn.reqs = append(lsn.reqs, pendingRequest{
		confirmedAtBlock: confirmedAt,
		req:              req,
		lb:               lb,
		utcTimestamp:     time.Now().UTC(),
	})
	lsn.reqAdded()
	lsn.reqsMu.Unlock()
}

func (lsn *listenerV2) markLogAsConsumed(lb log.Broadcast) {
	err := lsn.logBroadcaster.MarkConsumed(lb)
	lsn.l.ErrorIf(err, fmt.Sprintf("Unable to mark log %v as consumed", lb.String()))
}

// Close complies with job.Service
func (lsn *listenerV2) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		// wait on the request handler, log listener, and head listener to stop
		lsn.wg.Wait()
		return nil
	})
}

func (lsn *listenerV2) HandleLog(lb log.Broadcast) {
	wasOverCapacity := lsn.reqLogs.Deliver(lb)
	if wasOverCapacity {
		lsn.l.Error("Log mailbox is over capacity - dropped the oldest log")
	}
}

// Job complies with log.Listener
func (lsn *listenerV2) JobID() int32 {
	return lsn.job.ID
}

func toRequestSet(reqs []pendingRequest) map[string]struct{} {
	s := map[string]struct{}{}
	for _, r := range reqs {
		s[r.req.RequestId.String()] = struct{}{}
	}
	return s
}

// GasProofVerification is an upper limit on the gas used for verifying the VRF proof on-chain.
// It can be used to estimate the amount of LINK needed to fulfill a request.
const GasProofVerification uint32 = 200_000

// EstimateFeeJuels estimates the amount of link needed to fulfill a request
// given the callback gas limit, the gas price, and the wei per unit link.
func EstimateFeeJuels(callbackGasLimit uint32, maxGasPriceWei, weiPerUnitLink *big.Int) *big.Int {
	maxGasUsed := big.NewInt(int64(callbackGasLimit + GasProofVerification))
	costWei := new(big.Float).SetInt(
		maxGasUsed.Mul(maxGasUsed, maxGasPriceWei),
	)
	costLink := costWei.Quo(
		costWei,
		new(big.Float).SetInt(weiPerUnitLink),
	)
	costJuelsFloat := costLink.Mul(
		costLink,
		big.NewFloat(1e18),
	)
	costJuels, _ := costJuelsFloat.Int(nil)
	return costJuels
}
