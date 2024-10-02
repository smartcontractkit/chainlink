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

type ContractRead struct {
	BaseTask `mapstructure:",squash"`

	Address         string `json:"address"`
	Name            string `json:"name"`
	Method          string `json:"method"`
	ConfidenceLevel string `json:"confidenceLevel"`
	Params          string `json:"params"`

	RelayConfig map[string]interface{} `json:"config"`
	Relay       string                 `json:"relay"`

	csrm *contractReaderManager
}

var _ Task = (*ContractRead)(nil)

func (t *ContractRead) Type() TaskType {
	return TaskTypeContractRead
}

func (t *ContractRead) Run(ctx context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}
	var (
		address         StringParam
		name            StringParam
		method          StringParam
		confidenceLevel StringParam
		params          SliceParam
	)

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&address, From(VarExpr(t.Address, vars), NonemptyString(t.Address))), "address"),
		errors.Wrap(ResolveParam(&name, From(VarExpr(t.Name, vars), NonemptyString(t.Name))), "name"),
		errors.Wrap(ResolveParam(&method, From(VarExpr(t.Method, vars), NonemptyString(t.Method))), "method"),
		errors.Wrap(ResolveParam(&confidenceLevel, From(VarExpr(t.ConfidenceLevel, vars), t.ConfidenceLevel)), "confidenceLevel"),
		errors.Wrap(ResolveParam(&params, From(VarExpr(t.Params, vars), VarExpr(t.Params, vars), Inputs(inputs))), "params"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if confidenceLevel == "" {
		confidenceLevel = StringParam(primitives.Finalized)
	} else if confidenceLevel != StringParam(primitives.Finalized) && confidenceLevel != StringParam(primitives.Unconfirmed) {
		return Result{Error: errors.Wrap(err, "invalid confidence level")}, runInfo
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

	csr, rID, err := t.csrm.GetOrCreate(relayID, name.String(), address.String(), method.String(), crcb)

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
	err = csr.GetLatestValue(ctx, rID, primitives.ConfidenceLevel(confidenceLevel), methodParams, &response)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	return Result{Value: response}, runInfo
}
