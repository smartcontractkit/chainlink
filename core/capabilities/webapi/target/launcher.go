package target

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ job.ServiceCtx = &Launcher{}

// Config is the configuration for the Target capability
// Note that workflow executions have their own internal timeouts and retries set by the user
// that are separate from this configuration
type Config struct {
	TimeoutMs  uint32 // Timeout in milliseconds from the capability to the gateway
	RetryCount uint8  // Number of retries from the capability to the gateway
}

type Launcher struct {
	gc         connector.GatewayConnector
	targetCfg  Config
	registry   core.CapabilitiesRegistry
	lggr       logger.Logger
	capability *Capability
	handler    *Handler
}

func NewLauncher(gc connector.GatewayConnector, cfg string, registry core.CapabilitiesRegistry, lggr logger.Logger) (job.ServiceCtx, error) {
	var targetCfg Config
	if len(cfg) == 0 {
		return nil, errors.New("config is empty")
	}
	// TODO: is config JSON for standard capabilities?
	err := json.Unmarshal([]byte(cfg), &targetCfg)
	if err != nil {
		return nil, err
	}
	lggr = logger.Named(lggr, "WebAPITarget")
	// response channels and the mutex are shared between capability and handler
	responseChs := make(map[string]chan *api.Message)
	responseChsMu := &sync.Mutex{}
	capability, err := NewCapability(targetCfg, registry, gc, lggr, responseChs, responseChsMu)
	if err != nil {
		return nil, err
	}
	handler, err := NewHandler(gc, responseChs, responseChsMu, lggr)
	if err != nil {
		return nil, err
	}

	return &Launcher{
		gc:         gc,
		targetCfg:  targetCfg,
		registry:   registry,
		lggr:       lggr,
		capability: capability,
		handler:    handler,
	}, nil
}

func (l *Launcher) Start(ctx context.Context) error {
	if err := l.capability.Start(ctx); err != nil {
		return err
	}
	if err := l.handler.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (l *Launcher) Close() error {
	if err := l.capability.Close(); err != nil {
		return err
	}
	if err := l.handler.Close(); err != nil {
		return err
	}
	return nil
}
