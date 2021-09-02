package pipeline

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/chains/evm"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

//
// Return types:
//     []byte
//
type ETHCallTask struct {
	BaseTask   `mapstructure:",squash"`
	Contract   string `json:"contract"`
	Data       string `json:"data"`
	EVMChainID string `json:"evmChainID" mapstructure:"evmChainID"`

	chainSet evm.ChainSet
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
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&contractAddr, From(NonemptyString(t.Contract))), "contract"),
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars), JSONWithVarExprs(t.Data, vars, false))), "data"),
	)
	if err != nil {
		return Result{Error: err}
	} else if len(data) == 0 {
		return Result{Error: errors.Wrapf(ErrBadInput, "data param must not be empty")}
	}

	call := ethereum.CallMsg{
		To:   (*common.Address)(&contractAddr),
		Data: []byte(data),
	}

	chain, err := getChainByString(t.chainSet, t.EVMChainID)
	if err != nil {
		return Result{Error: err}
	}

	resp, err := chain.Client().CallContract(ctx, call, nil)
	if err != nil {
		return Result{Error: err}
	}
	return Result{Value: resp}
}
