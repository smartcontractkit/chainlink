package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/patrickmn/go-cache"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// defaultAllowListExpiration decides how long an upkeep's allow list info will be valid for.
	defaultAllowListExpiration = 20 * time.Minute
	// cleanupInterval decides when the expired items in cache will be deleted.
	cleanupInterval = 25 * time.Minute
	// validCheckBlockRange decides the max distance between the check block and the current block
	validCheckBlockRange = 128
)

var (
	ErrLogReadFailure                = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable              = fmt.Errorf("head not available")
	ErrInitializationFailure         = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled              = fmt.Errorf("context was cancelled")
	ErrABINotParsable                = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize    int64 = 1000
	FetchUpkeepConfigBatchSize       = 10
	reInitializationDelay            = 15 * time.Minute
	logEventLookback           int64 = 250
)

//go:generate mockery --quiet --name Registry --output ./mocks/ --case=underscore
type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (UpkeepInfo, error)
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

type LatestBlockGetter interface {
	LatestBlock() int64
}

func NewEVMRegistryService(addr common.Address, client evm.Chain, mc *models.MercuryCredentials, bs *BlockSubscriber, lggr logger.Logger) (*EvmRegistry, *EVMAutomationEncoder21, error) {
	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	packer := NewEvmRegistryPackerV2_1(keeperRegistryABI, utilsABI)
	logPacker := logprovider.NewLogEventsPacker(utilsABI)

	registry, err := iregistry21.NewIKeeperRegistryMaster(addr, client.Client())
	if err != nil {
		return nil, nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	filterStore := logprovider.NewUpkeepFilterStore()
	logEventProvider := logprovider.New(lggr, client.LogPoller(), logPacker, filterStore, nil)

	r := &EvmRegistry{
		ht:       client.HeadTracker(),
		lggr:     lggr.Named("EvmRegistry"),
		poller:   client.LogPoller(),
		addr:     addr,
		client:   client.Client(),
		txHashes: make(map[string]bool),
		registry: registry,
		active:   make(map[string]activeUpkeep),
		abi:      keeperRegistryABI,
		packer:   packer,
		headFunc: func(ocr2keepers.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred:           mc,
			abi:            feedLookupCompatibleABI,
			allowListCache: cache.New(defaultAllowListExpiration, cleanupInterval),
		},
		hc:               http.DefaultClient,
		enc:              EVMAutomationEncoder21{packer: packer},
		logEventProvider: logEventProvider,
		bs:               bs,
	}

	if err := r.registerEvents(client.ID().Uint64(), addr); err != nil {
		return nil, nil, fmt.Errorf("logPoller error while registering automation events: %w", err)
	}

	return r, &r.enc, nil
}

var upkeepStateEvents = []common.Hash{
	iregistry21.IKeeperRegistryMasterUpkeepRegistered{}.Topic(),       // adds new upkeep id to registry
	iregistry21.IKeeperRegistryMasterUpkeepReceived{}.Topic(),         // adds new upkeep id to registry via migration
	iregistry21.IKeeperRegistryMasterUpkeepGasLimitSet{}.Topic(),      // unpauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepUnpaused{}.Topic(),         // updates the gas limit for an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepPaused{}.Topic(),           // pauses an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepCanceled{}.Topic(),         // cancels an upkeep
	iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet{}.Topic(), // trigger config was changed
}

var upkeepActiveEvents = []common.Hash{
	iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
	iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
	iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
	iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
}

type checkResult struct {
	cr  []ocr2keepers.CheckResult
	err error
}

type activeUpkeep struct {
	ID              *big.Int
	PerformGasLimit uint32
	CheckData       []byte
}

type MercuryConfig struct {
	cred *models.MercuryCredentials
	abi  abi.ABI
	// allowListCache stores the admin address' privilege. in 2.1, this only includes a JSON bytes for allowed to use mercury
	allowListCache *cache.Cache
}

type EvmRegistry struct {
	ht               types.HeadTracker
	sync             utils.StartStopOnce
	lggr             logger.Logger
	poller           logpoller.LogPoller
	addr             common.Address
	client           client.Client
	registry         Registry
	abi              abi.ABI
	packer           *evmRegistryPackerV2_1
	chLog            chan logpoller.Log
	reInit           *time.Timer
	mu               sync.RWMutex
	txHashes         map[string]bool
	lastPollBlock    int64
	ctx              context.Context
	cancel           context.CancelFunc
	active           map[string]activeUpkeep
	headFunc         func(ocr2keepers.BlockKey)
	runState         int
	runError         error
	mercury          *MercuryConfig
	hc               HttpClient
	enc              EVMAutomationEncoder21
	bs               *BlockSubscriber
	logEventProvider logprovider.LogEventProvider
}

// GetActiveUpkeepIDs uses the latest head and map of all active upkeeps to build a
// slice of upkeep keys.
func (r *EvmRegistry) GetActiveUpkeepIDs(ctx context.Context) ([]ocr2keepers.UpkeepIdentifier, error) {
	return r.GetActiveUpkeepIDsByType(ctx)
}

// GetActiveUpkeepIDsByType returns all active upkeeps of the given trigger types.
func (r *EvmRegistry) GetActiveUpkeepIDsByType(ctx context.Context, triggers ...uint8) ([]ocr2keepers.UpkeepIdentifier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]ocr2keepers.UpkeepIdentifier, 0)

	for _, value := range r.active {
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(value.ID)
		if len(triggers) == 0 {
			keys = append(keys, *uid)
			continue
		}
		trigger := core.GetUpkeepType(*uid)
		for _, t := range triggers {
			if trigger == ocr2keepers.UpkeepType(t) {
				keys = append(keys, *uid)
				break
			}
		}
	}

	return keys, nil
}

