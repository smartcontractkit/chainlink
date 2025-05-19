package read

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
)

var errEmptyOutput = errors.New("rpc call output is empty (make sure that the contract method exists and rpc is healthy)")

const (
	// DefaultRpcBatchSizeLimit defines the maximum number of rpc requests to be included in a batch.
	DefaultRpcBatchSizeLimit = 100

	// DefaultRpcBatchBackOffMultiplier defines the rate of reducing the batch size limit for retried calls.
	// For example if limit is 20 and multiplier is 4:
	// 1.        20
	// 2. 20/4 = 5
	// 3. 5/4  = 1
	DefaultRpcBatchBackOffMultiplier = 5

	// DefaultMaxParallelRpcCalls defines the default maximum number of individual in-parallel rpc calls.
	DefaultMaxParallelRpcCalls = 10
)

// BatchResult is organised by contracts names, key is contract name.
type BatchResult map[string]ContractResults
type ContractResults []MethodCallResult
type MethodCallResult struct {
	Address     string
	MethodName  string
	ReturnValue any
	Err         error
}

type BatchCall []Call
type Call struct {
	ContractAddress        common.Address
	ContractName, ReadName string
	Params, ReturnVal      any
}

func (c BatchCall) String() string {
	callString := ""
	for _, call := range c {
		callString += fmt.Sprintf("%s\n", call.String())
	}
	return callString
}

// Implement the String method for the Call struct
func (c Call) String() string {
	return fmt.Sprintf("contractAddress: %s, contractName: %s, method: %s, params: %+v returnValType: %T",
		c.ContractAddress.Hex(), c.ContractName, c.ReadName, c.Params, c.ReturnVal)
}

type BatchCaller interface {
	// BatchCall executes all the provided BatchRequest and returns the results in the same order
	// of the calls. Pass blockNumber=0 to use the latest block.
	BatchCall(ctx context.Context, blockNumber uint64, batchRequests BatchCall) (BatchResult, error)
}

type EVMBatchCaller interface {
	BatchCallContext(context.Context, []rpc.BatchElem) error
}

// dynamicLimitedBatchCaller makes batched rpc calls and perform retries by reducing the batch size on each retry.
type dynamicLimitedBatchCaller struct {
	bc *defaultEvmBatchCaller
}

func NewDynamicLimitedBatchCaller(lggr logger.Logger, codec types.Codec, evmClient EVMBatchCaller, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint) BatchCaller {
	return &dynamicLimitedBatchCaller{
		bc: newDefaultEvmBatchCaller(lggr, evmClient, codec, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit),
	}
}

func (c *dynamicLimitedBatchCaller) BatchCall(ctx context.Context, blockNumber uint64, reqs BatchCall) (BatchResult, error) {
	return c.bc.batchCallDynamicLimitRetries(ctx, blockNumber, reqs)
}

type defaultEvmBatchCaller struct {
	lggr                  logger.Logger
	evmClient             EVMBatchCaller
	codec                 types.Codec
	batchSizeLimit        uint
	parallelRpcCallsLimit uint
	backOffMultiplier     uint
}

// NewDefaultEvmBatchCaller returns a new batch caller instance.
// batchCallLimit defines the maximum number of calls for BatchCallLimit method, pass 0 to keep the default.
// backOffMultiplier defines the back-off strategy for retries on BatchCallDynamicLimitRetries method, pass 0 to keep the default.
func newDefaultEvmBatchCaller(
	lggr logger.Logger, evmClient EVMBatchCaller, codec types.Codec, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint,
) *defaultEvmBatchCaller {
	batchSize := uint(DefaultRpcBatchSizeLimit)
	if batchSizeLimit > 0 {
		batchSize = batchSizeLimit
	}

	multiplier := uint(DefaultRpcBatchBackOffMultiplier)
	if backOffMultiplier > 0 {
		multiplier = backOffMultiplier
	}

	parallelRpcCalls := uint(DefaultMaxParallelRpcCalls)
	if parallelRpcCallsLimit > 0 {
		parallelRpcCalls = parallelRpcCallsLimit
	}

	return &defaultEvmBatchCaller{
		lggr:                  lggr,
		evmClient:             evmClient,
		codec:                 codec,
		batchSizeLimit:        batchSize,
		parallelRpcCallsLimit: parallelRpcCalls,
		backOffMultiplier:     multiplier,
	}
}

