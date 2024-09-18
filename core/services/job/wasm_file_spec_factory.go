package job

import (
	"context"
	"errors"
	"os"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type WasmFileSpecFactory struct{}

func (w WasmFileSpecFactory) Spec(ctx context.Context, lggr logger.Logger, rawSpec, config []byte) (sdk.WorkflowSpec, error) {
	moduleConfig := &host.ModuleConfig{Logger: lggr}
	spec, err := host.GetWorkflowSpec(moduleConfig, rawSpec, config)
	if err != nil {
		return sdk.WorkflowSpec{}, err
	} else if spec == nil {
		return sdk.WorkflowSpec{}, errors.New("workflow spec not found when running wasm")
	}

	return *spec, nil
}

func (w WasmFileSpecFactory) RawSpec(_ context.Context, wf string) ([]byte, error) {
	return os.ReadFile(wf)
}

var _ SDKWorkflowSpecFactory = (*WasmFileSpecFactory)(nil)
