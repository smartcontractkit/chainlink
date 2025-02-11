package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
)

// Return types:
//
//	[]byte
type ETHCallTask struct {
	BaseTask            `mapstructure:",squash"`
	Contract            string `json:"contract"`
	From                string `json:"from"`
	Data                string `json:"data"`
	Gas                 string `json:"gas"`
	GasPrice            string `json:"gasPrice"`
	GasTipCap           string `json:"gasTipCap"`
	GasFeeCap           string `json:"gasFeeCap"`
	GasUnlimited        string `json:"gasUnlimited"`
	ExtractRevertReason bool   `json:"extractRevertReason"`
	EVMChainID          string `json:"evmChainID" mapstructure:"evmChainID"`
	Block               string `json:"block"`

	specGasLimit *uint32
	legacyChains legacyevm.LegacyChainContainer
	config       Config
	jobType      string
}

var _ Task = (*ETHCallTask)(nil)

var (
	promETHCallTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_eth_call_execution_time",
		Help: "Time taken to fully execute the ETH call",
	},
		[]string{"pipeline_task_spec_id"},
	)
)

func (t *ETHCallTask) Type() TaskType {
	return TaskTypeETHCall
}

func (t *ETHCallTask) getEvmChainID() string {
	if t.EVMChainID == "" {
		t.EVMChainID = "$(jobSpec.evmChainID)"
	}
	return t.EVMChainID
}

func (t *ETHCallTask) Run(ctx context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		contractAddr AddressParam
		from         AddressParam
		data         BytesParam
		gas          Uint64Param
		gasPrice     MaybeBigIntParam
		gasTipCap    MaybeBigIntParam
		gasFeeCap    MaybeBigIntParam
		gasUnlimited BoolParam
		chainID      StringParam
		block        StringParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&contractAddr, From(VarExpr(t.Contract, vars), NonemptyString(t.Contract))), "contract"),
		errors.Wrap(ResolveParam(&from, From(VarExpr(t.From, vars), NonemptyString(t.From), utils.ZeroAddress)), "from"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), JSONWithVarExprs(t.Data, vars, false))), "data"),
		errors.Wrap(ResolveParam(&gas, From(VarExpr(t.Gas, vars), NonemptyString(t.Gas), 0)), "gas"),
		errors.Wrap(ResolveParam(&gasPrice, From(VarExpr(t.GasPrice, vars), t.GasPrice)), "gasPrice"),
		errors.Wrap(ResolveParam(&gasTipCap, From(VarExpr(t.GasTipCap, vars), t.GasTipCap)), "gasTipCap"),
		errors.Wrap(ResolveParam(&gasFeeCap, From(VarExpr(t.GasFeeCap, vars), t.GasFeeCap)), "gasFeeCap"),
		errors.Wrap(ResolveParam(&chainID, From(VarExpr(t.getEvmChainID(), vars), NonemptyString(t.getEvmChainID()), "")), "evmChainID"),
		errors.Wrap(ResolveParam(&gasUnlimited, From(VarExpr(t.GasUnlimited, vars), NonemptyString(t.GasUnlimited), false)), "gasUnlimited"),
		errors.Wrap(ResolveParam(&block, From(VarExpr(t.Block, vars), t.Block)), "block"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	} else if len(data) == 0 {
		return Result{Error: errors.Wrapf(ErrBadInput, "data param must not be empty")}, runInfo
	}

	chain, err := t.legacyChains.Get(string(chainID))
	if err != nil {
		err = fmt.Errorf("%w: %s: %w", ErrInvalidEVMChainID, chainID, err)
		return Result{Error: err}, runInfo
	}

	var selectedGas uint64
	if gasUnlimited {
		if gas > 0 {
			return Result{Error: errors.Wrapf(ErrBadInput, "gas must be zero when gasUnlimited is true")}, runInfo
		}
	} else {
		if gas > 0 {
			selectedGas = uint64(gas)
		} else {
			selectedGas = SelectGasLimit(chain.Config().EVM().GasEstimator(), t.jobType, t.specGasLimit)
		}
	}

	call := ethereum.CallMsg{
		To:        (*common.Address)(&contractAddr),
		From:      (common.Address)(from),
		Data:      []byte(data),
		Gas:       selectedGas,
		GasPrice:  gasPrice.BigInt(),
		GasTipCap: gasTipCap.BigInt(),
		GasFeeCap: gasFeeCap.BigInt(),
	}

	lggr = logger.Sugared(lggr).With("gas", call.Gas).
		With("gasPrice", call.GasPrice).
		With("gasTipCap", call.GasTipCap).
		With("gasFeeCap", call.GasFeeCap)

	start := time.Now()

	var resp []byte
	blockStr := block.String()
	if blockStr == "" || strings.ToLower(blockStr) == "latest" {
		resp, err = chain.Client().CallContract(ctx, call, nil)
	} else if strings.ToLower(blockStr) == "pending" {
		resp, err = chain.Client().PendingCallContract(ctx, call)
	}

	elapsed := time.Since(start)
	if err != nil {
		if t.ExtractRevertReason {
			rpcError, errExtract := evmclient.ExtractRPCError(err)
			if errExtract == nil {
				// Update error to unmarshalled RPCError with revert data.
				err = rpcError
			} else {
				lggr.Warnw("failed to extract rpc error", "err", err, "errExtract", errExtract)
				// Leave error as is.
			}
		}

		return Result{Error: err}, retryableRunInfo()
	}

	promETHCallTime.WithLabelValues(t.DotID()).Set(float64(elapsed))
	return Result{Value: resp}, runInfo
}
