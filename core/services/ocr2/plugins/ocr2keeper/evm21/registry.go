package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// defaultAllowListExpiration decides how long an upkeep's allow list info will be valid for.
	defaultAllowListExpiration = 20 * time.Minute
	// allowListCleanupInterval decides when the expired items in allowList cache will be deleted.
	allowListCleanupInterval = 5 * time.Minute
)

var (
	ErrLogReadFailure                = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable              = fmt.Errorf("head not available")
	ErrInitializationFailure         = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled              = fmt.Errorf("context was cancelled")
	ErrABINotParsable                = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize    int64 = 1000
	FetchUpkeepConfigBatchSize       = 10
	// This is the interval at which active upkeep list is fully refreshed from chain
	refreshInterval = 15 * time.Minute
	// This is the lookback for polling upkeep state event logs from latest block
	logEventLookback int64 = 250
)

//go:generate mockery --quiet --name Registry --output ./mocks/ --case=underscore
type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (encoding.UpkeepInfo, error)
	GetState(opts *bind.CallOpts) (iregistry21.GetState, error)
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	GetUpkeepTriggerConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	CheckCallback(opts *bind.TransactOpts, id *big.Int, values [][]byte, extraData []byte) (*coreTypes.Transaction, error)
	ParseLog(log coreTypes.Log) (generated.AbigenLog, error)
}

//go:generate mockery --quiet --name HttpClient --output ./mocks/ --case=underscore
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewEvmRegistry(
	lggr logger.Logger,
	addr common.Address,
	client evm.Chain,
	feedLookupCompatibleABI, keeperRegistryABI abi.ABI,
	registry *iregistry21.IKeeperRegistryMaster,
	mc *models.MercuryCredentials,
	al ActiveUpkeepList,
	logEventProvider logprovider.LogEventProvider,
	packer encoding.Packer,
	blockSub *BlockSubscriber,
) *EvmRegistry {
	return &EvmRegistry{
		lggr:         lggr.Named("EvmRegistry"),
		ht:           client.HeadTracker(),
		poller:       client.LogPoller(),
		addr:         addr,
		client:       client.Client(),
		logProcessed: make(map[string]bool),
		registry:     registry,
		abi:          keeperRegistryABI,
		active:       al,
		packer:       packer,
		headFunc:     func(ocr2keepers.BlockKey) {},
		chLog:        make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred:           mc,
			abi:            feedLookupCompatibleABI,
			allowListCache: cache.New(defaultAllowListExpiration, allowListCleanupInterval),
		},
		// TODO: is there some core node http client we should use
		hc:               http.DefaultClient,
		logEventProvider: logEventProvider,
		bs:               blockSub,
	}
}

var upkeepStateEvents = []common.Hash{
	iregistry21.IKeeperRegistryMasterUpkeepRegistered{}.Topic(),       // adds new upkeep id to registry
	iregistry21.IKeeperRegistryMasterUpkeepReceived{}.Topic(),         // adds new upkeep id to registry via migration
	iregistry21.IKeeperRegistryMasterUpkeepUnpaused{}.Topic(),         // unpauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepPaused{}.Topic(),           // pauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepMigrated{}.Topic(),         // migrated an upkeep, equivalent to cancel from this registry's perspective
	iregistry21.IKeeperRegistryMasterUpkeepCanceled{}.Topic(),         // cancels an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{}.Topic(), // trigger config was changed
}

type MercuryConfig struct {
	cred *models.MercuryCredentials
	abi  abi.ABI
	// allowListCache stores the upkeeps privileges. In 2.1, this only includes a JSON bytes for allowed to use mercury
	allowListCache *cache.Cache
}

