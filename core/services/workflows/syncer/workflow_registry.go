package syncer

import (
	"context"
	_ "embed"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const (
	workflowID    = ""
	workflowOwner = ""
	workflowName  = ""
)

var (
	//go:embed config.yaml
	config []byte

	//go:embed workflow.wasm
	workflow []byte
)

type WorkflowRegistry struct {
	services.StateMachine
	wg          sync.WaitGroup
	Logger      logger.Logger
	Registry    core.CapabilitiesRegistry
	Store       store.Store
	subServices []job.ServiceCtx
}

func (w *WorkflowRegistry) Start(ctx context.Context) error {
	w.wg.Add(1)
	go func() {
		w.Logger.Info("starting hardcoded workflow...")

		moduleConfig := &host.ModuleConfig{Logger: logger.NullLogger, IsUncompressed: true}
		spec, err := host.GetWorkflowSpec(ctx, moduleConfig, workflow, config)
		if err != nil {
			w.Logger.Errorf("failed to get workflow spec", err)
		}

		cfg := workflows.Config{
			Lggr:           w.Logger,
			Workflow:       *spec,
			WorkflowID:     workflowID,
			WorkflowOwner:  workflowOwner,
			WorkflowName:   workflowName,
			Registry:       w.Registry,
			Store:          w.Store,
			Config:         config,
			Binary:         workflow,
			SecretsFetcher: w,
		}
		engine, err := workflows.NewEngine(ctx, cfg)
		if err != nil {
			w.Logger.Errorf("failed to create engine: %w", err)
		}
		err = engine.Start(ctx)
		if err != nil {
			w.Logger.Errorf("failed to start hardcoded workflow: %w", err)
		}
		w.subServices = []job.ServiceCtx{engine}
	}()
	return nil
}

func (w *WorkflowRegistry) Close() error {
	for _, s := range w.subServices {
		err := s.Close()
		if err != nil {
			w.Logger.Errorf("could not close hardcoded engine: %w", err)
		}
	}

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
