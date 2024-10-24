package syncer

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type WorkflowRegistry struct {
	services.StateMachine
}

func (w *WorkflowRegistry) Start(ctx context.Context) error {
	return nil
}

func (w *WorkflowRegistry) Close() error {
	return nil
}

func (w *WorkflowRegistry) Ready() error {
	return nil
}

func (w *WorkflowRegistry) HealthReport() map[string]error {
	return nil
}

func (w *WorkflowRegistry) Name() string {
	return "WorkflowRegistrySyncer"
}

func (w *WorkflowRegistry) SecretsFor(workflowOwner, workflowName string) (map[string]string, error) {
	// TODO: actually get this from the right place.
	return map[string]string{}, nil
}

func NewWorkflowRegistry() *WorkflowRegistry {
	return &WorkflowRegistry{}
}
