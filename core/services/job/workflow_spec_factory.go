package job

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type WorkflowSpecFactory interface {
	Spec(ctx context.Context, workflow string, config []byte) (sdk.WorkflowSpec, []byte, string, error)
	RawSpec(ctx context.Context, workflow string, config []byte) ([]byte, error)
}

var workflowSpecFactories = map[WorkflowSpecType]WorkflowSpecFactory{
	YamlSpec:        YAMLSpecFactory{},
	WASMFile:        WasmFileSpecFactory{},
	DefaultSpecType: YAMLSpecFactory{},
}