type EvmRegistry struct {
	ht               types.HeadTracker
	sync             utils.StartStopOnce
	lggr             logger.Logger
	poller           logpoller.LogPoller
	addr             common.Address
	client           client.Client
	chainID          uint64
	registry         Registry
	abi              abi.ABI
	packer           encoding.Packer
	chLog            chan logpoller.Log
	reInit           *time.Timer
	mu               sync.RWMutex
	logProcessed     map[string]bool
	active           ActiveUpkeepList
	lastPollBlock    int64
	ctx              context.Context
	cancel           context.CancelFunc
	headFunc         func(ocr2keepers.BlockKey)
	runState         int
	runError         error
	mercury          *MercuryConfig
	hc               HttpClient
	bs               *BlockSubscriber
	logEventProvider logprovider.LogEventProvider
}

func (r *EvmRegistry) Name() string {
	return r.lggr.Name()
}

func (r *EvmRegistry) Start(ctx context.Context) error {
	return r.sync.StartOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.ctx, r.cancel = context.WithCancel(context.Background())
		r.reInit = time.NewTimer(refreshInterval)

		if err := r.registerEvents(r.chainID, r.addr); err != nil {
			return fmt.Errorf("logPoller error while registering automation events: %w", err)
		}

		// refresh the active upkeep keys; if the reInit timer returns, do it again
		{
			go func(cx context.Context, tmr *time.Timer, lggr logger.Logger, f func() error) {
				err := f()
				if err != nil {
					lggr.Errorf("failed to initialize upkeeps", err)
				}

				for {
					select {
					case <-tmr.C:
						err = f()
						if err != nil {
							lggr.Errorf("failed to re-initialize upkeeps", err)
						}
						tmr.Reset(refreshInterval)
					case <-cx.Done():
						return
					}
				}
			}(r.ctx, r.reInit, r.lggr, r.refreshActiveUpkeeps)
		}

		// start polling logs on an interval
		{
			go func(cx context.Context, lggr logger.Logger, f func() error) {
				ticker := time.NewTicker(time.Second)
				for {
					select {
					case <-ticker.C:
						err := f()
						if err != nil {
							lggr.Errorf("failed to poll logs for upkeeps", err)
						}
					case <-cx.Done():
						ticker.Stop()
						return
					}
				}
			}(r.ctx, r.lggr, r.pollUpkeepStateLogs)
		}

		// run process to process logs from log channel
		{
			go func(cx context.Context, ch chan logpoller.Log, lggr logger.Logger, f func(logpoller.Log) error) {
				for {
					select {
					case l := <-ch:
						err := f(l)
						if err != nil {
							lggr.Errorf("failed to process log for upkeep", err)
						}
					case <-cx.Done():
						return
					}
				}
			}(r.ctx, r.chLog, r.lggr, r.processUpkeepStateLog)
		}

		r.runState = 1
		return nil
	})
}

func (r *EvmRegistry) Close() error {
	return r.sync.StopOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.cancel()
		r.runState = 0
		r.runError = nil
		return nil
	})
}

func (r *EvmRegistry) Ready() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.runState == 1 {
		return nil
	}
	return r.sync.Ready()
}

func (r *EvmRegistry) HealthReport() map[string]error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.runState > 1 {
		r.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", r.runError))
	}
	return map[string]error{r.Name(): r.sync.Healthy()}
}

func (r *EvmRegistry) refreshActiveUpkeeps() error {
	// Allow for max timeout of refreshInterval
	ctx, cancel := context.WithTimeout(r.ctx, refreshInterval)
	defer cancel()

	r.lggr.Debugf("Refreshing active upkeeps list")
	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active upkeep ids from contract during refresh: %s", err)
	}
	r.active.Reset(ids...)

	// refresh trigger config log upkeeps
	// Note: We only update the trigger config for current active upkeeps and do not remove configs
	// for any upkeeps that might have been removed. We depend upon state events to remove the configs
	for _, id := range ids {
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(id)
		switch core.GetUpkeepType(*uid) {
		case ocr2keepers.LogTrigger:
			if err := r.updateTriggerConfig(id, nil); err != nil {
				r.lggr.Warnf("failed to update trigger config for upkeep ID %s: %s", id.String(), err)
			}
		default:
		}
	}

	return nil
}

