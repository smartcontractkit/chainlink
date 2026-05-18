package read

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
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
	ContractAddress          common.Address
	ContractName, MethodName string
	Params, ReturnVal        any
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
		c.ContractAddress.Hex(), c.ContractName, c.MethodName, c.Params, c.ReturnVal)
}

type BatchCaller interface {
	// BatchCall executes all the provided BatchRequest and returns the results in the same order
	// of the calls. Pass blockNumber=0 to use the latest block.
	BatchCall(ctx context.Context, blockNumber uint64, batchRequests BatchCall) (BatchResult, error)
}

// dynamicLimitedBatchCaller makes batched rpc calls and perform retries by reducing the batch size on each retry.
type dynamicLimitedBatchCaller struct {
	bc *defaultEvmBatchCaller
}

func NewDynamicLimitedBatchCaller(lggr logger.Logger, codec types.Codec, evmClient client.Client, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint) BatchCaller {
	return &dynamicLimitedBatchCaller{
		bc: newDefaultEvmBatchCaller(lggr, evmClient, codec, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit),
	}
}

func (c *dynamicLimitedBatchCaller) BatchCall(ctx context.Context, blockNumber uint64, reqs BatchCall) (BatchResult, error) {
	return c.bc.batchCallDynamicLimitRetries(ctx, blockNumber, reqs)
}

type defaultEvmBatchCaller struct {
	lggr                  logger.Logger
	evmClient             client.Client
	codec                 types.Codec
	batchSizeLimit        uint
	parallelRpcCallsLimit uint
	backOffMultiplier     uint
}

// NewDefaultEvmBatchCaller returns a new batch caller instance.
// batchCallLimit defines the maximum number of calls for BatchCallLimit method, pass 0 to keep the default.
// backOffMultiplier defines the back-off strategy for retries on BatchCallDynamicLimitRetries method, pass 0 to keep the default.
func newDefaultEvmBatchCaller(
	lggr logger.Logger, evmClient client.Client, codec types.Codec, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint,
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

func (c *defaultEvmBatchCaller) batchCall(ctx context.Context, blockNumber uint64, batchCall BatchCall) ([]dataAndErr, error) {
	if len(batchCall) == 0 {
		return nil, nil
	}

	packedOutputs := make([]string, len(batchCall))
	rpcBatchCalls := make([]rpc.BatchElem, len(batchCall))
	for i, call := range batchCall {
		data, err := c.codec.Encode(ctx, call.Params, codec.WrapItemType(call.ContractName, call.MethodName, true))
		if err != nil {
			return nil, err
		}

		blockNumStr := "latest"
		if blockNumber > 0 {
			blockNumStr = hexutil.EncodeBig(big.NewInt(0).SetUint64(blockNumber))
		}

		rpcBatchCalls[i] = rpc.BatchElem{
			Method: "eth_call",
			Args: []any{
				map[string]interface{}{
					"from": common.Address{},
					"to":   call.ContractAddress,
					"data": hexutil.Bytes(data),
				},
				blockNumStr,
			},
			Result: &packedOutputs[i],
		}
	}

	if err := c.evmClient.BatchCallContext(ctx, rpcBatchCalls); err != nil {
		return nil, fmt.Errorf("batch call context: %w", err)
	}

	results := make([]dataAndErr, len(batchCall))
	for i, call := range batchCall {
		results[i] = dataAndErr{
			address:      call.ContractAddress.Hex(),
			contractName: call.ContractName,
			methodName:   call.MethodName,
			returnVal:    call.ReturnVal,
		}

		if rpcBatchCalls[i].Error != nil {
			results[i].err = rpcBatchCalls[i].Error
			continue
		}

		if packedOutputs[i] == "" {
			// Some RPCs instead of returning "0x" are returning an empty string.
			// We are overriding this behaviour for consistent handling of this scenario.
			packedOutputs[i] = "0x"
		}

		b, err := hexutil.Decode(packedOutputs[i])
		if err != nil {
			return nil, fmt.Errorf("decode result %s: packedOutputs %s: %w", call, packedOutputs[i], err)
		}

		if err = c.codec.Decode(ctx, b, call.ReturnVal, codec.WrapItemType(call.ContractName, call.MethodName, false)); err != nil {
			if len(b) == 0 {
				results[i].err = fmt.Errorf("unpack result %s: %s: %w", call, err.Error(), errEmptyOutput)
			} else {
				results[i].err = fmt.Errorf("unpack result %s: %w", call, err)
			}
			continue
		}
		results[i].returnVal = call.ReturnVal
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
			return nil, errors.Wrapf(err, "calls %+v", calls)
		}

		newLim := lim / c.backOffMultiplier
		if newLim == 0 || newLim == lim {
			newLim = 1
		}
		lim = newLim
		c.lggr.Errorf("retrying batch call with %d calls and %d limit that failed with error=%s",
			len(calls), lim, err)
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
