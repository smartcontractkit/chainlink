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
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// DefaultAllowListExpiration decides how long an upkeep's allow list info will be valid for.
	DefaultAllowListExpiration = 20 * time.Minute
	// CleanupInterval decides when the expired items in cache will be deleted.
	CleanupInterval = 25 * time.Minute
)

var (
	ErrLogReadFailure                = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable              = fmt.Errorf("head not available")
	ErrRegistryCallFailure           = fmt.Errorf("registry chain call failure")
	ErrBlockKeyNotParsable           = fmt.Errorf("block identifier not parsable")
	ErrUpkeepKeyNotParsable          = fmt.Errorf("upkeep key not parsable")
	ErrInitializationFailure         = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled              = fmt.Errorf("context was cancelled")
	ErrABINotParsable                = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize    int64 = 1000
	FetchUpkeepConfigBatchSize       = 10
	separator                        = "|"
	reInitializationDelay            = 15 * time.Minute
	logEventLookback           int64 = 250
)

//go:generate mockery --quiet --name Registry --output ./mocks/ --case=underscore
type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (iregistry21.UpkeepInfo, error)
	GetState(opts *bind.CallOpts) (iregistry21.GetState, error)
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	GetActiveUpkeepIDsByType(opts *bind.CallOpts, startIndex *big.Int, endIndex *big.Int, trigger uint8) ([]*big.Int, error)
	GetUpkeepAdminOffchainConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
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

func NewEVMRegistryService(addr common.Address, client evm.Chain, mc *models.MercuryCredentials, lggr logger.Logger) (*EvmRegistry, error) {
	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	registry, err := iregistry21.NewIKeeperRegistryMaster(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	r := &EvmRegistry{
		HeadProvider: HeadProvider{
			ht:     client.HeadTracker(),
			hb:     client.HeadBroadcaster(),
			chHead: make(chan ocr2keepers.BlockKey, 1),
		},
		lggr:     lggr.Named("EvmRegistry"),
		poller:   client.LogPoller(),
		addr:     addr,
		client:   client.Client(),
		txHashes: make(map[string]bool),
		registry: registry,
		active:   make(map[string]activeUpkeep),
		abi:      keeperRegistryABI,
		packer:   &evmRegistryPackerV2_1{abi: keeperRegistryABI},
		headFunc: func(ocr2keepers.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		mercury: &MercuryConfig{
			cred:           mc,
			abi:            feedLookupCompatibleABI,
			allowListCache: cache.New(DefaultAllowListExpiration, CleanupInterval),
		},
		hc:  http.DefaultClient,
		enc: EVMAutomationEncoder21{},
	}

	if err := r.registerEvents(client.ID().Uint64(), addr); err != nil {
		return nil, fmt.Errorf("logPoller error while registering automation events: %w", err)
	}

	return r, nil
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
	ur  []EVMAutomationUpkeepResult21
	err error
}

type activeUpkeep struct {
	ID              *big.Int
	PerformGasLimit uint32
	CheckData       []byte
}

type MercuryConfig struct {
	cred           *models.MercuryCredentials
	abi            abi.ABI
	allowListCache *cache.Cache
}

type EvmRegistry struct {
	HeadProvider
	sync          utils.StartStopOnce
	lggr          logger.Logger
	poller        logpoller.LogPoller
	addr          common.Address
	client        client.Client
	registry      Registry
	abi           abi.ABI
	packer        *evmRegistryPackerV2_1
	chLog         chan logpoller.Log
	reInit        *time.Timer
	mu            sync.RWMutex
	txHashes      map[string]bool
	lastPollBlock int64
	ctx           context.Context
	cancel        context.CancelFunc
	active        map[string]activeUpkeep
	headFunc      func(ocr2keepers.BlockKey)
	runState      int
	runError      error
	mercury       *MercuryConfig
	hc            HttpClient
	enc           EVMAutomationEncoder21
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
		if len(triggers) == 0 {
			keys = append(keys, ocr2keepers.UpkeepIdentifier(value.ID.String()))
			continue
		}
		trigger := getUpkeepType(value.ID.Bytes())
		for _, t := range triggers {
			if trigger == upkeepType(t) {
				keys = append(keys, ocr2keepers.UpkeepIdentifier(value.ID.String()))
				break
			}
		}
	}

	return keys, nil
}

// TODO: should be called with ocr2keepers.UpkeepPayload
func (r *EvmRegistry) CheckUpkeep(ctx context.Context, mercuryEnabled bool, keys ...ocr2keepers.UpkeepKey) ([]ocr2keepers.UpkeepResult, error) {
	chResult := make(chan checkResult, 1)
	go r.doCheck(ctx, mercuryEnabled, keys, chResult)

	select {
	case rs := <-chResult:
		result := make([]ocr2keepers.UpkeepResult, len(rs.ur))
		for i := range rs.ur {
			result[i] = rs.ur[i]
		}
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
	return logpoller.FilterName("EvmRegistry - Upkeep events for", addr.String())
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
	case *iregistry21.IKeeperRegistryMasterUpkeepRegistered:
		trigger := getUpkeepType(ocr2keepers.UpkeepIdentifier(l.Id.Bytes()))
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s (trigger=%d) in transaction %s", l.Id.String(), trigger, hash)
		r.addToActive(l.Id, false)
	case *iregistry21.IKeeperRegistryMasterUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, false)
	case *iregistry21.IKeeperRegistryMasterUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, false)
	case *iregistry21.IKeeperRegistryMasterUpkeepGasLimitSet:
		r.lggr.Debugf("KeeperRegistryUpkeepGasLimitSet log detected for upkeep ID %s in transaction %s", l.Id.String(), hash)
		r.addToActive(l.Id, true)
	default:
		r.lggr.Debugf("Unknown log detected for log %+v in transaction %s", l, hash)
	}

	return nil
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
			r.lggr.Errorf("failed to get upkeep configs during adding active upkeep: %w", err)
			return
		}

		if len(actives) != 1 {
			return
		}

		r.active[id.String()] = actives[0]
	}
}

