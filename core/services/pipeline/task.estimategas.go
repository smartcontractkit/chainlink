package pipeline

import (
	"context"
	"math"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//
// Return types:
//   uint64
//
type EstimateGasLimitTask struct {
	BaseTask   `mapstructure:",squash"`
	Input      string `json:"input"`
	From       string `json:"from"`
	To         string `json:"to"`
	Multiplier string `json:"multiplier"`
	Data       string `json:"data"`
	EVMChainID string `json:"evmChainID" mapstructure:"evmChainID"`

	specGasLimit *uint32
	chainSet     evm.ChainSet
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

func (t *EstimateGasLimitTask) Run(_ context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		fromAddr   AddressParam
		toAddr     AddressParam
		data       BytesParam
		multiplier DecimalParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&fromAddr, From(VarExpr(t.From, vars), utils.ZeroAddress)), "from"),
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		// Default to 1, i.e. exactly what estimateGas suggests
		errors.Wrap(ResolveParam(&multiplier, From(VarExpr(t.Multiplier, vars), NonemptyString(t.Multiplier), decimal.New(1, 0))), "multiplier"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	chain, err := getChainByString(t.chainSet, t.EVMChainID)
	if err != nil {
		return Result{Error: err}, retryableRunInfo()
	}
	maximumGasLimit := chain.Config().EvmGasLimitDefault()
	if t.specGasLimit != nil {
		maximumGasLimit = uint64(*t.specGasLimit)
	}
	to := common.Address(toAddr)
	gasLimit, err := chain.Client().EstimateGas(context.Background(), ethereum.CallMsg{
		From: common.Address(fromAddr),
		To:   &to,
		Data: data,
	})
	if err != nil {
		// Fallback to the maximum conceivable gas limit
		// if we're unable to call estimate gas for whatever reason.
		lggr.Warnw("EstimateGas: unable to estimate, fallback to configured limit", "err", err, "fallback", maximumGasLimit)
		return Result{Value: maximumGasLimit}, runInfo
	}
	gasLimitDecimal, err := decimal.NewFromString(strconv.FormatUint(gasLimit, 10))
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
