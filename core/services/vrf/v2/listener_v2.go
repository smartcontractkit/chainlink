package v2

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/batch_vrf_coordinator_v2plus"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_owner"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

var (
	_                         job.ServiceCtx = &listenerV2{}
	coordinatorV2ABI                         = evmtypes.MustGetABI(vrf_coordinator_v2.VRFCoordinatorV2ABI)
	coordinatorV2PlusABI                     = evmtypes.MustGetABI(vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalABI)
	batchCoordinatorV2ABI                    = evmtypes.MustGetABI(batch_vrf_coordinator_v2.BatchVRFCoordinatorV2ABI)
	batchCoordinatorV2PlusABI                = evmtypes.MustGetABI(batch_vrf_coordinator_v2plus.BatchVRFCoordinatorV2PlusABI)
	vrfOwnerABI                              = evmtypes.MustGetABI(vrf_owner.VRFOwnerMetaData.ABI)
	// These are the transaction states used when summing up already reserved subscription funds that are about to be used in in-flight transactions
	reserveEthLinkQueryStates = []txmgrtypes.TxState{txmgrcommon.TxUnconfirmed, txmgrcommon.TxUnstarted, txmgrcommon.TxInProgress}
)

const (
	// GasAfterPaymentCalculation is the gas used after computing the payment
	GasAfterPaymentCalculation = 21000 + // base cost of the transaction
		100 + 5000 + // warm subscription balance read and update. See https://eips.ethereum.org/EIPS/eip-2929
		2*2100 + 20000 - // cold read oracle address and oracle balance and first time oracle balance update, note first time will be 20k, but 5k subsequently
		4800 + // request delete refund (refunds happen after execution), note pre-london fork was 15k. See https://eips.ethereum.org/EIPS/eip-3529
		6685 // Positive static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.

	// BatchFulfillmentIterationGasCost is the cost of a single iteration of the batch coordinator's
	// loop. This is used to determine the gas allowance for a batch fulfillment call.
	BatchFulfillmentIterationGasCost = 52_000

	// backoffFactor is the factor by which to increase the delay each time a request fails.
	backoffFactor = 1.3

	txMetaFieldSubId  = "SubId"
	txMetaGlobalSubId = "GlobalSubId"
)

func New(
	cfg vrfcommon.Config,
	feeCfg vrfcommon.FeeConfig,
	l logger.Logger,
	chain legacyevm.Chain,
	chainID *big.Int,
	ds sqlutil.DataSource,
	coordinator CoordinatorV2_X,
	batchCoordinator batch_vrf_coordinator_v2.BatchVRFCoordinatorV2Interface,
	vrfOwner vrf_owner.VRFOwnerInterface,
	aggregator *aggregator_v3_interface.AggregatorV3Interface,
	pipelineRunner pipeline.Runner,
	gethks keystore.Eth,
	job job.Job,
	reqAdded func(),
	inflightCache vrfcommon.InflightCache,
	fulfillmentDeduper *vrfcommon.LogDeduper,
) job.ServiceCtx {
	return &listenerV2{
		cfg:                   cfg,
		feeCfg:                feeCfg,
		l:                     logger.Sugared(l),
		chain:                 chain,
		chainID:               chainID,
		coordinator:           coordinator,
		batchCoordinator:      batchCoordinator,
		vrfOwner:              vrfOwner,
		pipelineRunner:        pipelineRunner,
		job:                   job,
		ds:                    ds,
		gethks:                gethks,
		chStop:                make(chan struct{}),
		reqAdded:              reqAdded,
		blockNumberToReqID:    pairing.New(),
		latestHeadMu:          sync.RWMutex{},
		wg:                    &sync.WaitGroup{},
		aggregator:            aggregator,
		inflightCache:         inflightCache,
		fulfillmentLogDeduper: fulfillmentDeduper,
	}
}

type listenerV2 struct {
	services.StateMachine
	cfg     vrfcommon.Config
	feeCfg  vrfcommon.FeeConfig
	l       logger.SugaredLogger
	chain   legacyevm.Chain
	chainID *big.Int

	coordinator      CoordinatorV2_X
	batchCoordinator batch_vrf_coordinator_v2.BatchVRFCoordinatorV2Interface
	vrfOwner         vrf_owner.VRFOwnerInterface

	pipelineRunner pipeline.Runner
	job            job.Job
	ds             sqlutil.DataSource
	gethks         keystore.Eth
	chStop         services.StopChan

	reqAdded func() // A simple debug helper

	// Data structures for reorg attack protection
	// We want a map so we can do an O(1) count update every fulfillment log we get.
	respCount map[string]uint64
	// This auxiliary heap is used when we need to purge the
	// respCount map - we repeatedly want to remove the minimum log.
	// You could use a sorted list if the completed logs arrive in order, but they may not.
	blockNumberToReqID *pairing.PairHeap

	// head tracking data structures
	latestHeadMu     sync.RWMutex
	latestHeadNumber uint64

	// Wait group to wait on all goroutines to shut down.
	wg *sync.WaitGroup

	// aggregator client to get link/eth feed prices from chain. Can be nil for VRF V2 plus
	aggregator aggregator_v3_interface.AggregatorV3InterfaceInterface

	// fulfillmentLogDeduper prevents re-processing fulfillment logs.
	// fulfillment logs are used to increment counts in the respCount map
	// and to update the blockNumberToReqID heap.
	fulfillmentLogDeduper *vrfcommon.LogDeduper

	// inflightCache is a cache of in-flight requests, used to prevent
	// re-processing of requests that are in-flight or already fulfilled.
	inflightCache vrfcommon.InflightCache
}

