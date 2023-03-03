package chain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/ocr2keepers/pkg/chain/gethwrappers/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

var (
	keeperRegistryABI = mustGetABI(keeper_registry_wrapper2_0.KeeperRegistryABI)
)

type outStruct struct {
	ur  []types.UpkeepResult
	err error
}

// evmRegistryv2_0 implements types.Registry interface
type evmRegistryv2_0 struct {
	address  common.Address
	registry *keeper_registry_wrapper2_0.KeeperRegistryCaller
	client   types.EVMClient
}

// NewEVMRegistryV2_0 is the constructor of evmRegistryv2_0
func NewEVMRegistryV2_0(address common.Address, client types.EVMClient) (*evmRegistryv2_0, error) {
	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistryCaller(address, client)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}

	return &evmRegistryv2_0{
		address:  address,
		registry: registry,
		client:   client,
	}, nil
}

func (r *evmRegistryv2_0) GetActiveUpkeepIDs(ctx context.Context) ([]types.UpkeepIdentifier, error) {
	opts, err := r.buildCallOpts(ctx, BlockKey("0"))
	if err != nil {
		return nil, err
	}

	state, err := r.registry.GetState(opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get contract state at block number %d", opts.BlockNumber.Int64())
	}

	keys := make([]types.UpkeepIdentifier, 0)
	for int64(len(keys)) < state.State.NumUpkeeps.Int64() {
		startIndex := int64(len(keys))
		maxCount := state.State.NumUpkeeps.Int64() - int64(len(keys))

		if maxCount > ActiveUpkeepIDBatchSize {
			maxCount = ActiveUpkeepIDBatchSize
		}

		nextRawKeys, err := r.registry.GetActiveUpkeepIDs(opts, big.NewInt(startIndex), big.NewInt(maxCount))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get active upkeep IDs from index %d to %d (both inclusive)", startIndex, startIndex+maxCount-1)
		}

		nextKeys := make([]types.UpkeepIdentifier, len(nextRawKeys))
		for i, next := range nextRawKeys {
			nextKeys[i] = types.UpkeepIdentifier(next.String())
		}

		if len(nextKeys) == 0 {
			break
		}

		buffer := make([]types.UpkeepIdentifier, len(keys), len(keys)+len(nextKeys))
		copy(keys, buffer)

		keys = append(buffer, nextKeys...)
	}

	return keys, nil
}

func (r *evmRegistryv2_0) checkUpkeeps(ctx context.Context, keys []types.UpkeepKey) ([]types.UpkeepResult, error) {
	var (
		checkReqs    = make([]rpc.BatchElem, len(keys))
		checkResults = make([]*string, len(keys))
	)

	for i, key := range keys {
		block, upkeepId, err := key.BlockKeyAndUpkeepID()
		if err != nil {
			return nil, err
		}

		upkeepIdInt, ok := upkeepId.BigInt()
		if !ok {
			return nil, ErrUpkeepKeyNotParsable
		}

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			return nil, err
		}

		payload, err := keeperRegistryABI.Pack("checkUpkeep", upkeepIdInt)
		if err != nil {
			return nil, err
		}

		var result string
		checkReqs[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.address.Hex(),
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
			results[i], err = unmarshalCheckUpkeepResult(keys[i], *checkResults[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return results, err
}

func (r *evmRegistryv2_0) simulatePerformUpkeeps(ctx context.Context, checkResults []types.UpkeepResult) ([]types.UpkeepResult, error) {
	var (
		performReqs     = make([]rpc.BatchElem, 0, len(checkResults))
		performResults  = make([]*string, 0, len(checkResults))
		performToKeyIdx = make([]int, 0, len(checkResults))
	)

	for i, checkResult := range checkResults {
		if checkResult.State == types.NotEligible {
			continue
		}

		block, upkeepId, err := checkResult.Key.BlockKeyAndUpkeepID()
		if err != nil {
			return nil, err
		}

		upkeepIdInt, ok := upkeepId.BigInt()
		if !ok {
			return nil, ErrUpkeepKeyNotParsable
		}

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			return nil, err
		}

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := keeperRegistryABI.Pack("simulatePerformUpkeep", upkeepIdInt, checkResult.PerformData)
		if err != nil {
			return nil, err
		}

		var result string
		performReqs = append(performReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"to":   r.address.Hex(),
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
			simulatePerformSuccess, err := unmarshalPerformUpkeepSimulationResult(*performResults[i])
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

func (r *evmRegistryv2_0) check(ctx context.Context, keys []types.UpkeepKey, ch chan outStruct) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		ch <- outStruct{
			err: err,
		}
		return
	}

	upkeepResults, err = r.simulatePerformUpkeeps(ctx, upkeepResults)
	if err != nil {
		ch <- outStruct{
			err: err,
		}
		return
	}

	ch <- outStruct{
		ur: upkeepResults,
	}
}

func (r *evmRegistryv2_0) CheckUpkeep(ctx context.Context, keys ...types.UpkeepKey) (types.UpkeepResults, error) {
	chResult := make(chan outStruct, 1)
	go r.check(ctx, keys, chResult)

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

func (r *evmRegistryv2_0) buildCallOpts(ctx context.Context, block types.BlockKey) (*bind.CallOpts, error) {
	b := new(big.Int)
	_, ok := b.SetString(block.String(), 10)

	if !ok {
		return nil, fmt.Errorf("%w: requires big int", ErrBlockKeyNotParsable)
	}

	if b == nil || b.Int64() == 0 {
		// fetch the current block number so batched GetLatestActiveUpkeepKeys calls can be performed on the same block
		header, err := r.client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: EVM failed to fetch block header", err, ErrRegistryCallFailure)
		}

		b = header.Number
	}

	return &bind.CallOpts{
		Context:     ctx,
		BlockNumber: b,
	}, nil
}
