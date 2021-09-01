package pipeline

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

//
// Return types:
//     []byte
//
type ETHCallTask struct {
	BaseTask            `mapstructure:",squash"`
	Contract            string `json:"contract"`
	Data                string `json:"data"`
	Gas                 string `json:"gas"`
	ExtractRevertReason bool   `json:"extractRevertReason"`
	EVMChainID          string `json:"evmChainID" mapstructure:"evmChainID"`

	chainSet evm.ChainSet
	config   Config
}

var _ Task = (*ETHCallTask)(nil)

func (t *ETHCallTask) Type() TaskType {
	return TaskTypeETHCall
}

func (t *ETHCallTask) Run(ctx context.Context, vars Vars, inputs []Result) (result Result) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}
	}

	var (
		contractAddr AddressParam
		data         BytesParam
		gas          Uint64Param
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&contractAddr, From(VarExpr(t.Contract, vars), NonemptyString(t.Contract))), "contract"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), JSONWithVarExprs(t.Data, vars, false))), "data"),
		errors.Wrap(ResolveParam(&gas, From(VarExpr(t.Gas, vars), NonemptyString(t.Gas), 0)), "gas"),
	)
	if err != nil {
		return Result{Error: err}
	} else if len(data) == 0 {
		return Result{Error: errors.Wrapf(ErrBadInput, "data param must not be empty")}
	}

	call := ethereum.CallMsg{
		To:   (*common.Address)(&contractAddr),
		Data: []byte(data),
		Gas:  uint64(gas),
	}

	chain, err := getChainByString(t.chainSet, t.EVMChainID)
	if err != nil {
		return Result{Error: err}
	}

	resp, err := chain.Client().CallContract(ctx, call, nil)
	if err != nil {
		if t.ExtractRevertReason {
			err = t.retrieveRevertReason(err)
		}

		return Result{Error: err}
	}
	return Result{Value: resp}
}

func (t *ETHCallTask) retrieveRevertReason(baseErr error) error {
	reason, err := eth.ExtractRevertReasonFromRPCError(baseErr)
	if err != nil {
		logger.Default.WithError(err).Errorw("failed to extract revert reason", "baseErr", baseErr)
		return baseErr
	}

	return errors.Wrap(baseErr, reason)
}