func (lsn *listenerV2) HealthReport() map[string]error {
	return map[string]error{lsn.Name(): lsn.Healthy()}
}

func (lsn *listenerV2) Name() string { return lsn.l.Name() }

// Start starts listenerV2.
func (lsn *listenerV2) Start(ctx context.Context) error {
	return lsn.StartOnce("VRFListenerV2", func() error {
		// Check gas limit configuration
		confCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		conf, err := lsn.coordinator.GetConfig(&bind.CallOpts{Context: confCtx})
		gasLimit := lsn.feeCfg.LimitDefault()
		vrfLimit := lsn.feeCfg.LimitJobType().VRF()
		if vrfLimit != nil {
			gasLimit = uint64(*vrfLimit)
		}
		if err != nil {
			lsn.l.Criticalw("Error getting coordinator config for gas limit check, starting anyway.", "err", err)
		} else if uint64(conf.MaxGasLimit()+(GasProofVerification*2)) > gasLimit {
			lsn.l.Criticalw("Node gas limit setting may not be high enough to fulfill all requests; it should be increased. Starting anyway.",
				"currentGasLimit", gasLimit,
				"neededGasLimit", conf.MaxGasLimit()+(GasProofVerification*2),
				"callbackGasLimit", conf.MaxGasLimit(),
				"proofVerificationGas", GasProofVerification)
		}

		spec := job.LoadDefaultVRFPollPeriod(*lsn.job.VRFSpec)

		var respCount map[string]uint64
		respCount, err = lsn.GetStartingResponseCountsV2(ctx)
		if err != nil {
			return err
		}
		lsn.respCount = respCount

		if lsn.job.VRFSpec.CustomRevertsPipelineEnabled && lsn.vrfOwner != nil && lsn.job.VRFSpec.VRFOwnerAddress != nil {
			// Start reverted txns handler in background
			lsn.wg.Add(1)
			go func() {
				defer lsn.wg.Done()
				lsn.runRevertedTxnsHandler(spec.PollPeriod)
			}()
		}

		// Log listener gathers request logs and processes them
		lsn.wg.Add(1)
		go func() {
			defer lsn.wg.Done()
			lsn.runLogListener(spec.PollPeriod, spec.MinIncomingConfirmations)
		}()

		return nil
	})
}

func (lsn *listenerV2) GetStartingResponseCountsV2(ctx context.Context) (respCount map[string]uint64, err error) {
	respCounts := map[string]uint64{}
	var latestBlockNum *big.Int
	// Retry client call for LatestBlockHeight if fails
	// Want to avoid failing startup due to potential faulty RPC call
	err = retry.Do(func() error {
		latestBlockNum, err = lsn.chain.Client().LatestBlockHeight(ctx)
		return err
	}, retry.Attempts(10), retry.Delay(500*time.Millisecond))
	if err != nil {
		return nil, err
	}
	if latestBlockNum == nil {
		return nil, errors.New("LatestBlockHeight return nil block num")
	}
	confirmedBlockNum := latestBlockNum.Int64() - int64(lsn.chain.Config().EVM().FinalityDepth())
	// Only check as far back as the evm finality depth for completed transactions.
	var counts []vrfcommon.RespCountEntry
	counts, err = vrfcommon.GetRespCounts(ctx, lsn.chain.TxManager(), lsn.chainID, confirmedBlockNum)
	if err != nil {
		// Continue with an empty map, do not block job on this.
		lsn.l.Errorw("Unable to read previous confirmed fulfillments", "err", err)
		return respCounts, nil
	}

	for _, c := range counts {
		// Remove the quotes from the json
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		// Remove the 0x prefix
		b, err := hex.DecodeString(req[2:])
		if err != nil {
			lsn.l.Errorw("Unable to read fulfillment", "err", err, "reqID", c.RequestID)
			continue
		}
		bi := new(big.Int).SetBytes(b)
		respCounts[bi.String()] = uint64(c.Count)
	}
	return respCounts, nil
}

func (lsn *listenerV2) setLatestHead(head logpoller.Block) {
	lsn.latestHeadMu.Lock()
	defer lsn.latestHeadMu.Unlock()
	num := uint64(head.BlockNumber)
	if num > lsn.latestHeadNumber {
		lsn.latestHeadNumber = num
	}
}

func (lsn *listenerV2) getLatestHead() uint64 {
	lsn.latestHeadMu.RLock()
	defer lsn.latestHeadMu.RUnlock()
	return lsn.latestHeadNumber
}

// Close complies with job.Service
func (lsn *listenerV2) Close() error {
	return lsn.StopOnce("VRFListenerV2", func() error {
		close(lsn.chStop)
		// wait on the request handler, log listener
		lsn.wg.Wait()
		return nil
	})
}