func (r *EvmRegistry) CheckUpkeeps(ctx context.Context, keys ...ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
	r.lggr.Debugw("Checking upkeeps", "upkeeps", keys)
	chResult := make(chan checkResult, 1)
	go r.doCheck(ctx, keys, chResult)

	select {
	case rs := <-chResult:
		result := make([]ocr2keepers.CheckResult, len(rs.cr))
		copy(result, rs.cr)
		return result, rs.err
	case <-ctx.Done():
		// safety on context done to provide an error on context cancellation
		// contract calls through the geth wrappers are a bit of a black box
		// so this safety net ensures contexts are fully respected and contract
		// call functions have a more graceful closure outside the scope of
		// CheckUpkeep needing to return immediately.
		return nil, fmt.Errorf("%w: failed to check upkeep on registry", ErrContextCancelled)
	}
}

func (r *EvmRegistry) Name() string {
	return r.lggr.Name()
}

func (r *EvmRegistry) Start(ctx context.Context) error {
	return r.sync.StartOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.ctx, r.cancel = context.WithCancel(context.Background())
		r.reInit = time.NewTimer(reInitializationDelay)

		// initialize the upkeep keys; if the reInit timer returns, do it again
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
						tmr.Reset(reInitializationDelay)
					case <-cx.Done():
						return
					}
				}
			}(r.ctx, r.reInit, r.lggr, r.initialize)
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
			}(r.ctx, r.lggr, r.pollLogs)
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

		// Start log event provider
		{
			go func(ctx context.Context, lggr logger.Logger, f func(context.Context) error, c func() error) {
				for ctx.Err() == nil {
					if err := f(ctx); err != nil {
						lggr.Errorf("failed to start log event provider", err)
					}
					if err := c(); err != nil {
						lggr.Errorf("failed to close log event provider", err)
					}
				}
			}(r.ctx, r.lggr, r.logEventProvider.Start, r.logEventProvider.Close)
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

func (r *EvmRegistry) LogEventProvider() logprovider.LogEventProvider {
	return r.logEventProvider
}

func (r *EvmRegistry) initialize() error {
	startupCtx, cancel := context.WithTimeout(r.ctx, reInitializationDelay)
	defer cancel()

	idMap := make(map[string]activeUpkeep)

	r.lggr.Debugf("Re-initializing active upkeeps list")
	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(startupCtx)
	if err != nil {
		return fmt.Errorf("failed to get ids from contract: %s", err)
	}

	var offset int
	for offset < len(ids) {
		batch := FetchUpkeepConfigBatchSize
		if len(ids)-offset < batch {
			batch = len(ids) - offset
		}

		actives, err := r.getUpkeepConfigs(startupCtx, ids[offset:offset+batch])
		if err != nil {
			return fmt.Errorf("failed to get configs for id batch (length '%d'): %s", batch, err)
		}

		for _, active := range actives {
			idMap[active.ID.String()] = active
		}

		offset += batch
	}

	r.mu.Lock()
	r.active = idMap
	r.mu.Unlock()

	// register upkeep ids for log triggers
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

func (r *EvmRegistry) pollLogs() error {
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

	{
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
	}

	return nil
}

func UpkeepFilterName(addr common.Address) string {
	return logpoller.FilterName("KeeperRegistry Events", addr.String())
}

func (r *EvmRegistry) registerEvents(chainID uint64, addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	return r.poller.RegisterFilter(logpoller.Filter{
		Name:      UpkeepFilterName(addr),
		EventSigs: append(upkeepStateEvents, upkeepActiveEvents...),
		Addresses: []common.Address{addr},
	})
}

func (r *EvmRegistry) processUpkeepStateLog(l logpoller.Log) error {
	hash := l.TxHash.String()
	if _, ok := r.txHashes[hash]; ok {
		return nil
	}
	r.txHashes[hash] = true

	rawLog := l.ToGethLog()
	abilog, err := r.registry.ParseLog(rawLog)
	if err != nil {
		return err
	}

	switch l := abilog.(type) {
	case *iregistry21.IKeeperRegistryMasterUpkeepPaused:
		r.lggr.Debugf("KeeperRegistryUpkeepPaused log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepCanceled:
		r.lggr.Debugf("KeeperRegistryUpkeepCanceled log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.removeFromActive(l.Id)
	case *iregistry21.IKeeperRegistryMasterUpkeepTriggerConfigSet:
		r.lggr.Debugf("KeeperRegistryUpkeepTriggerConfigSet log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		// passing nil instead of l.TriggerConfig to protect against reorgs,
		// as we'll fetch the latest config from the contract
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config for upkeep ID %s: %s", l.Id.String(), err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepRegistered:
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(l.Id)
		trigger := core.GetUpkeepType(*uid)
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s (trigger=%d) in transaction %s", l.Id.String(), trigger, hash)
		r.addToActive(l.Id, false)
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, false)
	case *iregistry21.IKeeperRegistryMasterUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, false)
		if err := r.updateTriggerConfig(l.Id, nil); err != nil {
			r.lggr.Warnf("failed to update trigger config for upkeep ID %s: %s", err)
		}
	case *iregistry21.IKeeperRegistryMasterUpkeepGasLimitSet:
		r.lggr.Debugf("KeeperRegistryUpkeepGasLimitSet log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, true)
	default:
		r.lggr.Debugf("Unknown log detected for log %+v in transaction %s", l, hash)
	}

	return nil
}

func (r *EvmRegistry) removeFromActive(id *big.Int) {
	r.mu.Lock()
	delete(r.active, id.String())
	r.mu.Unlock()

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

func (r *EvmRegistry) addToActive(id *big.Int, force bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.active == nil {
		r.active = make(map[string]activeUpkeep)
	}

	if _, ok := r.active[id.String()]; !ok || force {
		actives, err := r.getUpkeepConfigs(r.ctx, []*big.Int{id})
		if err != nil {
			r.lggr.Warnf("failed to get upkeep configs during adding active upkeep: %w", err)
			return
		}

		if len(actives) != 1 {
			return
		}

		r.active[id.String()] = actives[0]
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

func (r *EvmRegistry) doCheck(ctx context.Context, keys []ocr2keepers.UpkeepPayload, chResult chan checkResult) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	upkeepResults = r.feedLookup(ctx, upkeepResults)

	upkeepResults, err = r.simulatePerformUpkeeps(ctx, upkeepResults)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	chResult <- checkResult{
		cr: upkeepResults,
	}
}

// getBlockAndUpkeepId retrieves check block number, block hash from trigger and upkeep id
func (r *EvmRegistry) getBlockAndUpkeepId(upkeepID ocr2keepers.UpkeepIdentifier, trigger ocr2keepers.Trigger) (*big.Int, common.Hash, *big.Int) {
	block := new(big.Int).SetInt64(int64(trigger.BlockNumber))
	return block, common.BytesToHash(trigger.BlockHash[:]), upkeepID.BigInt()
}

func (r *EvmRegistry) checkUpkeeps(ctx context.Context, payloads []ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
	var (
		checkReqs    []rpc.BatchElem
		checkResults []*string
		results      = make([]ocr2keepers.CheckResult, len(payloads))
	)
	indices := map[int]int{}

	for i, p := range payloads {
		block, checkHash, upkeepId := r.getBlockAndUpkeepId(p.UpkeepID, p.Trigger)
		state, retryable := r.verifyCheckBlock(ctx, block, upkeepId, checkHash)
		if state != NoPipelineError {
			results[i] = getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, state, retryable)
			continue
		}

		opts := r.buildCallOpts(ctx, block)
		var payload []byte
		var err error
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(upkeepId)
		switch core.GetUpkeepType(*uid) {
		case ocr2keepers.LogTrigger:
			reason, state, retryable := r.verifyLogExists(upkeepId, p)
			if reason != UpkeepFailureReasonNone || state != NoPipelineError {
				results[i] = getIneligibleCheckResultWithoutPerformData(p, reason, state, retryable)
				continue
			}

			// check data will include the log trigger config
			payload, err = r.abi.Pack("checkUpkeep", upkeepId, p.CheckData)
			if err != nil {
				// pack error, no retryable
				r.lggr.Warnf("failed to pack log trigger checkUpkeep data for upkeepId %s with check data %s: %s", upkeepId, hexutil.Encode(p.CheckData), err)
				results[i] = getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false)
				continue
			}
		default:
			// checkUpkeep is overloaded on the contract for conditionals and log upkeeps
			// Need to use the first function (checkUpkeep0) for conditionals
			payload, err = r.abi.Pack("checkUpkeep0", upkeepId)
			if err != nil {
				// pack error, no retryable
				r.lggr.Warnf("failed to pack conditional checkUpkeep data for upkeepId %s with check data %s: %s", upkeepId, hexutil.Encode(p.CheckData), err)
				results[i] = getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, PackUnpackDecodeFailed, false)
				continue
			}
		}
		indices[len(checkReqs)] = i
		results[i] = getIneligibleCheckResultWithoutPerformData(p, UpkeepFailureReasonNone, NoPipelineError, false)

		var result string
		checkReqs = append(checkReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		})

		checkResults = append(checkResults, &result)
	}

	if len(checkResults) > 0 {
		// In contrast to CallContext, BatchCallContext only returns errors that have occurred
		// while sending the request. Any error specific to a request is reported through the
		// Error field of the corresponding BatchElem.
		// hence, if BatchCallContext returns an error, it will be an error which will terminate the pipeline
		if err := r.client.BatchCallContext(ctx, checkReqs); err != nil {
			r.lggr.Errorf("failed to batch call for checkUpkeeps: %s", err)
			return nil, err
		}
	}

	for i, req := range checkReqs {
		index := indices[i]
		if req.Error != nil {
			// individual upkeep failed in a batch call, retryable
			r.lggr.Warnf("error encountered in check result for upkeepId %s: %s", results[index].UpkeepID.String(), req.Error)
			results[index].Retryable = true
			results[index].PipelineExecutionState = uint8(RpcFlakyFailure)
		} else {
			var err error
			results[index], err = r.packer.UnpackCheckResult(payloads[index], *checkResults[i])
			if err != nil {
				r.lggr.Warnf("failed to unpack check result: %s", err)
			}
		}
	}

	return results, nil
}

func (r *EvmRegistry) simulatePerformUpkeeps(ctx context.Context, checkResults []ocr2keepers.CheckResult) ([]ocr2keepers.CheckResult, error) {
	var (
		performReqs     = make([]rpc.BatchElem, 0, len(checkResults))
		performResults  = make([]*string, 0, len(checkResults))
		performToKeyIdx = make([]int, 0, len(checkResults))
	)

	for i, cr := range checkResults {
		if !cr.Eligible {
			continue
		}

		block, _, upkeepId := r.getBlockAndUpkeepId(cr.UpkeepID, cr.Trigger)

		opts := r.buildCallOpts(ctx, block)

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := r.abi.Pack("simulatePerformUpkeep", upkeepId, cr.PerformData)
		if err != nil {
			// pack failed, not retryable
			r.lggr.Warnf("failed to pack perform data %s for %s: %s", hexutil.Encode(cr.PerformData), upkeepId, err)
			checkResults[i].Eligible = false
			checkResults[i].PipelineExecutionState = uint8(PackUnpackDecodeFailed)
			continue
		}

		var result string
		performReqs = append(performReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		})

		performResults = append(performResults, &result)
		performToKeyIdx = append(performToKeyIdx, i)
	}

	if len(performReqs) > 0 {
		if err := r.client.BatchCallContext(ctx, performReqs); err != nil {
			r.lggr.Errorf("failed to batch call for simulatePerformUpkeeps: %s", err)
			return nil, err
		}
	}

	for i, req := range performReqs {
		idx := performToKeyIdx[i]
		if req.Error != nil {
			// individual upkeep failed in a batch call, retryable
			r.lggr.Warnf("failed to simulate upkeepId %s: %s", checkResults[idx].UpkeepID.String(), req.Error)
			checkResults[idx].Retryable = true
			checkResults[idx].Eligible = false
			checkResults[idx].PipelineExecutionState = uint8(RpcFlakyFailure)
			continue
		}

		state, simulatePerformSuccess, err := r.packer.UnpackPerformResult(*performResults[i])
		if err != nil {
			// unpack failed, not retryable
			r.lggr.Warnf("failed to unpack simulate performUpkeep result for upkeepId %s for state %d: %s", checkResults[idx].UpkeepID.String(), state, req.Error)
			checkResults[idx].Retryable = false
			checkResults[idx].Eligible = false
			checkResults[idx].PipelineExecutionState = uint8(state)
			continue
		}

		if !simulatePerformSuccess {
			r.lggr.Warnf("upkeepId %s is not eligible after simulation", checkResults[idx].UpkeepID.String())
			checkResults[performToKeyIdx[i]].Eligible = false
		}
	}

	return checkResults, nil
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) getUpkeepConfigs(ctx context.Context, ids []*big.Int) ([]activeUpkeep, error) {
	if len(ids) == 0 {
		return []activeUpkeep{}, nil
	}

	var (
		uReqs    = make([]rpc.BatchElem, len(ids))
		uResults = make([]*string, len(ids))
	)

	for i, id := range ids {
		opts := r.buildCallOpts(ctx, nil)

		payload, err := r.abi.Pack("getUpkeep", id)
		if err != nil {
			return nil, fmt.Errorf("failed to pack id with abi: %s", err)
		}

		var result string
		uReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.addr.Hex(),
					"data": hexutil.Bytes(payload),
				},
				hexutil.EncodeBig(opts.BlockNumber),
			},
			Result: &result,
		}

		uResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, uReqs); err != nil {
		return nil, fmt.Errorf("rpc error: %s", err)
	}

	var (
		multiErr error
		results  = make([]activeUpkeep, len(ids))
	)

	for i, req := range uReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for config id %s with message '%s' in get config", ids[i], req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			info, err := r.packer.UnpackUpkeepInfo(ids[i], *uResults[i])
			if err != nil {
				return nil, fmt.Errorf("failed to unpack result: %s", err)
			}
			results[i] = activeUpkeep{ // TODO
				ID:              ids[i],
				PerformGasLimit: info.PerformGas,
				CheckData:       info.CheckData,
			}
		}
	}

	return results, multiErr
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
	block, err := r.client.BlockByNumber(r.ctx, blockNumber)
	if err != nil {
		return [32]byte{}, err
	}

	return block.Hash(), nil
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
		b, err := r.client.BlockByNumber(ctx, checkBlock)
		if err != nil {
			r.lggr.Warnf("failed to query block %s: %s", checkBlock, err.Error())
			return RpcFlakyFailure, true
		}
		h = b.Hash().Hex()
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
