package logevent

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const ID = "log-event-trigger-%s-%d@1.0.0"

const defaultSendChannelBufferSize = 1000

// Log Event Trigger Capability Input
type Input struct {
}

// Log Event Trigger Capability Payload
type Payload struct {
	// Time that Log Event Trigger's task execution occurred (RFC3339Nano formatted)
	ActualExecutionTime string
}

// Log Event Trigger Capability Response
type Response struct {
	capabilities.TriggerEvent
	Metadata struct{}
	Payload  Payload
}

// Log Event Trigger Capabilities Manager
// Manages different log event triggers using an underlying triggerStore
type LogEventTriggerService struct {
	capabilities.CapabilityInfo
	capabilities.Validator[RequestConfig, Input, capabilities.TriggerResponse]
	lggr           logger.Logger
	triggers       CapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]
	relayer        core.Relayer
	logEventConfig LogEventConfig
}

// Common capability level config across all workflows
type LogEventConfig struct {
	ChainId        uint64 `json:"chainId"`
	Network        string `json:"network"`
	LookbackBlocks uint64 `json:"lookbakBlocks"`
	PollPeriod     uint64 `json:"pollPeriod"`
}

func (config LogEventConfig) Version(capabilityVersion string) string {
	return fmt.Sprintf(capabilityVersion, config.Network, config.ChainId)
}

type Params struct {
	Logger         logger.Logger
	Relayer        core.Relayer
	RelayerSet     core.RelayerSet
	LogEventConfig LogEventConfig
}

var _ capabilities.TriggerCapability = (*LogEventTriggerService)(nil)
var _ services.Service = &LogEventTriggerService{}

// Creates a new Cron Trigger Service.
// Scheduling will commence on calling .Start()
func NewLogEventTriggerService(p Params) *LogEventTriggerService {
	l := logger.Named(p.Logger, "LogEventTriggerCapabilityService: ")

	logEventStore := NewCapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]()

	s := &LogEventTriggerService{
		lggr:           l,
		triggers:       logEventStore,
		relayer:        p.Relayer,
		logEventConfig: p.LogEventConfig,
	}
	s.CapabilityInfo, _ = s.Info(context.Background())
	s.Validator = capabilities.NewValidator[RequestConfig, Input, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: s.CapabilityInfo})
	return s
}

func (s *LogEventTriggerService) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		s.logEventConfig.Version(ID),
		capabilities.CapabilityTypeTrigger,
		"A trigger that listens for specific contract log events and starts a workflow run.",
	)
}

// Register a new trigger
// Can register triggers before the service is actively scheduling
func (s *LogEventTriggerService) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if req.Config == nil {
		return nil, errors.New("config is required to register a log event trigger")
	}
	reqConfig, err := s.ValidateConfig(req.Config)
	if err != nil {
		return nil, err
	}
	// Add log event trigger with Contract details to CapabilitiesStore
	respCh, err := s.triggers.InsertIfNotExists(req.TriggerID, func() (*logEventTrigger, chan capabilities.TriggerResponse, error) {
		return newLogEventTrigger(ctx, reqConfig, s.logEventConfig, s.relayer)
	})
	if err != nil {
		return nil, fmt.Errorf("LogEventTrigger %v", err)
	}
	s.lggr.Debugw("RegisterTrigger", "triggerId", req.TriggerID)
	return respCh, nil
}

func (s *LogEventTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	trigger, ok := s.triggers.Read(req.TriggerID)
	if !ok {
		return fmt.Errorf("triggerId %s not found", req.TriggerID)
	}
	// Close callback channel and stop log event trigger listener
	trigger.Stop()
	// Remove from triggers context
	s.triggers.Delete(req.TriggerID)
	s.lggr.Debugw("UnregisterTrigger", "triggerId", req.TriggerID)
	return nil
}

// Start the service.
func (s *LogEventTriggerService) Start(ctx context.Context) error {
	if s.relayer == nil {
		return errors.New("service has shutdown, it must be built again to restart")
	}

	return nil
}

// Close stops the Service.
// After this call the Service cannot be started again,
// The service will need to be re-built to start scheduling again.
func (s *LogEventTriggerService) Close() error {
	return nil
}

func (s *LogEventTriggerService) Ready() error {
	return nil
}

func (s *LogEventTriggerService) HealthReport() map[string]error {
	return map[string]error{s.Name(): nil}
}

func (s *LogEventTriggerService) Name() string {
	return "Service"
}
