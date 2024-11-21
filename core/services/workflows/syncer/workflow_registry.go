package syncer

import (
	"context"
	_ "embed"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const (
	// Compute Fetch Workflow
	workflowID    = "924eef66516e5387b6e8ab8cc544685dfe50dfc837886f22beecebced5063968"
	workflowOwner = "00000000000000000000000000000000000000ab"
	workflowName  = "trueusdpor"

	// Chain Read Workflow
	workflow2ID    = "00000066516e5387b6e8ab8cc544685dfe50dfc837886f22beecebced5063968"
	workflow2Owner = "00000000000000000000000000000000000000ab"
	workflow2Name  = "ethseppor"
)

var (
	// Compute Fetch Workflow
	//go:embed config.yaml
	config []byte

	//go:embed workflow.wasm.br
	workflow []byte

	// Chain Read Workflow
	//go:embed config.yaml
	config2 []byte

	//go:embed por-read-chain.wasm.br
	workflow2 []byte
)

type WorkflowRegistry struct {
	services.StateMachine
	wg          sync.WaitGroup
	Logger      logger.Logger
	Registry    core.CapabilitiesRegistry
	Store       store.Store
	DS          sqlutil.DataSource
	subServices []job.ServiceCtx
}

func (w *WorkflowRegistry) Start(ctx context.Context) error {
	go func() {
		timeout := time.After(5 * time.Minute)
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-timeout:
				w.Logger.Info("timed out setting up hardcoded workflow")
				return
			case <-ticker.C:
				success1 := w.trySetup(workflowID, workflowName, workflowOwner, workflow, config)
				success2 := w.trySetup(workflow2ID, workflow2Name, workflow2Owner, workflow2, config2)
				if success1 && success2 {
					return
				}
			}
		}
	}()
	return nil
}

func (w *WorkflowRegistry) trySetup(id, name, owner string, binary, config []byte) bool {
	ctx := context.Background()
	w.Logger.Info("starting hardcoded workflow...")

	// HACK: don't load the workflow if we aren't a workflow node.
	_, err := w.Registry.Get(ctx, "offchain_reporting@1.0.0")
	if err != nil {
		w.Logger.Info("not a workflow node, skipping hardcoded workflow")
		return false
	}

	jb := job.WorkflowSpec{
		Workflow:      "a string",
		Config:        "a config",
		WorkflowID:    id,
		WorkflowName:  name,
		WorkflowOwner: owner,
	}
	sql := `INSERT INTO workflow_specs (workflow, workflow_id, workflow_owner, workflow_name, created_at, updated_at, spec_type, config)
	VALUES (:workflow, :workflow_id, :workflow_owner, :workflow_name, NOW(), NOW(), :spec_type, :config)
	RETURNING id;`
	_, err = w.DS.NamedExecContext(ctx, sql, jb)
	if err != nil {
		w.Logger.Info("failed to create entry: %w", err)
	}

	moduleConfig := &host.ModuleConfig{Logger: logger.NullLogger}
	spec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
	if err != nil {
		w.Logger.Errorf("failed to get workflow spec", err)
		return false
	}

	cfg := workflows.Config{
		Lggr:           w.Logger,
		Workflow:       *spec,
		WorkflowID:     id,
		WorkflowOwner:  owner,
		WorkflowName:   name,
		Registry:       w.Registry,
		Store:          w.Store,
		Config:         config,
		Binary:         binary,
		SecretsFetcher: w,
	}
	engine, err := workflows.NewEngine(ctx, cfg)
	if err != nil {
		w.Logger.Errorf("failed to create engine: %w", err)
		return false
	}
	err = engine.Start(ctx)
	if err != nil {
		w.Logger.Errorf("failed to start hardcoded workflow: %w", err)
		return false
	}
	w.subServices = []job.ServiceCtx{engine}
	return true
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
