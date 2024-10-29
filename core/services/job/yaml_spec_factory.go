package job

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type YAMLSpecFactory struct{}

var _ WorkflowSpecFactory = (*YAMLSpecFactory)(nil)

func (y YAMLSpecFactory) Spec(_ context.Context, workflow, _ string) (sdk.WorkflowSpec, []byte, string, error) {
	spec, err := workflows.ParseWorkflowSpecYaml(workflow)
	return spec, []byte(workflow), fmt.Sprintf("%x", sha256.Sum256([]byte(workflow))), err
}

func (y YAMLSpecFactory) RawSpec(_ context.Context, workflow, _ string) ([]byte, error) {
	return []byte(workflow), nil
}

func (y YAMLSpecFactory) Config(_ context.Context, config string) ([]byte, error) {
	return []byte(config), nil
}
