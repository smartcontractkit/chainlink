package pipeline

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
)

// Return types:
//
//	uint64
type EstimateGasLimitTask struct {
	BaseTask   `mapstructure:",squash"`
	Input      string `json:"input"`
	From       string `json:"from"`
	To         string `json:"to"`
	Multiplier string `json:"multiplier"`
	Data       string `json:"data"`
	EVMChainID string `json:"evmChainID" mapstructure:"evmChainID"`
	Block      string `json:"block"`

	specGasLimit *uint32
	legacyChains legacyevm.LegacyChainContainer
	jobType      string
}

type GasEstimator interface {
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
}

var (
	_                    Task = (*EstimateGasLimitTask)(nil)
	ErrInvalidMultiplier      = errors.New("Invalid multiplier")
)

func (t *EstimateGasLimitTask) Type() TaskType {
	return TaskTypeEstimateGasLimit
}

func (t *EstimateGasLimitTask) getEvmChainID() string {
	if t.EVMChainID == "" {
		t.EVMChainID = "$(jobSpec.evmChainID)"
	}
	return t.EVMChainID
}

func (t *EstimateGasLimitTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		fromAddr   AddressParam
		toAddr     AddressParam
		data       BytesParam
		multiplier DecimalParam
		chainID    StringParam
		block      StringParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddr, From(VarExpr(t.From, vars), utils.ZeroAddress)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		// Default to 1, i.e. exactly what estimateGas suggests
		errors.Wrap(ResolveParam(&multiplier, From(VarExpr(t.Multiplier, vars), NonemptyString(t.Multiplier), decimal.New(1, 0))), "multiplier"),
		errors.Wrap(ResolveParam(&chainID, From(VarExpr(t.getEvmChainID(), vars), NonemptyString(t.getEvmChainID()), "")), "evmChainID"),
		errors.Wrap(ResolveParam(&block, From(VarExpr(t.Block, vars), t.Block)), "block"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	chain, err := t.legacyChains.Get(string(chainID))
	if err != nil {
		err = fmt.Errorf("%w: %s: %w", ErrInvalidEVMChainID, chainID, err)
		return Result{Error: err}, runInfo
	}

	maximumGasLimit := SelectGasLimit(chain.Config().EVM().GasEstimator(), t.jobType, t.specGasLimit)
	to := common.Address(toAddr)
	var gasLimit hexutil.Uint64
	args := map[string]interface{}{
		"from":  common.Address(fromAddr),
		"to":    &to,
		"input": hexutil.Bytes([]byte(data)),
	}

	selectedBlock, err := selectBlock(string(block))
	if err != nil {
		return Result{Error: err}, runInfo
	}
	err = chain.Client().CallContext(ctx,
		&gasLimit,
		"eth_estimateGas",
		args,
		selectedBlock,
	)

	if err != nil {
		// Fallback to the maximum conceivable gas limit
		// if we're unable to call estimate gas for whatever reason.
		lggr.Warnw("EstimateGas: unable to estimate, fallback to configured limit", "err", err, "fallback", maximumGasLimit)
		return Result{Value: maximumGasLimit}, runInfo
	}

	gasLimitDecimal, err := decimal.NewFromString(strconv.FormatUint(uint64(gasLimit), 10))
	if err != nil {
		return Result{Error: err}, retryableRunInfo()
	}
	newExp := int64(gasLimitDecimal.Exponent()) + int64(multiplier.Decimal().Exponent())
	if newExp > math.MaxInt32 || newExp < math.MinInt32 {
		return Result{Error: ErrMultiplyOverlow}, retryableRunInfo()
	}
	gasLimitWithMultiplier := gasLimitDecimal.Mul(multiplier.Decimal()).Truncate(0).BigInt()
	if !gasLimitWithMultiplier.IsUint64() {
		return Result{Error: ErrInvalidMultiplier}, retryableRunInfo()
	}
	gasLimitFinal := gasLimitWithMultiplier.Uint64()
	if gasLimitFinal > maximumGasLimit {
		lggr.Warnw("EstimateGas: estimated amount is greater than configured limit, fallback to configured limit",
			"estimate", gasLimitFinal,
			"fallback", maximumGasLimit,
		)
		gasLimitFinal = maximumGasLimit
	}
	return Result{Value: gasLimitFinal}, runInfo
}
