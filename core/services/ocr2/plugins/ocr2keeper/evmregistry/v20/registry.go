package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	coreTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/multierr"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// DefaultUpkeepExpiration decides how long an upkeep info will be valid for. after it expires, a getUpkeepInfo
	// call will be made to the registry to obtain the most recent upkeep info and refresh this cache.
	DefaultUpkeepExpiration = 10 * time.Minute
	// DefaultCooldownExpiration decides how long a Mercury upkeep will be put in cool down for the first time. within
	// 10 minutes, subsequent failures will result in double amount of cool down period.
	DefaultCooldownExpiration = 5 * time.Second
	// DefaultApiErrExpiration decides a running sum of total errors of an upkeep in this 10 minutes window. it is used
	// to decide how long the cool down period will be.
	DefaultApiErrExpiration = 10 * time.Minute
	// CleanupInterval decides when the expired items in cache will be deleted.
	CleanupInterval = 15 * time.Minute
)

var (
	ErrLogReadFailure                = errors.New("failure reading logs")
	ErrHeadNotAvailable              = errors.New("head not available")
	ErrRegistryCallFailure           = errors.New("registry chain call failure")
	ErrBlockKeyNotParsable           = errors.New("block identifier not parsable")
	ErrUpkeepKeyNotParsable          = errors.New("upkeep key not parsable")
	ErrInitializationFailure         = errors.New("failed to initialize registry")
	ErrContextCancelled              = errors.New("context was cancelled")
	ErrABINotParsable                = errors.New("error parsing abi")
	ActiveUpkeepIDBatchSize    int64 = 1000
	FetchUpkeepConfigBatchSize       = 50
	separator                        = "|"
	reInitializationDelay            = 15 * time.Minute
	logEventLookback           int64 = 250
)

type Registry interface {
	GetUpkeep(opts *bind.CallOpts, id *big.Int) (keeper_registry_wrapper2_0.UpkeepInfo, error)
	GetState(opts *bind.CallOpts) (keeper_registry_wrapper2_0.GetState, error)
	GetActiveUpkeepIDs(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	ParseLog(log coreTypes.Log) (generated.AbigenLog, error)
}

type LatestBlockGetter interface {
	LatestBlock() int64
}

func NewEVMRegistryService(addr common.Address, client legacyevm.Chain, lggr logger.Logger) (*EvmRegistry, error) {
	keeperRegistryABI, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrABINotParsable, err)
	}

	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistry(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	r := &EvmRegistry{
		HeadProvider: HeadProvider{
			ht:     client.HeadTracker(),
			hb:     client.HeadBroadcaster(),
			chHead: make(chan ocr2keepers.BlockKey, 1),
		},
		lggr:     lggr.Named("AutomationRegistry"),
		poller:   client.LogPoller(),
		addr:     addr,
		client:   client.Client(),
		txHashes: make(map[string]bool),
		registry: registry,
		abi:      keeperRegistryABI,
		active:   make(map[string]activeUpkeep),
		packer:   &evmRegistryPackerV2_0{abi: keeperRegistryABI},
		headFunc: func(ocr2keepers.BlockKey) {},
		chLog:    make(chan logpoller.Log, 1000),
		enc:      EVMAutomationEncoder20{},
	}

	r.stopCh = make(chan struct{})
	r.reInit = time.NewTimer(reInitializationDelay)

	ctx, cancel := r.stopCh.NewCtx()
	defer cancel()
	if err := r.registerEvents(ctx, client.ID().Uint64(), addr); err != nil {
		return nil, fmt.Errorf("logPoller error while registering automation events: %w", err)
	}

	return r, nil
}

var upkeepStateEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered{}.Topic(),  // adds new upkeep id to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived{}.Topic(),    // adds new upkeep id to registry via migration
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet{}.Topic(), // unpauses an upkeep
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused{}.Topic(),    // updates the gas limit for an upkeep
}

var upkeepActiveEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformed{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryReorgedUpkeepReport{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryInsufficientFundsUpkeepReport{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryStaleUpkeepReport{}.Topic(),
}

type checkResult struct {
	ur  []EVMAutomationUpkeepResult20
	err error
}

type activeUpkeep struct {
	ID              *big.Int
	PerformGasLimit uint32
	CheckData       []byte
}

type EvmRegistry struct {
	HeadProvider
	sync          services.StateMachine
	lggr          logger.Logger
	poller        logpoller.LogPoller
	addr          common.Address
	client        client.Client
	registry      Registry
	abi           abi.ABI
	packer        *evmRegistryPackerV2_0
	chLog         chan logpoller.Log
	reInit        *time.Timer
	mu            sync.RWMutex
	txHashes      map[string]bool
	lastPollBlock int64
	stopCh        services.StopChan
	active        map[string]activeUpkeep
	headFunc      func(ocr2keepers.BlockKey)
	runState      int
	runError      error
	enc           EVMAutomationEncoder20
}

// GetActiveUpkeepKeys uses the latest head and map of all active upkeeps to build a
// slice of upkeep keys.
func (r *EvmRegistry) GetActiveUpkeepIDs(context.Context) ([]ocr2keepers.UpkeepIdentifier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]ocr2keepers.UpkeepIdentifier, len(r.active))
	var i int
	for _, value := range r.active {
		keys[i] = ocr2keepers.UpkeepIdentifier(value.ID.String())
		i++
	}

	return keys, nil
}

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

