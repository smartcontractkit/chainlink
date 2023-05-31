package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/multierr"
)

type upkeepState int32

const (
	stateActive upkeepState = iota
	stateInactive
)

type upkeepType int32

const (
	blockTrigger upkeepType = iota
	logTrigger
)

// upkeepInfoEntry holds the needed info of an upkeep
type upkeepInfoEntry struct {
	id *big.Int
	// state is used to determine if the upkeep is active or inactive
	state upkeepState
	// updatedAt is used to manage entry expiration
	updatedAt time.Time

	// target is used by mercury lookup to find the target contract
	target common.Address

	performGasLimit uint32
	// offchainConfig is the upkeep config.
	// used by log triggers to specify the corresponding filter
	offchainConfig []byte

	// TBD, for log triggers
	// balance *big.Int
}

// fetchUpkeep calls the contract to get the upkeep info for the given id, use this instead of fetchUpkeeps if you have a relatively.
func (r *EvmRegistry) fetchUpkeep(ctx context.Context, block *big.Int, id *big.Int) (upkeepInfoEntry, error) {
	if id == nil {
		return upkeepInfoEntry{}, fmt.Errorf("id is nil")
	}
	opts, err := r.buildCallOpts(ctx, block)
	if err != nil {
		return upkeepInfoEntry{}, fmt.Errorf("failed to build call opts: %w", err)
	}
	info, err := r.registry.GetUpkeep(opts, id)
	if err != nil {
		return upkeepInfoEntry{}, fmt.Errorf("failed to get upkeep info: %w", err)
	}
	upkeep := upkeepInfoEntry{
		id:              id,
		state:           stateActive,
		target:          info.Target,
		performGasLimit: info.ExecuteGas,
		offchainConfig:  info.OffchainConfig,
	}
	if info.Paused {
		upkeep.state = stateInactive
	}

	return upkeep, nil
}

// fetchUpkeeps fetches the upkeep info for the given ids, returns nil for missing/errored entries.
// It breaks the ids into batches.
func (r *EvmRegistry) fetchUpkeeps(ctx context.Context, block *big.Int, ids []*big.Int) ([]upkeepInfoEntry, error) {
	upkeeps := make([]upkeepInfoEntry, 0)

	var multiErr error
	var offset int
	for offset < len(ids) {
		batch := FetchUpkeepConfigBatchSize
		if len(ids)-offset < batch {
			batch = len(ids) - offset
		}

		currentBatch := ids[offset : offset+batch]
		if len(currentBatch) == 0 {
			break
		}
		infos, err := r.fetchUpkeepsBatch(ctx, block, currentBatch)
		if err != nil {
			multierr.AppendInto(&multiErr, fmt.Errorf("failed to get configs for id batch (length '%d'): %s", batch, err))
			// TBD: return or continue?
			// return fmt.Errorf("failed to get configs for id batch (length '%d'): %s", batch, err)
		}
		upkeeps = append(upkeeps, infos...)

		offset += batch
	}

	return upkeeps, multiErr
}

// fetchUpkeepsBatch fetches the upkeep info for the given ids.
// NOTE: this function doesn't take into account batch size, use getUpkeeps instead.
func (r *EvmRegistry) fetchUpkeepsBatch(ctx context.Context, block *big.Int, ids []*big.Int) ([]upkeepInfoEntry, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var (
		uReqs    = make([]rpc.BatchElem, len(ids))
		uResults = make([]*string, len(ids))
	)

	for i, id := range ids {
		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			r.lggr.Warnw("failed to build call opts", "err", err, "id", id)
			uResults[i] = nil
			uReqs[i] = rpc.BatchElem{
				Error: fmt.Errorf("failed to build call opts: %w", err),
			}
			continue
		}

		payload, err := r.abi.Pack("getUpkeep", id)
		if err != nil {
			r.lggr.Warnw("failed to pack id with abi", "err", err, "id", id)
			uResults[i] = nil
			uReqs[i] = rpc.BatchElem{
				Error: fmt.Errorf("failed to pack id with abi: %w", err),
			}
			continue
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
		results  = make([]upkeepInfoEntry, len(ids))
	)

	for i, req := range uReqs {
		if req.Error != nil {
			r.lggr.Debugf("error encountered for config id %s with message '%s' in get config", ids[i], req.Error)
			multierr.AppendInto(&multiErr, req.Error)
			continue
		}
		res, err := r.packer.UnpackUpkeepInfo(ids[i], *uResults[i])
		if err != nil {
			multierr.AppendInto(&multiErr, fmt.Errorf("failed to unpack result: %s", err))
			continue
		}
		results[i] = res
	}

	return results, multiErr
}
