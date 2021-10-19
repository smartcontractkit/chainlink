package pipeline

import (
	"context"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

//
// Return types:
//   uint64
//
type EstimateGasLimitTask struct {
	BaseTask   `mapstructure:",squash"`
	Input      string `json:"input"`
	To         string `json:"to"`
	Multiplier string `json:"multiplier"`
	Data       string `json:"data"`

	EvmGasLimit  uint64
	GasEstimator GasEstimator
}

type GasEstimator interface {
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
}

var (
	_ Task = (*EstimateGasLimitTask)(nil)
)

func (t *EstimateGasLimitTask) Type() TaskType {
	return TaskTypeEstimateGasLimit
}

func (t *EstimateGasLimitTask) Run(_ context.Context, vars Vars, inputs []Result) (result Result) {
	var (
		toAddr     AddressParam
		data       BytesParam
		multiplier DecimalParam
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&toAddr, From(VarExpr(t.To, vars), NonemptyString(t.To))), "to"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), NonemptyString(t.Data))), "data"),
		// Default to 1, i.e. exactly what estimateGas suggests
		errors.Wrap(ResolveParam(&multiplier, From(VarExpr(t.Multiplier, vars), NonemptyString(t.Multiplier), decimal.New(1, 0))), "multiplier"),
	)
	if err != nil {
		return Result{Error: err}
	}

	to := common.Address(toAddr)
	gasLimit, err := t.GasEstimator.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &to,
		Data: data,
	})
	if err != nil {
		// Fallback to the maximum conceivable gas limit
		// if we're unable to call estimate gas for whatever reason.
		logger.Warnw("EstimateGas: unable to estimate, fallback to configured limit", "err", err, "fallback", t.EvmGasLimit)
		return Result{Value: t.EvmGasLimit}
	}
	gasLimitDecimal, err := decimal.NewFromString(strconv.FormatUint(gasLimit, 10))
	if err != nil {
		return Result{Error: err}
	}
	gasLimitWithMultiplier := gasLimitDecimal.Mul(multiplier.Decimal()).Truncate(0).BigInt()
	if !gasLimitWithMultiplier.IsUint64() {
		return Result{Error: errors.New("Invalid multiplier")}
	}
	gasLimitFinal := gasLimitWithMultiplier.Uint64()
	if gasLimitFinal > t.EvmGasLimit {
		gasLimitFinal = t.EvmGasLimit
	}
	return Result{Value: gasLimitFinal}
}
