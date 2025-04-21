package v2

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type EngineConfig struct {
	Lggr       logger.Logger
	WorkflowID string
	Module     host.ModuleV2
}

type Engine struct {
	services.Service
	srvcEng *services.Engine

	cfg EngineConfig
}

func NewEngine(ctx context.Context, cfg EngineConfig) (*Engine, error) {
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
	return nil
}

func (e *Engine) close() error {
	return nil
}