func (r *EvmRegistry) buildCallOpts(ctx context.Context, block *big.Int) (*bind.CallOpts, error) {
	opts := bind.CallOpts{
		Context:     ctx,
		BlockNumber: nil,
	}

	if block == nil || block.Int64() == 0 {
		if r.LatestBlock() != 0 {
			opts.BlockNumber = big.NewInt(r.LatestBlock())
		}
	} else {
		opts.BlockNumber = block
	}

	return &opts, nil
}

func (r *EvmRegistry) getLatestIDsFromContract(ctx context.Context) ([]*big.Int, error) {
	opts, err := r.buildCallOpts(ctx, nil)
	if err != nil {
		return nil, err
	}

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

func (r *EvmRegistry) doCheck(ctx context.Context, mercuryEnabled bool, keys []ocr2keepers.UpkeepKey, chResult chan checkResult) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	if mercuryEnabled {
		if r.mercury.cred == nil || !r.mercury.cred.Validate() {
			chResult <- checkResult{
				err: errors.New("mercury credential is empty or not provided but FeedLookup feature is enabled on registry"),
			}
			return
		}
		upkeepResults, err = r.feedLookup(ctx, upkeepResults)
		if err != nil {
			chResult <- checkResult{
				err: err,
			}
			return
		}
	}

	upkeepResults, err = r.simulatePerformUpkeeps(ctx, upkeepResults)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	for i, res := range upkeepResults {
		r.mu.RLock()
		up, ok := r.active[res.ID.String()]
		r.mu.RUnlock()

		if ok {
			upkeepResults[i].ExecuteGas = up.PerformGasLimit
		}
	}

	chResult <- checkResult{
		ur: upkeepResults,
	}
}

func splitKey(key ocr2keepers.UpkeepKey) (*big.Int, *big.Int, error) {
	var (
		block *big.Int
		id    *big.Int
		ok    bool
	)

	parts := strings.Split(string(key), separator)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("unsplittable key")
	}

	if block, ok = new(big.Int).SetString(parts[0], 10); !ok {
		return nil, nil, fmt.Errorf("could not get block from key")
	}

	if id, ok = new(big.Int).SetString(parts[1], 10); !ok {
		return nil, nil, fmt.Errorf("could not get id from key")
	}

	return block, id, nil
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) checkUpkeeps(ctx context.Context, keys []ocr2keepers.UpkeepKey) ([]EVMAutomationUpkeepResult21, error) {
	var (
		checkReqs    = make([]rpc.BatchElem, len(keys))
		checkResults = make([]*string, len(keys))
	)

	for i, key := range keys {
		block, upkeepId, err := splitKey(key)
		if err != nil {
			return nil, err
		}

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			return nil, err
		}
		var payload []byte
		switch getUpkeepType(upkeepId.Bytes()) {
		case logTrigger:
			payload, err = r.abi.Pack("checkUpkeep", upkeepId, []byte{}) // TODO: pass log data
			if err != nil {
				return nil, err
			}
		default:
			payload, err = r.abi.Pack("checkUpkeep", upkeepId)
			if err != nil {
				return nil, err
			}
		}

		var result string
		checkReqs[i] = rpc.BatchElem{
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

		checkResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, checkReqs); err != nil {
		return nil, err
	}

	var (
		multiErr error
		results  = make([]EVMAutomationUpkeepResult21, len(keys))
	)

	for i, req := range checkReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for key %s with message '%s' in check", keys[i], req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			var err error
			r.lggr.Debugf("UnpackCheckResult key %s checkResult: %s", string(keys[i]), *checkResults[i])
			results[i], err = r.packer.UnpackCheckResult(keys[i], *checkResults[i])
			if err != nil {
				return nil, errors.Wrap(err, "failed to unpack check result")
			}
		}
	}

	return results, multiErr
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) simulatePerformUpkeeps(ctx context.Context, checkResults []EVMAutomationUpkeepResult21) ([]EVMAutomationUpkeepResult21, error) {
	var (
		performReqs     = make([]rpc.BatchElem, 0, len(checkResults))
		performResults  = make([]*string, 0, len(checkResults))
		performToKeyIdx = make([]int, 0, len(checkResults))
	)

	for i, checkResult := range checkResults {
		if !checkResult.Eligible {
			continue
		}

		opts, err := r.buildCallOpts(ctx, big.NewInt(int64(checkResult.Block)))
		if err != nil {
			return nil, err
		}

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := r.abi.Pack("simulatePerformUpkeep", checkResult.ID, checkResult.PerformData)
		if err != nil {
			return nil, err
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
			return nil, err
		}
	}

	var multiErr error

	for i, req := range performReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for key %d|%s with message '%s' in simulate perform", checkResults[i].Block, checkResults[i].ID, req.Error)
			multierr.AppendInto(&multiErr, req.Error)
		} else {
			simulatePerformSuccess, err := r.packer.UnpackPerformResult(*performResults[i])
			if err != nil {
				return nil, err
			}

			if !simulatePerformSuccess {
				checkResults[performToKeyIdx[i]].Eligible = false
			}
		}
	}

	return checkResults, multiErr
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
		opts, err := r.buildCallOpts(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get call opts: %s", err)
		}

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
				PerformGasLimit: info.ExecuteGas,
				CheckData:       info.CheckData,
			}
		}
	}

	return results, multiErr
}
