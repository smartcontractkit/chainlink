package syncer

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/secrets"
)

const name = "WorkflowRegistrySyncer"

type Syncer interface {
	Sync(ctx context.Context, isInitialSync bool) error
}

type WorkflowRegistrySyncer interface {
	services.Service
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

type workflowRegistry struct {
	services.StateMachine
	lggr         logger.Logger
	orm          WorkflowRegistryDS
	secretSyncer secrets.Syncer
}

func (w *workflowRegistry) Start(ctx context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		return w.secretSyncer.Start(ctx)
	})
}

func (w *workflowRegistry) Close() error {
	return w.StopOnce(w.Name(), func() error {
		return w.secretSyncer.Close()
	})
}

func (w *workflowRegistry) Ready() error {
	return nil
}

func (w *workflowRegistry) HealthReport() map[string]error {
	return nil
}

func (w *workflowRegistry) Name() string {
	return name
}

func (w *workflowRegistry) SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error) {
	return w.secretSyncer.SecretsFor(ctx, workflowOwner, workflowName)
}

func NewWorkflowRegistry(
	lggr logger.Logger,
	orm WorkflowRegistryDS,
	secretSyncer secrets.Syncer,
	opts ...func(*workflowRegistry),
) *workflowRegistry {
	wr := &workflowRegistry{
		lggr:         lggr.Named(name),
		orm:          orm,
		secretSyncer: secretSyncer,
	}
	for _, opt := range opts {
		opt(wr)
	}
	return wr
}
