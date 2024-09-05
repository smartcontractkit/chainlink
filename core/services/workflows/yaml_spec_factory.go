package workflows

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type YAMLSpecFactory struct{}

var _ SDKWorkflowSpecFactory = (*YAMLSpecFactory)(nil)

func (y YAMLSpecFactory) GetSpec(rawSpec, _ []byte) (sdk.WorkflowSpec, error) {
	return workflows.ParseWorkflowSpecYaml(string(rawSpec))
}

func (y YAMLSpecFactory) GetRawSpec(wf string) ([]byte, error) {
	return []byte(wf), nil
}
