package job

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type YAMLSpecFactory struct{}

var _ SDKWorkflowSpecFactory = (*YAMLSpecFactory)(nil)

func (y YAMLSpecFactory) Spec(_ context.Context, rawSpec, _ []byte) (sdk.WorkflowSpec, error) {
	return workflows.ParseWorkflowSpecYaml(string(rawSpec))
}

func (y YAMLSpecFactory) RawSpec(_ context.Context, wf string) ([]byte, error) {
	return []byte(wf), nil
}
