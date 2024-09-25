package job

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var ErrInvalidWorkflowType = errors.New("invalid workflow type")

type WorkflowSpecFactory interface {
	Spec(ctx context.Context, lggr logger.Logger, workflow string, config []byte) (sdk.WorkflowSpec, string, error)
}

var workflowSpecFactories = map[WorkflowSpecType]WorkflowSpecFactory{
	YamlSpec:        YAMLSpecFactory{},
	WASMFile:        WasmFileSpecFactory{},
	DefaultSpecType: YAMLSpecFactory{},
}
