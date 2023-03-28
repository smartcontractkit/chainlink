package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NOTE: Currently only returns latest block, could be extended in future to
// return block by number or hash

// Return types:
//
// map[string]interface{}
//
// Fields:
// - number: int64
// - hash: common.Hash
// - parentHash: common.Hash
// - timestamp: time.Time
// - baseFeePerGas: *big.Int
// - receiptsRoot: common.Hash
// - transactionsRoot: common.Hash
// - stateRoot: common.Hash
type ETHGetBlockTask struct {
	BaseTask   `mapstructure:",squash"`
	EVMChainID string `json:"evmChainID" mapstructure:"evmChainID"`

	chainSet evm.ChainSet
	config   Config
}

var _ Task = (*ETHGetBlockTask)(nil)

func (t *ETHGetBlockTask) Type() TaskType {
	return TaskTypeETHGetBlock
}

func (t *ETHGetBlockTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var chainID StringParam
	err = errors.Wrap(ResolveParam(&chainID, From(VarExpr(t.EVMChainID, vars), NonemptyString(t.EVMChainID), "")), "evmChainID")
	if err != nil {
		return Result{Error: err}, runInfo
	}

	chain, err := getChainByString(t.chainSet, string(chainID))
	if err != nil {
		return Result{Error: err}, runInfo
	}
	// Use the headtracker's view of the latest block, this is very fast since
	// it doesn't make any external network requests, and it is the
	// headtracker's job to ensure it has an up-to-date view of the chain based
	// on responses from all available RPC nodes
	latestHead := chain.HeadTracker().LatestChain()
	if latestHead == nil {
		logger.Sugared(lggr).AssumptionViolation("HeadTracker unexpectedly returned nil head, falling back to RPC call")
		latestHead, err = chain.Client().HeadByNumber(ctx, nil)
		if err != nil {
			return Result{Error: err}, runInfo
		}
	}

	h := make(map[string]interface{})
	h["number"] = latestHead.Number
	h["hash"] = latestHead.Hash
	h["parentHash"] = latestHead.ParentHash
	h["timestamp"] = latestHead.Timestamp
	h["baseFeePerGas"] = latestHead.BaseFeePerGas
	h["receiptsRoot"] = latestHead.ReceiptsRoot
	h["transactionsRoot"] = latestHead.TransactionsRoot
	h["stateRoot"] = latestHead.StateRoot

	return Result{Value: h}, runInfo
}
