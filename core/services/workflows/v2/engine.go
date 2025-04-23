package v2

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/internal"
)

type Engine struct {
	services.Service
	srvcEng *services.Engine

	cfg       EngineConfig
	localNode capabilities.Node
}

type LifecycleHooks struct {
	OnInitialized       func(err error)
	OnExecutionFinished func(executionID string)
	OnRateLimited       func(executionID string)
}

func NewEngine(ctx context.Context, cfg EngineConfig) (*Engine, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	engine := &Engine{
		cfg: cfg,
	}
	engine.Service, engine.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: engine.start,
		Close: engine.close,
	}.NewServiceEngine(cfg.Lggr.Named("WorkflowEngine").With("workflowID", cfg.WorkflowID))
	return engine, nil
}

func (e *Engine) start(_ context.Context) error {
	// TODO(CAPPL-733): apply global workflow limits

	e.cfg.Module.Start()
	e.srvcEng.Go(e.init)
	return nil
}

func (e *Engine) init(ctx context.Context) {
	// retrieve info about the current node we are running on
	retryErr := internal.RunWithRetries(ctx, e.cfg.Lggr, time.Millisecond*time.Duration(e.cfg.Limits.CapRegistryAccessRetryIntervalMs), e.cfg.Limits.MaxCapRegistryAccessRetries, func() error {
		// retry until the underlying peerWrapper service is ready
		node, err := e.cfg.CapRegistry.LocalNode(ctx)
		if err != nil {
			return fmt.Errorf("failed to get donInfo: %w", err)
		}
		e.localNode = node
		return nil
	})

	if retryErr != nil {
		e.cfg.Lggr.Errorw("Workflow Engine initialization failed", "workflowID", e.cfg.WorkflowID, "err", retryErr)
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(retryErr)
		return
	}

	err := e.runTriggerSubscriptionPhase(ctx)
	if err != nil {
		e.cfg.Lggr.Errorw("Workflow Engine initialization failed", "workflowID", e.cfg.WorkflowID, "err", err)
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	e.cfg.Lggr.Infow("Workflow Engine initialized", "workflowID", e.cfg.WorkflowID)
	e.cfg.Hooks.OnInitialized(nil)
}

func (e *Engine) runTriggerSubscriptionPhase(_ context.Context) error {
	// TODO (CAPPL-734): Subscription Phase:
	//   - call into WASM to get triggers
	//   - register to triggers
	//   - start goroutines that wait for events from each trigger
	return nil
}

func (e *Engine) close() error {
	e.cfg.Module.Close()
	return nil
}