// batchCall formats a batch, calls the rpc client, and unpacks results.
// this function only returns errors of type Error which should wrap lower errors.
func (c *defaultEvmBatchCaller) batchCall(ctx context.Context, blockNumber uint64, batchCall BatchCall) ([]dataAndErr, error) {
	if len(batchCall) == 0 {
		return nil, nil
	}

	blockNumStr := "latest"
	if blockNumber > 0 {
		blockNumStr = hexutil.EncodeBig(big.NewInt(0).SetUint64(blockNumber))
	}

	rpcBatchCalls, hexEncodedOutputs, err := c.createBatchCalls(ctx, batchCall, blockNumStr)
	if err != nil {
		return nil, err
	}

	if err = c.evmClient.BatchCallContext(ctx, rpcBatchCalls); err != nil {
		// return a basic read error with no detail or result since this is a general client
		// error instead of an error for a specific batch call.
		return nil, Error{
			Err:  fmt.Errorf("%w: batch call context: %s", types.ErrInternal, err.Error()),
			Type: batchReadType,
		}
	}

	results, err := c.unpackBatchResults(ctx, batchCall, rpcBatchCalls, hexEncodedOutputs, blockNumStr)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *defaultEvmBatchCaller) createBatchCalls(
	ctx context.Context,
	batchCall BatchCall,
	block string,
) ([]rpc.BatchElem, []string, error) {
	rpcBatchCalls := make([]rpc.BatchElem, len(batchCall))
	hexEncodedOutputs := make([]string, len(batchCall))

	for idx, call := range batchCall {
		data, err := c.codec.Encode(ctx, call.Params, codec.WrapItemType(call.ContractName, call.ReadName, true))
		if err != nil {
			return nil, nil, newErrorFromCall(
				fmt.Errorf("%w: encode params: %s", types.ErrInvalidConfig, err.Error()),
				call,
				block,
				batchReadType,
			)
		}

		rpcBatchCalls[idx] = rpc.BatchElem{
			Method: "eth_call",
			Args: []any{
				map[string]interface{}{
					"from": common.Address{},
					"to":   call.ContractAddress,
					"data": hexutil.Bytes(data),
				},
				block,
			},
			Result: &hexEncodedOutputs[idx],
		}
	}

	return rpcBatchCalls, hexEncodedOutputs, nil
}

func (c *defaultEvmBatchCaller) unpackBatchResults(
	ctx context.Context,
	batchCall BatchCall,
	rpcBatchCalls []rpc.BatchElem,
	hexEncodedOutputs []string,
	block string,
) ([]dataAndErr, error) {
	results := make([]dataAndErr, len(batchCall))

	for idx, call := range batchCall {
		results[idx] = dataAndErr{
			address:      call.ContractAddress.Hex(),
			contractName: call.ContractName,
			methodName:   call.ReadName,
			returnVal:    call.ReturnVal,
		}

		if rpcBatchCalls[idx].Error != nil {
			results[idx].err = newErrorFromCall(
				fmt.Errorf("%w: rpc call error: %w", types.ErrInternal, rpcBatchCalls[idx].Error),
				call, block, batchReadType,
			)

			continue
		}

		if hexEncodedOutputs[idx] == "" {
			// Some RPCs instead of returning "0x" are returning an empty string.
			// We are overriding this behaviour for consistent handling of this scenario.
			hexEncodedOutputs[idx] = "0x"
		}

		packedBytes, err := hexutil.Decode(hexEncodedOutputs[idx])
		if err != nil {
			callErr := newErrorFromCall(
				fmt.Errorf("%w: hex decode result: %s", types.ErrInternal, err.Error()),
				call, block, batchReadType,
			)

			callErr.Result = &hexEncodedOutputs[idx]

			return nil, callErr
		}

		// the codec can't do anything with no bytes, so skip decoding and allow
		// the result to be the empty struct or value
		if len(packedBytes) > 0 {
			if err = c.codec.Decode(
				ctx,
				packedBytes,
				call.ReturnVal,
				codec.WrapItemType(call.ContractName, call.ReadName, false),
			); err != nil {
				if len(packedBytes) == 0 {
					callErr := newErrorFromCall(
						fmt.Errorf("%w: %w: %s", types.ErrInternal, errEmptyOutput, err.Error()),
						call, block, batchReadType,
					)

					callErr.Result = &hexEncodedOutputs[idx]

					results[idx].err = callErr
				} else {
					callErr := newErrorFromCall(
						fmt.Errorf("%w: codec decode result: %s", types.ErrInvalidType, err.Error()),
						call, block, batchReadType,
					)

					callErr.Result = &hexEncodedOutputs[idx]
					results[idx].err = callErr
				}

				continue
			}
		}

		results[idx].returnVal = call.ReturnVal
	}

	return results, nil
}

