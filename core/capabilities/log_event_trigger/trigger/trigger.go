package trigger

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const ID = "cron-trigger@1.0.0"

const defaultSendChannelBufferSize = 1000

var logEventTriggerInfo = capabilities.MustNewCapabilityInfo(
	ID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that listens for specific contract log events and starts a workflow run.",
)

// Log Event Trigger Capability Config
type Config struct {
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

type logEventTrigger struct {
	ch chan<- capabilities.TriggerResponse
}

// Log Event Trigger Capabilities Manager
// Manages different log event triggers using an underlying triggerStore
type LogEventTriggerManager struct {
	capabilities.CapabilityInfo
	capabilities.Validator[Config, struct{}, capabilities.TriggerResponse]
	lggr     logger.Logger
	triggers CapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]
}

type Params struct {
	Logger logger.Logger
}

var _ capabilities.TriggerCapability = (*LogEventTriggerManager)(nil)
var _ services.Service = &LogEventTriggerManager{}

// Creates a new Cron Trigger Service.
// Scheduling will commence on calling .Start()
func New(p Params) *LogEventTriggerManager {
	l := logger.Named(p.Logger, "Log Event Trigger Capability Service")

	logEventStore := NewCapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]()

	return &LogEventTriggerManager{
		CapabilityInfo: logEventTriggerInfo,
		lggr:           l,
		triggers:       logEventStore,
	}
}

// Register a new trigger
// Can register triggers before the service is actively scheduling
func (s *LogEventTriggerManager) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if req.Config == nil {
		return nil, errors.New("config is required to register a log event trigger")
	}
	_, err := s.ValidateConfig(req.Config)
	if err != nil {
		return nil, err
	}
	respCh, err := s.triggers.InsertIfNotExists(req.TriggerID, func() (logEventTrigger, chan capabilities.TriggerResponse) {
		callbackCh := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
		return logEventTrigger{
			ch: callbackCh,
		}, callbackCh
	})
	if err != nil {
		return nil, fmt.Errorf("log_event_trigger %v", err)
	}
	return respCh, nil
}

func (s *LogEventTriggerManager) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	return nil
}

// Start the service.
func (s *LogEventTriggerManager) Start(ctx context.Context) error {
	return nil
}

// Close stops the Service.
// After this call the Service cannot be started again,
// The service will need to be re-built to start scheduling again.
func (s *LogEventTriggerManager) Close() error {
	return nil
}

func (s *LogEventTriggerManager) Ready() error {
	return nil
}

func (s *LogEventTriggerManager) HealthReport() map[string]error {
	return map[string]error{s.Name(): nil}
}

func (s *LogEventTriggerManager) Name() string {
	return "Service"
}
