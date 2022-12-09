package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	ErrLogReadFailure              = fmt.Errorf("failure reading logs")
	ErrHeadNotAvailable            = fmt.Errorf("head not available")
	ErrRegistryCallFailure         = fmt.Errorf("registry chain call failure")
	ErrBlockKeyNotParsable         = fmt.Errorf("block identifier not parsable")
	ErrUpkeepKeyNotParsable        = fmt.Errorf("upkeep key not parsable")
	ErrInitializationFailure       = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled            = fmt.Errorf("context was cancelled")
	ErrABINotParsable              = fmt.Errorf("error parsing abi")
	ActiveUpkeepIDBatchSize  int64 = 1000
	separator                      = "|"
)

func NewEVMRegistryServiceV2_0(addr common.Address, client evm.Chain) (*EvmRegistry, error) {
	abi, err := abi.JSON(strings.NewReader(keeper_registry_wrapper2_0.KeeperRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistry(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	r := &EvmRegistry{
		HeadWatcher: HeadWatcher{
			client:  client.Client(),
			chReady: make(chan struct{}, 1),
		},
		poller:   client.LogPoller(),
		addr:     addr,
		client:   client.Client(),
		registry: registry,
		abi:      abi,
		packer:   &evmRegistryPackerV2_0{abi: abi},
		headFunc: func(types.BlockKey) {},
		active:   make(map[int64]activeUpkeep),
		chLog:    make(chan logpoller.Log, 1000),
	}

	if err := r.registerEvents(addr); err != nil {
		return nil, fmt.Errorf("logPoller error while registering automation events: %w", err)
	}

	return r, nil
}

var upkeepStateEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigrated{}.Topic(),   // removes upkeep id and detail from registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered{}.Topic(), // adds new upkeep id to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived{}.Topic(),   // adds multiple new upkeep ids to registry
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdated{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused{}.Topic(),
	keeper_registry_wrapper2_0.KeeperRegistryFundsAdded{}.Topic(),
}

var upkeepActiveEvents = []common.Hash{
	keeper_registry_wrapper2_0.KeeperRegistryUpkeepPerformed{}.Topic(),
}

type activeUpkeep struct {
	ID *big.Int
}

type checkResult struct {
	ur  []types.UpkeepResult
	err error
}

type EvmRegistry struct {
	HeadWatcher
	sync          utils.StartStopOnce
	mu            sync.RWMutex
	poller        logpoller.LogPoller
	filterID      int
	lastPollBlock int64
	addr          common.Address
	client        client.Client
	registry      *keeper_registry_wrapper2_0.KeeperRegistry
	abi           abi.ABI
	ctx           context.Context
	cancel        context.CancelFunc
	active        map[int64]activeUpkeep
	packer        *evmRegistryPackerV2_0
	headFunc      func(types.BlockKey)
	chLog         chan logpoller.Log
	runState      int
	runError      error
}

// GetActiveUpkeepKeys uses the latest head and map of all active upkeeps to build a
// slice of upkeep keys.
func (r *EvmRegistry) GetActiveUpkeepKeys(context.Context, types.BlockKey) ([]types.UpkeepKey, error) {
	if r.LatestBlock() == 0 {
		return nil, fmt.Errorf("%w: service probably not yet started", ErrHeadNotAvailable)
	}

	keys := make([]types.UpkeepKey, len(r.active))
	var i int
	for _, value := range r.active {
		keys[i] = blockAndIdToKey(big.NewInt(r.LatestBlock()), value.ID)
		i++
	}
	return keys, nil
}

func (r *EvmRegistry) CheckUpkeep(ctx context.Context, keys ...types.UpkeepKey) (types.UpkeepResults, error) {
	chResult := make(chan checkResult, 1)
	go r.doCheck(ctx, keys, chResult)

	select {
	case rs := <-chResult:
		return rs.ur, rs.err
	case <-ctx.Done():
		// safety on context done to provide an error on context cancellation
		// contract calls through the geth wrappers are a bit of a black box
		// so this safety net ensures contexts are fully respected and contract
		// call functions have a more graceful closure outside the scope of
		// CheckUpkeep needing to return immediately.
		return nil, fmt.Errorf("%w: failed to check upkeep on registry", ErrContextCancelled)
	}
}

func (r *EvmRegistry) IdentifierFromKey(key types.UpkeepKey) (types.UpkeepIdentifier, error) {
	_, id, err := blockAndIdFromKey(key)
	if err != nil {
		return nil, err
	}

	return id.Bytes(), nil
}

func (r *EvmRegistry) Start(ctx context.Context) error {
	return r.sync.StartOnce("AutomationRegistry", func() error {
		ctx, cancel := context.WithCancel(ctx)
		r.ctx = ctx
		r.cancel = cancel
		if err := r.initialize(); err != nil {
			return err
		}

		// start polling logs on an interval
		{
			go func(ctx context.Context, f func() error) {
				ticker := time.NewTicker(time.Second)

				for {
					select {
					case <-ticker.C:
						_ = f()
					case <-ctx.Done():
						ticker.Stop()
						return
					}
				}
			}(r.ctx, r.pollLogs)
		}

		if err := r.Watch(ctx); err != nil {
			return err
		}

		// run process to process logs from log channel
		{
			go func(ctx context.Context, ch chan logpoller.Log, f func(logpoller.Log) error) {
				for {
					select {
					case log := <-ch:
						_ = f(log)
					case <-ctx.Done():
						return
					}
				}
			}(r.ctx, r.chLog, r.processUpkeepStateLog)
		}

		r.runState = 1
		return nil
	})
}

func (r *EvmRegistry) Close() error {
	return r.sync.StopOnce("AutomationRegistry", func() error {
		r.cancel()
		r.runState = 0
		r.runError = nil
		if r.filterID > 0 {
			return r.poller.UnregisterFilter(r.filterID)
		}
		return nil
	})
}

func (r *EvmRegistry) Ready() error {
	if r.runState == 1 {
		return nil
	}
	return r.sync.Ready()
}

func (r *EvmRegistry) Healthy() error {
	if r.runState > 1 {
		return fmt.Errorf("failed run state: %w", r.runError)
	}
	return r.sync.Healthy()
}

func (r *EvmRegistry) initialize() error {
	startupCtx, cancel := context.WithTimeout(r.ctx, 10*time.Second)
	defer cancel()

	// get active upkeep ids from contract
	ids, err := r.getLatestIDsFromContract(startupCtx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		r.active[id.Int64()] = activeUpkeep{
			ID: id,
		}
	}

	return nil
}

func (r *EvmRegistry) pollLogs() error {
	var start int64
	var end int64
	var err error

	if end, err = r.poller.LatestBlock(); err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}

	r.mu.Lock()
	start = r.lastPollBlock
	r.lastPollBlock = end
	r.mu.Unlock()

	// if start and end are the same, no polling needs to be done
	if start == 0 || start == end {
		return nil
	}

	{
		var logs []logpoller.Log

		if logs, err = r.poller.LogsWithSigs(
			start,
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

	r.mu.Lock()
	r.lastPollBlock = end
	r.mu.Unlock()

	return nil
}

func (r *EvmRegistry) registerEvents(addr common.Address) error {
	// Add log filters for the log poller so that it can poll and find the logs that
	// we need
	filterID, err := r.poller.RegisterFilter(logpoller.Filter{
		EventSigs: append(upkeepStateEvents, upkeepActiveEvents...),
		Addresses: []common.Address{addr},
	})
	if err != nil {
		r.filterID = filterID
	}
	return err
}

func (r *EvmRegistry) processUpkeepStateLog(log logpoller.Log) error {
	rawLog := log.ToGethLog()
	abilog, err := r.registry.ParseLog(rawLog)
	if err != nil {
		return err
	}

	switch l := abilog.(type) {
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepMigrated: // removes upkeep id and detail from registry
		r.removeFromActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCanceled:
		r.removeFromActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepPaused:
		r.removeFromActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepRegistered: // adds new upkeep id to registry
		r.addToActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepReceived: // adds multiple new upkeep ids to registry
		r.addToActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepUnpaused:
		r.addToActive(l.Id)
	case *keeper_registry_wrapper2_0.KeeperRegistryFundsAdded:
		r.addToActive(l.Id)
		// case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepCheckDataUpdated:
		// case *keeper_registry_wrapper2_0.KeeperRegistryUpkeepGasLimitSet:
	}

	return nil
}

func (r *EvmRegistry) removeFromActive(id *big.Int) {
	if _, ok := r.active[id.Int64()]; ok {
		delete(r.active, id.Int64())
	}
}

func (r *EvmRegistry) addToActive(id *big.Int) {
	if _, ok := r.active[id.Int64()]; !ok {
		r.active[id.Int64()] = activeUpkeep{
			ID: id,
		}
	}
}

func (r *EvmRegistry) buildCallOpts(ctx context.Context, block *big.Int) (*bind.CallOpts, error) {
	opts := bind.CallOpts{
		Context: ctx,
	}

	if block == nil || block.Int64() == 0 {
		if r.LatestBlock() == 0 {
			// fetch the current block number so batched GetActiveUpkeepKeys calls can be performed on the same block
			bl, err := r.poller.LatestBlock()
			if err != nil {
				return nil, fmt.Errorf("%w: %s: EVM failed to fetch block header", err, ErrRegistryCallFailure)
			}

			block = new(big.Int).SetInt64(bl)
		} else {
			block = new(big.Int).SetInt64(r.LatestBlock())
		}
		opts.Pending = true
	} else if block.Int64() == r.LatestBlock() {
		opts.Pending = true
	} else {
		opts.BlockNumber = new(big.Int).Add(block, big.NewInt(1))
	}

	return &opts, nil
}

func (r *EvmRegistry) getLatestIDsFromContract(ctx context.Context) ([]*big.Int, error) {
	opts, err := r.buildCallOpts(ctx, nil)
	if err != nil {
		return nil, err
	}

	state, err := r.registry.KeeperRegistryCaller.GetState(opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get contract state at block number %d", opts.BlockNumber.Int64())
	}

	ids := make([]*big.Int, 0)
	for int64(len(ids)) < state.State.NumUpkeeps.Int64() {
		startIndex := int64(len(ids))
		maxCount := state.State.NumUpkeeps.Int64() - startIndex

		if maxCount > ActiveUpkeepIDBatchSize {
			maxCount = ActiveUpkeepIDBatchSize
		}

		batchIDs, err := r.registry.KeeperRegistryCaller.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get active upkeep IDs from index %d to %d (both inclusive)", startIndex, startIndex+maxCount-1)
		}

		if len(batchIDs) == 0 {
			break
		}

		buffer := make([]*big.Int, len(ids), len(ids)+len(batchIDs))
		copy(ids, buffer)

		ids = append(buffer, batchIDs...)
	}

	return ids, nil
}

func (r *EvmRegistry) doCheck(ctx context.Context, keys []types.UpkeepKey, chResult chan checkResult) {
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

	chResult <- checkResult{
		ur: upkeepResults,
	}
}

func (r *EvmRegistry) checkUpkeeps(ctx context.Context, keys []types.UpkeepKey) ([]types.UpkeepResult, error) {
	var (
		checkReqs    = make([]rpc.BatchElem, len(keys))
		checkResults = make([]*string, len(keys))
	)

	for i, key := range keys {
		block, upkeepId, err := blockAndIdFromKey(key)
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
		err     error
		results = make([]types.UpkeepResult, len(keys))
	)

	for i, req := range checkReqs {
		if req.Error != nil {
			if strings.Contains(req.Error.Error(), "reverted") {
				// subscription was canceled
				// NOTE: would we want to publish the fact that it is inactive?
				continue
			}
			// some other error
			multierr.AppendInto(&err, req.Error)
		} else {
			results[i], err = r.packer.UnpackCheckResult(keys[i], *checkResults[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return results, err
}

func (r *EvmRegistry) simulatePerformUpkeeps(ctx context.Context, checkResults []types.UpkeepResult) ([]types.UpkeepResult, error) {
	var (
		performReqs     = make([]rpc.BatchElem, 0, len(checkResults))
		performResults  = make([]*string, 0, len(checkResults))
		performToKeyIdx = make([]int, 0, len(checkResults))
	)

	for i, checkResult := range checkResults {
		if checkResult.State == types.NotEligible {
			continue
		}

		block, upkeepId, err := blockAndIdFromKey(checkResult.Key)
		if err != nil {
			return nil, err
		}

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			return nil, err
		}

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := r.abi.Pack("simulatePerformUpkeep", upkeepId, checkResult.PerformData)
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

	var err error

	for i, req := range performReqs {
		if req.Error != nil {
			if strings.Contains(req.Error.Error(), "reverted") {
				// subscription was canceled
				// NOTE: would we want to publish the fact that it is inactive?
				continue
			}
			// some other error
			multierr.AppendInto(&err, req.Error)
		} else {
			simulatePerformSuccess, err := r.packer.UnpackPerformResult(*performResults[i])
			if err != nil {
				return nil, err
			}

			if !simulatePerformSuccess {
				checkResults[performToKeyIdx[i]].State = types.NotEligible
			}
		}
	}

	return checkResults, nil
}

func blockAndIdToKey(block *big.Int, id *big.Int) types.UpkeepKey {
	return types.UpkeepKey(fmt.Sprintf("%s%s%s", block, separator, id))
}

func blockAndIdFromKey(key types.UpkeepKey) (*big.Int, *big.Int, error) {
	parts := strings.Split(string(key), separator)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("%w: missing data in upkeep key", ErrUpkeepKeyNotParsable)
	}

	block := new(big.Int)
	_, ok := block.SetString(parts[0], 10)
	if !ok {
		return nil, nil, fmt.Errorf("%w: must be big int", ErrUpkeepKeyNotParsable)
	}

	id := new(big.Int)
	_, ok = id.SetString(parts[1], 10)
	if !ok {
		return nil, nil, fmt.Errorf("%w: must be big int", ErrUpkeepKeyNotParsable)
	}

	return block, id, nil
}