func (c *defaultEvmBatchCaller) batchCallDynamicLimitRetries(ctx context.Context, blockNumber uint64, calls BatchCall) (BatchResult, error) {
	lim := c.batchSizeLimit

	// Limit the batch size to the number of calls
	if uint(len(calls)) < lim {
		lim = uint(len(calls))
	}

	for {
		results, err := c.batchCallLimit(ctx, blockNumber, calls, lim)
		if err == nil {
			return results, nil
		}

		if lim <= 1 {
			return nil, Error{
				Err:  fmt.Errorf("%w: limited call: call data: %+v", err, calls),
				Type: batchReadType,
			}
		}

		newLim := lim / c.backOffMultiplier
		if newLim == 0 || newLim == lim {
			newLim = 1
		}

		lim = newLim

		c.lggr.Debugf("retrying batch call with %d calls and %d limit that failed with error=%s", len(calls), lim, err)
	}
}

type dataAndErr struct {
	address                  string
	contractName, methodName string
	returnVal                any
	err                      error
}

func (c *defaultEvmBatchCaller) batchCallLimit(ctx context.Context, blockNumber uint64, calls BatchCall, batchSizeLimit uint) (BatchResult, error) {
	if batchSizeLimit <= 0 {
		res, err := c.batchCall(ctx, blockNumber, calls)

		return convertToBatchResult(res), err
	}

	type job struct {
		blockNumber uint64
		calls       BatchCall
		results     []dataAndErr
	}

	jobs := make([]job, 0)
	for i := 0; i < len(calls); i += int(batchSizeLimit) {
		idxFrom := i

		idxTo := idxFrom + int(batchSizeLimit)
		if idxTo > len(calls) {
			idxTo = len(calls)
		}

		jobs = append(jobs, job{blockNumber: blockNumber, calls: calls[idxFrom:idxTo], results: nil})
	}

	if c.parallelRpcCallsLimit > 1 {
		eg := new(errgroup.Group)
		eg.SetLimit(int(c.parallelRpcCallsLimit))

		for jobIdx := range jobs {
			jobIdx := jobIdx

			eg.Go(func() error {
				res, err := c.batchCall(ctx, jobs[jobIdx].blockNumber, jobs[jobIdx].calls)
				if err != nil {
					return err
				}

				jobs[jobIdx].results = res

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, err
		}
	} else {
		var err error

		for jobIdx := range jobs {
			jobs[jobIdx].results, err = c.batchCall(ctx, jobs[jobIdx].blockNumber, jobs[jobIdx].calls)
			if err != nil {
				return nil, err
			}
		}
	}

	var results []dataAndErr
	for _, jb := range jobs {
		results = append(results, jb.results...)
	}

	return convertToBatchResult(results), nil
}

func convertToBatchResult(data []dataAndErr) BatchResult {
	if data == nil {
		return nil
	}

	batchResult := make(BatchResult)
	for _, d := range data {
		methodCall := MethodCallResult{
			Address:     d.address,
			MethodName:  d.methodName,
			ReturnValue: d.returnVal,
			Err:         d.err,
		}

		if _, exists := batchResult[d.contractName]; !exists {
			batchResult[d.contractName] = ContractResults{}
		}

		batchResult[d.contractName] = append(batchResult[d.contractName], methodCall)
	}

	return batchResult
}
