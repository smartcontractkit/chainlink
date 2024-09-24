package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type OnChainRead struct {
	BaseTask `mapstructure:",squash"`

	ContractAddress string                 `json:"contractAddress"`
	ContractName    string                 `json:"contractName"`
	MethodName      string                 `json:"methodName"`
	Params          string                 `json:"params"`
	RelayConfig     map[string]interface{} `json:"config"`
	Relay           string                 `json:"relay"`

	csrm *contractReaderManager
}

var _ Task = (*OnChainRead)(nil)

func (t *OnChainRead) Type() TaskType {
	return TaskTypeOnchainRead
}

func (t *OnChainRead) Run(ctx context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}
	var (
		contractAddress StringParam
		contractName    StringParam
		methodName      StringParam
		params          SliceParam
	)

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&contractAddress, From(VarExpr(t.ContractAddress, vars), NonemptyString(t.ContractAddress))), "contractAddress"),
		errors.Wrap(ResolveParam(&contractName, From(VarExpr(t.ContractName, vars), NonemptyString(t.ContractName))), "contractName"),
		errors.Wrap(ResolveParam(&methodName, From(VarExpr(t.MethodName, vars), NonemptyString(t.MethodName))), "methodName"),
		errors.Wrap(ResolveParam(&params, From(VarExpr(t.Params, vars), VarExpr(t.Params, vars), Inputs(inputs))), "params"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	//Fetch network and chainID to create the RelayID
	c, ok := t.RelayConfig["chainID"]
	if !ok {
		return Result{Error: fmt.Errorf("cannot get chainID")}, runInfo
	}
	chainID, ok := c.(string)
	if !ok {
		return Result{Error: fmt.Errorf("cannot get chainID,expected string but got %T", c)}, runInfo
	}

	relayID := types.NewRelayID(t.Relay, chainID)
	crc, ok := t.RelayConfig["chainReader"]
	if !ok {
		return Result{Error: fmt.Errorf("cannot find chainReader config")}, runInfo
	}

	crcb, err := json.Marshal(crc)
	if err != nil {
		return Result{Error: fmt.Errorf("cannot marshal chainReader config")}, runInfo
	}

	csr, rID, err := t.csrm.GetOrCreate(relayID, contractName.String(), contractAddress.String(), methodName.String(), crcb)

	if err != nil {
		return Result{Error: err}, runInfo
	}

	methodParams := map[string]any{}
	if json.Valid([]byte(t.Params)) {
		err = json.Unmarshal([]byte(t.Params), &methodParams)
		if err != nil {
			return Result{Error: err}, runInfo
		}
	}

	var response any
	err = csr.GetLatestValue(ctx, rID, primitives.Finalized, methodParams, &response)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	return Result{Value: response}, runInfo
}
