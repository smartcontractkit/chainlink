package chainreadertest

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

var (
	// ChainReader is a static implementation of [types.ContractReader], [testtypes.Evaluator] and [types.PluginProvider
	// it is used for testing the [types.PluginProvider] interface
	ChainReader = staticChainReader{
		contractName:   "anyContract",
		contractMethod: "anyMethod",
		latestValue:    map[string]any{"ret1": "latestValue1", "ret2": "latestValue2"},
		params:         map[string]any{"param1": "value1", "param2": "value2"},
	}
)

// staticChainReader is a static implementation of ChainReaderTester
type staticChainReader struct {
	contractName   string
	contractMethod string
	latestValue    map[string]any
	params         map[string]any
}

var _ testtypes.Evaluator[types.ContractReader] = staticChainReader{}
var _ types.ContractReader = staticChainReader{}

func (c staticChainReader) Start(_ context.Context) error { return nil }

func (c staticChainReader) Close() error { return nil }

func (c staticChainReader) Ready() error { panic("unimplemented") }

func (c staticChainReader) Name() string { panic("unimplemented") }

func (c staticChainReader) HealthReport() map[string]error { panic("unimplemented") }

func (c staticChainReader) Bind(_ context.Context, _ []types.BoundContract) error {
	return nil
}

func (c staticChainReader) GetLatestValue(_ context.Context, contractName, method string, params, returnVal any) error {
	if !assert.ObjectsAreEqual(contractName, c.contractName) {
		return fmt.Errorf("%w: expected report context %v but got %v", types.ErrInvalidType, c.contractName, contractName)
	}
	if method != c.contractMethod {
		return fmt.Errorf("%w: expected generic contract method %v but got %v", types.ErrInvalidType, c.contractMethod, method)
	}
	//gotParams, ok := params.(*map[string]string)
	gotParams, ok := params.(*map[string]any)
	if !ok {
		return fmt.Errorf("%w: Invalid parameter type received in GetLatestValue. Expected %T but received %T", types.ErrInvalidEncoding, c.params, params)
	}
	if (*gotParams)["param1"] != c.params["param1"] || (*gotParams)["param2"] != c.params["param2"] {
		return fmt.Errorf("%w: Wrong params value received in GetLatestValue. Expected %v but received %v", types.ErrInvalidEncoding, c.params, *gotParams)
	}

	ret, ok := returnVal.(*map[string]any)
	if !ok {
		return fmt.Errorf("%w: Wrong type passed for retVal param to GetLatestValue impl (expected %T instead of %T", types.ErrInvalidType, c.latestValue, returnVal)
	}

	(*ret)["ret1"] = c.latestValue["ret1"]
	(*ret)["ret2"] = c.latestValue["ret2"]

	return nil
}

func (c staticChainReader) QueryKey(_ context.Context, _ string, _ query.KeyFilter, _ query.LimitAndSort, _ any) ([]types.Sequence, error) {
	return nil, nil
}

func (c staticChainReader) Evaluate(ctx context.Context, cr types.ContractReader) error {
	gotLatestValue := make(map[string]any)
	err := cr.GetLatestValue(ctx, c.contractName, c.contractMethod, &c.params, &gotLatestValue)
	if err != nil {
		return fmt.Errorf("failed to call GetLatestValue(): %w", err)
	}

	if !assert.ObjectsAreEqual(gotLatestValue, c.latestValue) {
		return fmt.Errorf("GetLatestValue: expected %v but got %v", c.latestValue, gotLatestValue)
	}
	return nil
}
