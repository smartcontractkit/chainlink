package job

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type YAMLSpecFactory struct{}

var _ WorkflowSpecFactory = (*YAMLSpecFactory)(nil)

func (y YAMLSpecFactory) Spec(ctx context.Context, lggr logger.Logger, workflow string, config []byte) (sdk.WorkflowSpec, string, error) {
	spec, err := workflows.ParseWorkflowSpecYaml(workflow)
	return spec, fmt.Sprintf("%x", sha256.Sum256([]byte(workflow))), err
}
