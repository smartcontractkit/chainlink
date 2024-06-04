package pipeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type OnChainRead struct {
	BaseTask `mapstructure:",squash"`

	ChainID         string `json:"chainID"`
	Network         string `json:"network"`
	ContractAddress string `json:"contractAddress"`
	ContractName    string `json:"contractName"`
	MethodName      string `json:"methodName"`
	Params          string `json:"params"`
	Config          string `json:"config"`

	relayers map[types.RelayID]loop.Relayer
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
		chainID         StringParam
		network         StringParam
		contractAddress StringParam
		contractName    StringParam
		methodName      StringParam
		params          SliceParam
		//config          StringParam
	)

	err = multierr.Combine(
		errors.Wrap(ResolveParam(&chainID, From(VarExpr(t.ChainID, vars), NonemptyString(t.ChainID))), "chainID"),
		errors.Wrap(ResolveParam(&network, From(VarExpr(t.Network, vars), NonemptyString(t.Network))), "network"),
		errors.Wrap(ResolveParam(&contractAddress, From(VarExpr(t.ContractAddress, vars), NonemptyString(t.ContractAddress))), "contractAddress"),
		errors.Wrap(ResolveParam(&contractName, From(VarExpr(t.ContractName, vars), NonemptyString(t.ContractName))), "contractName"),
		errors.Wrap(ResolveParam(&methodName, From(VarExpr(t.MethodName, vars), NonemptyString(t.MethodName))), "methodName"),
		errors.Wrap(ResolveParam(&params, From(VarExpr(t.Params, vars), VarExpr(t.Params, vars), Inputs(inputs))), "params"),
		//errors.Wrap(ResolveParam(&config, From(VarExpr(t.Config, vars), NonemptyString(t.Config))), "config"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	//Create relayID
	relayID := types.NewRelayID(network.String(), chainID.String())

	r, ok := t.relayers[relayID]

	if !ok {
		return Result{Error: fmt.Errorf("relayer not found for network %q and chainID: %q ", network.String(), chainID.String())}, runInfo
	}

	c := `{
	"contracts": {
		"median": {
			"contractABI": "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"calldata\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTransmissionDetails\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"},{\"internalType\":\"int192\",\"name\":\"latestAnswer\",\"type\":\"int192\"},{\"internalType\":\"uint64\",\"name\":\"latestTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"}],\"name\":\"owedPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
			"configs": {
				"latestTransmissionDetails": "{\n  \"chainSpecificName\": \"latestTransmissionDetails\"\n}\n",
				"hasAccess": "{\"chainSpecificName\":\"hasAccess\"}",
				"owedPayment": "{\n  \"chainSpecificName\": \"owedPayment\"\n}\n"
			}
		}
	}
}` //TODO: @george-dorin move config out

	csr, err := r.NewContractStateReader(ctx, []byte(c))
	if err != nil {
		return Result{Error: err}, runInfo
	}

	err = csr.Bind(ctx, []types.BoundContract{{
		Address: contractAddress.String(),
		Name:    contractName.String(),
		Pending: false,
	}})
	if err != nil {
		return Result{Error: err}, runInfo
	}

	methodParams := map[string]any{}
	err = json.Unmarshal([]byte(t.Params), &methodParams)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	var response any
	err = csr.GetLatestValue(ctx, t.ContractName, t.MethodName, methodParams, &response)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	b, err := json.Marshal(response)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	return Result{Value: string(b)}, runInfo
}