func (r *EvmRegistry) Start(_ context.Context) error {
	return r.sync.StartOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		// initialize the upkeep keys; if the reInit timer returns, do it again
		{
			go func(tmr *time.Timer, lggr logger.Logger, f func(context.Context) error) {
				ctx, cancel := r.stopCh.NewCtx()
				defer cancel()
				err := f(ctx)
				if err != nil {
					lggr.Errorf("failed to initialize upkeeps; error %v", err)
				}

				for {
					select {
					case <-tmr.C:
						err = f(ctx)
						if err != nil {
							lggr.Errorf("failed to re-initialize upkeeps; error %v", err)
						}
						tmr.Reset(reInitializationDelay)
					case <-ctx.Done():
						return
					}
				}
			}(r.reInit, r.lggr, r.initialize)
		}

		// start polling logs on an interval
		{
			go func(lggr logger.Logger, f func(context.Context) error) {
				ctx, cancel := r.stopCh.NewCtx()
				defer cancel()
				ticker := time.NewTicker(time.Second)

				for {
					select {
					case <-ticker.C:
						err := f(ctx)
						if err != nil {
							lggr.Errorf("failed to poll logs for upkeeps; error %v", err)
						}
					case <-ctx.Done():
						ticker.Stop()
						return
					}
				}
			}(r.lggr, r.pollLogs)
		}

		// run process to process logs from log channel
		{
			go func(ch chan logpoller.Log, lggr logger.Logger, f func(context.Context, logpoller.Log) error) {
				ctx, cancel := r.stopCh.NewCtx()
				defer cancel()
				for {
					select {
					case l := <-ch:
						err := f(ctx, l)
						if err != nil {
							lggr.Errorf("failed to process log for upkeep; error %v", err)
						}
					case <-ctx.Done():
						return
					}
				}
			}(r.chLog, r.lggr, r.processUpkeepStateLog)
		}

		r.runState = 1
		return nil
	})
}