func (r *EvmRegistry) pollUpkeepStateLogs() error {
	var latest int64
	var end int64
	var err error

	if end, err = r.poller.LatestBlock(pg.WithParentCtx(r.ctx)); err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	r.mu.Lock()
	latest = r.lastPollBlock
	r.lastPollBlock = end
	r.mu.Unlock()

	// if start and end are the same, no polling needs to be done
	if latest == 0 || latest == end {
		return nil
	}

	var logs []logpoller.Log
	if logs, err = r.poller.LogsWithSigs(
		end-logEventLookback,
		end,
		upkeepStateEvents,
		r.addr,
		pg.WithParentCtx(r.ctx),
	); err != nil {
		return fmt.Errorf("%w: %s", ErrLogReadFailure, err)
	}

	for _, log := range logs {
		r.chLog <- log
	}

	return nil
}

func (r *EvmRegistry) processUpkeepStateLog(l logpoller.Log) error {
	lid := fmt.Sprintf("%s%d", l.TxHash.String(), l.LogIndex)
	r.mu.Lock()
	if _, ok := r.logProcessed[lid]; ok {
		r.mu.Unlock()
		return nil
	}
	r.logProcessed[lid] = true
	r.mu.Unlock()
	txHash := l.TxHash.String()

	rawLog := l.ToGethLog()
	abilog, err := r.registry.ParseLog(rawLog)
	if err != nil {
		return err
	}

	switch l := abilog.(type) {
	case *iregistry21.IKeeperRegistryMasterUpkeepPaused:
		r.lggr.Debugf("KeeperRegistryUpkeepPaused log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepCanceled:
		r.lggr.Debugf("KeeperRegistryUpkeepCanceled log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepMigrated:
		r.lggr.Debugf("KeeperRegistryMasterUpkeepMigrated log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet:
		r.lggr.Debugf("KeeperRegistryUpkeepTriggerConfigSet log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		if err := r.updateTriggerConfig(l.Id, l.TriggerConfig); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepTriggerConfigSet for upkeep ID %s: %s", l.Id.String(), err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepRegistered:
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(l.Id)
		trigger := core.GetUpkeepType(*uid)
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s (trigger=%d) in transaction %s", l.Id.String(), trigger, txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepRegistered for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepReceived for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id.String(), txHash)
		r.active.Add(l.Id)
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config upon KeeperRegistryMasterUpkeepUnpaused for upkeep ID %s: %s", err)
		}
	default:
		r.lggr.Debugf("Unknown log detected for log %+v in transaction %s", l, txHash)
	}

	return nil
}

func RegistryUpkeepFilterName(addr common.Address) string {
	return logpoller.FilterName("KeeperRegistry Events", addr.String())
}

func (r *EvmRegistry) registerEvents(chainID uint64, addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	return r.poller.RegisterFilter(logpoller.Filter{
		Name:      RegistryUpkeepFilterName(addr),
		EventSigs: upkeepStateEvents,
		Addresses: []common.Address{addr},
	})
}

// Removes an upkeepID from active list and unregisters the log filter for log upkeeps
func (r *EvmRegistry) removeFromActive(id *big.Int) {
	r.active.Remove(id)

	uid := &ocr2keepers.UpkeepIdentifier{}
	uid.FromBigInt(id)
	trigger := core.GetUpkeepType(*uid)
	switch trigger {
	case ocr2keepers.LogTrigger:
		if err := r.logEventProvider.UnregisterFilter(id); err != nil {
			r.lggr.Warnw("failed to unregister log filter", "upkeepID", id.String())
		}
		r.lggr.Debugw("unregistered log filter", "upkeepID", id.String())
	default:
	}
}

func (r *EvmRegistry) buildCallOpts(ctx context.Context, block *big.Int) *bind.CallOpts {
	opts := bind.CallOpts{
		Context: ctx,
	}

	if block == nil || block.Int64() == 0 {
		l := r.ht.LatestChain()
		if l != nil && l.BlockNumber() != 0 {
			opts.BlockNumber = big.NewInt(l.BlockNumber())
		}
	} else {
		opts.BlockNumber = block
	}

	return &opts
}

func (r *EvmRegistry) getLatestIDsFromContract(ctx context.Context) ([]*big.Int, error) {
	opts := r.buildCallOpts(ctx, nil)

	state, err := r.registry.GetState(opts)
	if err != nil {
		n := "latest"
		if opts.BlockNumber != nil {
			n = fmt.Sprintf("%d", opts.BlockNumber.Int64())
		}

		return nil, fmt.Errorf("%w: failed to get contract state at block number '%s'", err, n)
	}

	ids := make([]*big.Int, 0, int(state.State.NumUpkeeps.Int64()))
	for int64(len(ids)) < state.State.NumUpkeeps.Int64() {
		startIndex := int64(len(ids))
		maxCount := state.State.NumUpkeeps.Int64() - startIndex

		if maxCount == 0 {
			break
		}

		if maxCount > ActiveUpkeepIDBatchSize {
			maxCount = ActiveUpkeepIDBatchSize
		}

		batchIDs, err := r.registry.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
		if err != nil {
			return nil, fmt.Errorf("%w: failed to get active upkeep IDs from index %d to %d (both inclusive)", err, startIndex, startIndex+maxCount-1)
		}

		ids = append(ids, batchIDs...)
	}

	return ids, nil
}

func (r *EvmRegistry) updateTriggerConfig(id *big.Int, cfg []byte) error {
	uid := &ocr2keepers.UpkeepIdentifier{}
	uid.FromBigInt(id)
	switch core.GetUpkeepType(*uid) {
	case ocr2keepers.LogTrigger:
		if len(cfg) == 0 {
			fetched, err := r.fetchTriggerConfig(id)
			if err != nil {
				return errors.Wrap(err, "failed to fetch log upkeep config")
			}
			cfg = fetched
		}
		parsed, err := r.packer.UnpackLogTriggerConfig(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to unpack log upkeep config")
		}
		if err := r.logEventProvider.RegisterFilter(id, logprovider.LogTriggerConfig(parsed)); err != nil {
			return errors.Wrap(err, "failed to register log filter")
		}
		r.lggr.Debugw("registered log filter", "upkeepID", id.String(), "cfg", parsed)
	default:
	}
	return nil
}

// updateTriggerConfig gets invoked upon changes in the trigger config of an upkeep.
func (r *EvmRegistry) fetchTriggerConfig(id *big.Int) ([]byte, error) {
	opts := r.buildCallOpts(r.ctx, nil)
	cfg, err := r.registry.GetUpkeepTriggerConfig(opts, id)
	if err != nil {
		r.lggr.Warnw("failed to get trigger config", "err", err)
		return nil, err
	}
	return cfg, nil
}

func (r *EvmRegistry) getBlockHash(blockNumber *big.Int) (common.Hash, error) {
	blocks, err := r.poller.GetBlocksRange(r.ctx, []uint64{blockNumber.Uint64()})
	if err != nil {
		return [32]byte{}, err
	}
	if len(blocks) == 0 {
		return [32]byte{}, fmt.Errorf("could not find block %d in log poller", blockNumber.Uint64())
	}

	return blocks[0].BlockHash, nil
}

func (r *EvmRegistry) getTxBlock(txHash common.Hash) (*big.Int, common.Hash, error) {
	// TODO: we need to differentiate here b/w flaky errors and tx not found errors
	txr, err := r.client.TransactionReceipt(r.ctx, txHash)
	if err != nil {
		return nil, common.Hash{}, err
	}

	return txr.BlockNumber, txr.BlockHash, nil
}

// getIneligibleCheckResultWithoutPerformData returns an ineligible check result with ineligibility reason and pipeline execution state but without perform data
func getIneligibleCheckResultWithoutPerformData(p ocr2keepers.UpkeepPayload, reason UpkeepFailureReason, state PipelineExecutionState, retryable bool) ocr2keepers.CheckResult {
	return ocr2keepers.CheckResult{
		IneligibilityReason:    uint8(reason),
		PipelineExecutionState: uint8(state),
		Retryable:              retryable,
		UpkeepID:               p.UpkeepID,
		Trigger:                p.Trigger,
		WorkID:                 p.WorkID,
		FastGasWei:             big.NewInt(0),
		LinkNative:             big.NewInt(0),
	}
}

// verifyCheckBlock checks that the check block and hash are valid, returns the pipeline execution state and retryable
func (r *EvmRegistry) verifyCheckBlock(ctx context.Context, checkBlock, upkeepId *big.Int, checkHash common.Hash) (state PipelineExecutionState, retryable bool) {
	// verify check block number is not too old
	if r.bs.latestBlock.Load()-checkBlock.Int64() > validCheckBlockRange {
		r.lggr.Warnf("latest block is %d, check block number %s is too old for upkeepId %s", r.bs.latestBlock.Load(), checkBlock, upkeepId)
		return CheckBlockTooOld, false
	}
	r.lggr.Warnf("latestBlock=%d checkBlock=%d", r.bs.latestBlock.Load(), checkBlock.Int64())

	var h string
	var ok bool
	// verify check block number and hash are valid
	h, ok = r.bs.queryBlocksMap(checkBlock.Int64())
	if !ok {
		r.lggr.Warnf("check block %s does not exist in block subscriber for upkeepId %s, querying eth client", checkBlock, upkeepId)
		hash, err := r.getBlockHash(checkBlock)
		if err != nil {
			r.lggr.Warnf("failed to query block %s: %s", checkBlock, err.Error())
			return RpcFlakyFailure, true
		}
		h = hash.Hex()
	}
	if checkHash.Hex() != h {
		r.lggr.Warnf("check block %s hash do not match. %s from block subscriber vs %s from trigger for upkeepId %s", checkBlock, h, checkHash.Hex(), upkeepId)
		return CheckBlockInvalid, false
	}
	return NoPipelineError, false
}

// verifyLogExists checks that the log still exists on chain, returns failure reason, pipeline error, and retryable
func (r *EvmRegistry) verifyLogExists(upkeepId *big.Int, p ocr2keepers.UpkeepPayload) (UpkeepFailureReason, PipelineExecutionState, bool) {
	logBlockNumber := int64(p.Trigger.LogTriggerExtension.BlockNumber)
	logBlockHash := common.BytesToHash(p.Trigger.LogTriggerExtension.BlockHash[:])
	// if log block number is populated, check log block number and block hash
	if logBlockNumber != 0 {
		h, ok := r.bs.queryBlocksMap(logBlockNumber)
		if ok && h == logBlockHash.Hex() {
			r.lggr.Debugf("tx hash %s exists on chain at block number %d for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), logBlockNumber, upkeepId)
			return UpkeepFailureReasonNone, NoPipelineError, false
		}
		r.lggr.Debugf("log block %d does not exist in block subscriber for upkeepId %s, querying eth client", logBlockNumber, upkeepId)
	} else {
		r.lggr.Debugf("log block not provided, querying eth client for tx hash %s for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
	}
	// query eth client as a fallback
	bn, _, err := r.getTxBlock(p.Trigger.LogTriggerExtension.TxHash)
	if err != nil {
		// primitive way of checking errors
		if strings.Contains(err.Error(), "missing required field") || strings.Contains(err.Error(), "not found") {
			return UpkeepFailureReasonTxHashNoLongerExists, NoPipelineError, false
		}
		r.lggr.Warnf("failed to query tx hash %s for upkeepId %s: %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId, err.Error())
		return UpkeepFailureReasonNone, RpcFlakyFailure, true
	}
	if bn == nil {
		r.lggr.Warnf("tx hash %s does not exist on chain for upkeepId %s.", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
		return UpkeepFailureReasonLogBlockNoLongerExists, NoPipelineError, false
	}
	r.lggr.Debugf("tx hash %s exists on chain for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
	return UpkeepFailureReasonNone, NoPipelineError, false
}
