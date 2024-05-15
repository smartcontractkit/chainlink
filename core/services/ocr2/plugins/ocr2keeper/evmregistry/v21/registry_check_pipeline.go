package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/gasprice"
)

const (
	// checkBlockTooOldRange is the number of blocks that can be behind the latest block before
	// we return a CheckBlockTooOld error
	checkBlockTooOldRange = 128
	zeroAddress           = "0x0000000000000000000000000000000000000000"
)

type checkResult struct {
	cr  []ocr2keepers.CheckResult
	err error
}

func (r *EvmRegistry) CheckUpkeeps(ctx context.Context, keys ...ocr2keepers.UpkeepPayload) ([]ocr2keepers.CheckResult, error) {
	r.lggr.Debugw("Checking upkeeps", "upkeeps", keys)
	for i := range keys {
		if keys[i].Trigger.BlockNumber == 0 { // check block was not populated, use latest
			latest := r.bs.latestBlock.Load()
			if latest == nil {
				return nil, fmt.Errorf("no latest block available")
			}
			copy(keys[i].Trigger.BlockHash[:], latest.Hash[:])
			keys[i].Trigger.BlockNumber = latest.Number
			r.lggr.Debugf("Check upkeep key had no trigger block number, using latest block %v", keys[i].Trigger.BlockNumber)
		}
	}

	chResult := make(chan checkResult, 1)

	r.threadCtrl.GoCtx(ctx, func(ctx context.Context) {
		r.doCheck(ctx, keys, chResult)
	})

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

func (r *EvmRegistry) doCheck(ctx context.Context, keys []ocr2keepers.UpkeepPayload, chResult chan checkResult) {
	upkeepResults, err := r.checkUpkeeps(ctx, keys)
	if err != nil {
		chResult <- checkResult{
			err: err,
		}
		return
	}

	upkeepResults = r.streams.Lookup(ctx, upkeepResults)

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

// verifyCheckBlock checks that the check block and hash are valid, returns the pipeline execution state and retryable
func (r *EvmRegistry) verifyCheckBlock(_ context.Context, checkBlock, upkeepId *big.Int, checkHash common.Hash) (state encoding.PipelineExecutionState, retryable bool) {
	// verify check block number and hash are valid
	h, ok := r.bs.queryBlocksMap(checkBlock.Int64())
	// if this block number/hash combo exists in block subscriber, this check block and hash still exist on chain and are valid
	// the block hash in block subscriber might be slightly outdated, if it doesn't match then we fetch the latest from RPC.
	if ok && h == checkHash.Hex() {
		r.lggr.Debugf("check block hash %s exists on chain at block number %d for upkeepId %s", checkHash.Hex(), checkBlock, upkeepId)
		return encoding.NoPipelineError, false
	}
	r.lggr.Warnf("check block %s does not exist in block subscriber or hash does not match for upkeepId %s. this may be caused by block subscriber outdated due to re-org, querying eth client to confirm", checkBlock, upkeepId)
	b, err := r.getBlockHash(checkBlock)
	if err != nil {
		r.lggr.Warnf("failed to query block %s: %s", checkBlock, err.Error())
		return encoding.RpcFlakyFailure, true
	}
	if checkHash.Hex() != b.Hex() {
		r.lggr.Warnf("check block %s hash do not match. %s from block subscriber vs %s from trigger for upkeepId %s", checkBlock, h, checkHash.Hex(), upkeepId)
		return encoding.CheckBlockInvalid, false
	}
	return encoding.NoPipelineError, false
}

// verifyLogExists checks that the log still exists on chain, returns failure reason, pipeline error, and retryable
func (r *EvmRegistry) verifyLogExists(upkeepId *big.Int, p ocr2keepers.UpkeepPayload) (encoding.UpkeepFailureReason, encoding.PipelineExecutionState, bool) {
	logBlockNumber := int64(p.Trigger.LogTriggerExtension.BlockNumber)
	logBlockHash := common.BytesToHash(p.Trigger.LogTriggerExtension.BlockHash[:])
	checkBlockHash := common.BytesToHash(p.Trigger.BlockHash[:])
	if checkBlockHash.String() == logBlockHash.String() {
		// log verification would be covered by checkBlock verification as they are the same. Return early from
		// log verificaion. This also helps in preventing some racy conditions when rpc does not return the tx receipt
		// for a very new log
		return encoding.UpkeepFailureReasonNone, encoding.NoPipelineError, false
	}
	// if log block number is populated, check log block number and block hash
	if logBlockNumber != 0 {
		h, ok := r.bs.queryBlocksMap(logBlockNumber)
		// if this block number/hash combo exists in block subscriber, this block and tx still exists on chain and is valid
		// the block hash in block subscriber might be slightly outdated, if it doesn't match then we fetch the latest from RPC.
		if ok && h == logBlockHash.Hex() {
			r.lggr.Debugf("tx hash %s exists on chain at block number %d, block hash %s for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), logBlockHash.Hex(), logBlockNumber, upkeepId)
			return encoding.UpkeepFailureReasonNone, encoding.NoPipelineError, false
		}
		// if this block does not exist in the block subscriber, the block which this log lived on was probably re-orged
		// hence, check eth client for this log's tx hash to confirm
		r.lggr.Debugf("log block %d does not exist in block subscriber or block hash does not match for upkeepId %s. this may be caused by block subscriber outdated due to re-org, querying eth client to confirm", logBlockNumber, upkeepId)
	} else {
		r.lggr.Debugf("log block not provided, querying eth client for tx hash %s for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
	}
	// query eth client as a fallback
	bn, bh, err := core.GetTxBlock(r.ctx, r.client, p.Trigger.LogTriggerExtension.TxHash)
	if err != nil {
		// primitive way of checking errors
		if strings.Contains(err.Error(), "missing required field") || strings.Contains(err.Error(), "not found") {
			return encoding.UpkeepFailureReasonTxHashNoLongerExists, encoding.NoPipelineError, false
		}
		r.lggr.Warnf("failed to query tx hash %s for upkeepId %s: %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId, err.Error())
		return encoding.UpkeepFailureReasonNone, encoding.RpcFlakyFailure, true
	}
	if bn == nil {
		r.lggr.Warnf("tx hash %s does not exist on chain for upkeepId %s.", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
		return encoding.UpkeepFailureReasonTxHashNoLongerExists, encoding.NoPipelineError, false
	}
	if bh.Hex() != logBlockHash.Hex() {
		r.lggr.Warnf("tx hash %s reorged from expected blockhash %s to %s for upkeepId %s.", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), logBlockHash.Hex(), bh.Hex(), upkeepId)
		return encoding.UpkeepFailureReasonTxHashReorged, encoding.NoPipelineError, false
	}
	r.lggr.Debugf("tx hash %s exists on chain for upkeepId %s", hexutil.Encode(p.Trigger.LogTriggerExtension.TxHash[:]), upkeepId)
	return encoding.UpkeepFailureReasonNone, encoding.NoPipelineError, false
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
		if state != encoding.NoPipelineError {
			results[i] = encoding.GetIneligibleCheckResultWithoutPerformData(p, encoding.UpkeepFailureReasonNone, state, retryable)
			continue
		}

		opts := r.buildCallOpts(ctx, block)
		var payload []byte
		var err error
		uid := &ocr2keepers.UpkeepIdentifier{}
		uid.FromBigInt(upkeepId)
		switch core.GetUpkeepType(*uid) {
		case types.LogTrigger:
			reason, state, retryable := r.verifyLogExists(upkeepId, p)
			if reason != encoding.UpkeepFailureReasonNone || state != encoding.NoPipelineError {
				results[i] = encoding.GetIneligibleCheckResultWithoutPerformData(p, reason, state, retryable)
				continue
			}

			// check data will include the log trigger config
			payload, err = r.abi.Pack("checkUpkeep", upkeepId, p.CheckData)
			if err != nil {
				// pack error, no retryable
				r.lggr.Warnf("failed to pack log trigger checkUpkeep data for upkeepId %s with check data %s: %s", upkeepId, hexutil.Encode(p.CheckData), err)
				results[i] = encoding.GetIneligibleCheckResultWithoutPerformData(p, encoding.UpkeepFailureReasonNone, encoding.PackUnpackDecodeFailed, false)
				continue
			}
		default:
			// checkUpkeep is overloaded on the contract for conditionals and log upkeeps
			// Need to use the first function (checkUpkeep0) for conditionals
			payload, err = r.abi.Pack("checkUpkeep0", upkeepId)
			if err != nil {
				// pack error, no retryable
				r.lggr.Warnf("failed to pack conditional checkUpkeep data for upkeepId %s with check data %s: %s", upkeepId, hexutil.Encode(p.CheckData), err)
				results[i] = encoding.GetIneligibleCheckResultWithoutPerformData(p, encoding.UpkeepFailureReasonNone, encoding.PackUnpackDecodeFailed, false)
				continue
			}
		}
		indices[len(checkReqs)] = i
		results[i] = encoding.GetIneligibleCheckResultWithoutPerformData(p, encoding.UpkeepFailureReasonNone, encoding.NoPipelineError, false)

		var result string
		checkReqs = append(checkReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"from": zeroAddress,
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
			latestBlockNumber := int64(0)
			latestBlock := r.bs.latestBlock.Load()
			if latestBlock != nil {
				latestBlockNumber = int64(latestBlock.Number)
			}
			checkBlock, _, _ := r.getBlockAndUpkeepId(payloads[index].UpkeepID, payloads[index].Trigger)
			// Exploratory: remove reliance on primitive way of checking errors
			blockNotFound := (strings.Contains(req.Error.Error(), "header not found") || strings.Contains(req.Error.Error(), "missing trie node"))
			if blockNotFound && latestBlockNumber-checkBlock.Int64() > checkBlockTooOldRange {
				// Check block not found in RPC and it is too old, non-retryable error
				r.lggr.Warnf("block not found error encountered in check result for upkeepId %s, check block %d, latest block %d: %s", results[index].UpkeepID.String(), checkBlock.Int64(), latestBlockNumber, req.Error)
				results[index].Retryable = false
				results[index].PipelineExecutionState = uint8(encoding.CheckBlockTooOld)
			} else {
				// individual upkeep failed in a batch call, likely a flay RPC error, consider retryable
				r.lggr.Warnf("rpc error encountered in check result for upkeepId %s: %s", results[index].UpkeepID.String(), req.Error)
				results[index].Retryable = true
				results[index].PipelineExecutionState = uint8(encoding.RpcFlakyFailure)
			}
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

		oc, err := r.fetchUpkeepOffchainConfig(upkeepId)
		if err != nil {
			// this is mostly caused by RPC flakiness
			r.lggr.Errorw("failed get offchain config, gas price check will be disabled", "err", err, "upkeepId", upkeepId, "block", block)
		}
		fr := gasprice.CheckGasPrice(ctx, upkeepId, oc, r.ge, r.lggr)
		if uint8(fr) == uint8(encoding.UpkeepFailureReasonGasPriceTooHigh) {
			r.lggr.Debugf("upkeep %s upkeep failure reason is %d", upkeepId, fr)
			checkResults[i].Eligible = false
			checkResults[i].Retryable = false
			checkResults[i].IneligibilityReason = uint8(fr)
			continue
		}

		// Since checkUpkeep is true, simulate perform upkeep to ensure it doesn't revert
		payload, err := r.abi.Pack("simulatePerformUpkeep", upkeepId, cr.PerformData)
		if err != nil {
			// pack failed, not retryable
			r.lggr.Warnf("failed to pack perform data %s for %s: %s", hexutil.Encode(cr.PerformData), upkeepId, err)
			checkResults[i].Eligible = false
			checkResults[i].PipelineExecutionState = uint8(encoding.PackUnpackDecodeFailed)
			continue
		}

		opts := r.buildCallOpts(ctx, block)
		var result string
		performReqs = append(performReqs, rpc.BatchElem{
			Method: "eth_call",
			Args: []interface{}{
				map[string]interface{}{
					"from": zeroAddress,
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
			checkResults[idx].PipelineExecutionState = uint8(encoding.RpcFlakyFailure)
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
			r.lggr.Warnf("upkeepId %s is not eligible after simulation of perform", checkResults[idx].UpkeepID.String())
			checkResults[performToKeyIdx[i]].Eligible = false
			checkResults[performToKeyIdx[i]].IneligibilityReason = uint8(encoding.UpkeepFailureReasonSimulationFailed)
		}
	}

	return checkResults, nil
}