func (r *EvmRegistry) Close() error {
	return r.sync.StopOnce("AutomationRegistry", func() error {
		r.mu.Lock()
		defer r.mu.Unlock()
		close(r.stopCh)
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

func (r *EvmRegistry) initialize(ctx context.Context) error {
	startupCtx, cancel := context.WithTimeout(ctx, reInitializationDelay)
	defer cancel()

	idMap := make(map[string]activeUpkeep)

	r.lggr.Debugf("Re-initializing active upkeeps list")
	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(startupCtx)
	if err != nil {
		return fmt.Errorf("failed to get ids from contract: %w", err)
	}

	var offset int
	for offset < len(ids) {
		batch := FetchUpkeepConfigBatchSize
		if len(ids)-offset < batch {
			batch = len(ids) - offset
		}

		actives, err := r.getUpkeepConfigs(startupCtx, ids[offset:offset+batch])
		if err != nil {
			return fmt.Errorf("failed to get configs for id batch (length '%d'): %w", batch, err)
		}

		for _, active := range actives {
			idMap[active.ID.String()] = active
		}

		offset += batch

		// Do not bombard RPC will calls, wait a bit
		time.Sleep(100 * time.Millisecond)
	}

	r.mu.Lock()
	r.active = idMap
	r.mu.Unlock()

	return nil
}

func (r *EvmRegistry) pollLogs(ctx context.Context) error {
	var latest int64
	var end logpoller.Block
	var err error

	if end, err = r.poller.LatestBlock(ctx); err != nil {
		return fmt.Errorf("%w: %w", ErrHeadNotAvailable, err)
	}

	r.mu.Lock()
	latest = r.lastPollBlock
	r.lastPollBlock = end.BlockNumber
	r.mu.Unlock()

	// if start and end are the same, no polling needs to be done
	if latest == 0 || latest == end.BlockNumber {
		return nil
	}

	{
		var logs []logpoller.Log
		if logs, err = r.poller.LogsWithSigs(
			ctx,
			end.BlockNumber-logEventLookback,
			end.BlockNumber,
			upkeepStateEvents,
			r.addr,
		); err != nil {
			return fmt.Errorf("%w: %w", ErrLogReadFailure, err)
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

func (r *EvmRegistry) registerEvents(ctx context.Context, chainID uint64, addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	return r.poller.RegisterFilter(ctx, logpoller.Filter{
		Name:      UpkeepFilterName(addr),
		EventSigs: append(upkeepStateEvents, upkeepActiveEvents...),
		Addresses: []common.Address{addr},
	})
}

func (r *EvmRegistry) processUpkeepStateLog(ctx context.Context, l logpoller.Log) error {
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
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered:
		r.lggr.Debugf("KeeperRegistryUpkeepRegistered log detected for upkeep ID %s in transaction %s", l.Id, hash)
		r.addToActive(ctx, l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived:
		r.lggr.Debugf("KeeperRegistryUpkeepReceived log detected for upkeep ID %s in transaction %s", l.Id, hash)
		r.addToActive(ctx, l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused:
		r.lggr.Debugf("KeeperRegistryUpkeepUnpaused log detected for upkeep ID %s in transaction %s", l.Id, hash)
		r.addToActive(ctx, l.Id, false)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet:
		r.lggr.Debugf("KeeperRegistryUpkeepGasLimitSet log detected for upkeep ID %s in transaction %s", l.Id, hash)
		r.addToActive(ctx, l.Id, true)
	}

	return nil
}

func (r *EvmRegistry) addToActive(ctx context.Context, id *big.Int, force bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.active == nil {
		r.active = make(map[string]activeUpkeep)
	}

	if _, ok := r.active[id.String()]; !ok || force {
		actives, err := r.getUpkeepConfigs(ctx, []*big.Int{id})
		if err != nil {
			r.lggr.Errorf("failed to get upkeep configs during adding active upkeep: %v", err)
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
			n = strconv.FormatInt(opts.BlockNumber.Int64(), 10)
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

func (r *EvmRegistry) doCheck(ctx context.Context, _ bool, keys []ocr2keepers.UpkeepKey, chResult chan checkResult) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
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
		return nil, nil, errors.New("unsplittable key")
	}

	if block, ok = new(big.Int).SetString(parts[0], 10); !ok {
		return nil, nil, errors.New("could not get block from key")
	}

	if id, ok = new(big.Int).SetString(parts[1], 10); !ok {
		return nil, nil, errors.New("could not get id from key")
	}

	return block, id, nil
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) checkUpkeeps(ctx context.Context, keys []ocr2keepers.UpkeepKey) ([]EVMAutomationUpkeepResult20, error) {
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

		payload, err := r.abi.Pack("checkUpkeep", upkeepId)
		if err != nil {
			return nil, err
		}

		args := []interface{}{
			map[string]interface{}{
				"to":   r.addr.Hex(),
				"data": hexutil.Bytes(payload),
			},
		}

		if opts.BlockNumber != nil {
			args = append(args, hexutil.EncodeBig(opts.BlockNumber))
		}

		var result string
		checkReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args:   args,
			Result: &result,
		}

		checkResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, checkReqs); err != nil {
		return nil, err
	}

	var (
		multiErr error
		results  = make([]EVMAutomationUpkeepResult20, len(keys))
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
				return nil, err
			}
		}
	}

	return results, multiErr
}

// TODO (AUTO-2013): Have better error handling to not return nil results in case of partial errors
func (r *EvmRegistry) simulatePerformUpkeeps(ctx context.Context, checkResults []EVMAutomationUpkeepResult20) ([]EVMAutomationUpkeepResult20, error) {
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

		args := []interface{}{
			map[string]interface{}{
				"to":   r.addr.Hex(),
				"data": hexutil.Bytes(payload),
			},
		}

		if opts.BlockNumber != nil {
			args = append(args, hexutil.EncodeBig(opts.BlockNumber))
		}

		var result string
		performReqs = append(performReqs, rpc.BatchElem{
			Method: "eth_call",
			Args:   args,
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
			return nil, fmt.Errorf("failed to get call opts: %w", err)
		}

		payload, err := r.abi.Pack("getUpkeep", id)
		if err != nil {
			return nil, fmt.Errorf("failed to pack id with abi: %w", err)
		}

		args := []interface{}{
			map[string]interface{}{
				"to":   r.addr.Hex(),
				"data": hexutil.Bytes(payload),
			},
		}

		if opts.BlockNumber != nil {
			args = append(args, hexutil.EncodeBig(opts.BlockNumber))
		}

		var result string
		uReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args:   args,
			Result: &result,
		}

		uResults[i] = &result
	}

	if err := r.client.BatchCallContext(ctx, uReqs); err != nil {
		return nil, fmt.Errorf("rpc error: %w", err)
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
			var err error
			results[i], err = r.packer.UnpackUpkeepResult(ids[i], *uResults[i])
			if err != nil {
				return nil, fmt.Errorf("failed to unpack result: %w", err)
			}
		}
	}

	return results, multiErr
}
