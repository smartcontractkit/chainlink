package direct

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func CreateWorkflow() (*workflow.Spec, error) {
	trigger := NewHttpTrigger("http", "httpTrigger")
	root, builder := workflow.NewWorkflowBuilder[*HttpTrigger](trigger)

	return root.Build()
}
