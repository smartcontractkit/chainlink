package logevent

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
)

const ID = "log-event-trigger-%s-%s@1.0.0"

const defaultSendChannelBufferSize = 1000

// Log Event Trigger Capability Input
type Input struct {
}

// Log Event Trigger Capabilities Manager
// Manages different log event triggers using an underlying triggerStore
type TriggerService struct {
	services.StateMachine
	capabilities.CapabilityInfo
	capabilities.Validator[logeventcap.Config, Input, capabilities.TriggerResponse]
	lggr           logger.Logger
	triggers       CapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]
	relayer        core.Relayer
	logEventConfig Config
	stopCh         services.StopChan
}

// Common capability level config across all workflows
type Config struct {
	ChainID        string `json:"chainId"`
	Network        string `json:"network"`
	LookbackBlocks uint64 `json:"lookbakBlocks"`
	PollPeriod     uint32 `json:"pollPeriod"`
	QueryCount     uint64 `json:"queryCount"`
}

func (config Config) Version(capabilityVersion string) string {
	return fmt.Sprintf(capabilityVersion, config.Network, config.ChainID)
}

var _ capabilities.TriggerCapability = (*TriggerService)(nil)
var _ services.Service = &TriggerService{}

// Creates a new Cron Trigger Service.
// Scheduling will commence on calling .Start()
func NewTriggerService(ctx context.Context,
	lggr logger.Logger,
	relayer core.Relayer,
	logEventConfig Config) (*TriggerService, error) {
	l := logger.Named(lggr, "LogEventTriggerCapabilityService")

	logEventStore := NewCapabilitiesStore[logEventTrigger, capabilities.TriggerResponse]()

	s := &TriggerService{
		lggr:           l,
		triggers:       logEventStore,
		relayer:        relayer,
		logEventConfig: logEventConfig,
		stopCh:         make(services.StopChan),
	}
	var err error
	s.CapabilityInfo, err = s.Info(ctx)
	if err != nil {
		return s, err
	}
	s.Validator = capabilities.NewValidator[logeventcap.Config, Input, capabilities.TriggerResponse](capabilities.ValidatorArgs{Info: s.CapabilityInfo})
	return s, nil
}

func (s *TriggerService) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(
		s.logEventConfig.Version(ID),
		capabilities.CapabilityTypeTrigger,
		"A trigger that listens for specific contract log events and starts a workflow run.",
	)
}

// Register a new trigger
// Can register triggers before the service is actively scheduling
func (s *TriggerService) RegisterTrigger(ctx context.Context,
	req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if req.Config == nil {
		return nil, errors.New("config is required to register a log event trigger")
	}
	reqConfig, err := s.ValidateConfig(req.Config)
	if err != nil {
		return nil, err
	}
	// Add log event trigger with Contract details to CapabilitiesStore
	var respCh chan capabilities.TriggerResponse
	ok := s.IfNotStopped(func() {
		respCh, err = s.triggers.InsertIfNotExists(req.TriggerID, func() (*logEventTrigger, chan capabilities.TriggerResponse, error) {
			l, ch, tErr := newLogEventTrigger(ctx, s.lggr, req.Metadata.WorkflowID, reqConfig, s.logEventConfig, s.relayer)
			if tErr != nil {
				return l, ch, tErr
			}
			tErr = l.Start(ctx)
			return l, ch, tErr
		})
	})
	if !ok {
		return nil, errors.New("cannot create new trigger since LogEventTriggerCapabilityService has been stopped")
	}
	if err != nil {
		return nil, fmt.Errorf("create new trigger failed %w", err)
	}
	s.lggr.Infow("RegisterTrigger", "triggerId", req.TriggerID, "WorkflowID", req.Metadata.WorkflowID)
	return respCh, nil
}

func (s *TriggerService) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	trigger, ok := s.triggers.Read(req.TriggerID)
	if !ok {
		return fmt.Errorf("triggerId %s not found", req.TriggerID)
	}
	// Close callback channel and stop log event trigger listener
	err := trigger.Close()
	if err != nil {
		return fmt.Errorf("error closing trigger %s (chainID %s): %w", req.TriggerID, s.logEventConfig.ChainID, err)
	}
	// Remove from triggers context
	s.triggers.Delete(req.TriggerID)
	s.lggr.Infow("UnregisterTrigger", "triggerId", req.TriggerID, "WorkflowID", req.Metadata.WorkflowID)
	return nil
}

// Start the service.
func (s *TriggerService) Start(ctx context.Context) error {
	return s.StartOnce("LogEventTriggerCapabilityService", func() error {
		s.lggr.Info("Starting LogEventTriggerCapabilityService")
		return nil
	})
}

// Close stops the Service.
// After this call the Service cannot be started again,
// The service will need to be re-built to start scheduling again.
func (s *TriggerService) Close() error {
	return s.StopOnce("LogEventTriggerCapabilityService", func() error {
		s.lggr.Infow("Stopping LogEventTriggerCapabilityService")
		triggers := s.triggers.ReadAll()
		return services.MultiCloser(triggers).Close()
	})
}

func (s *TriggerService) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

func (s *TriggerService) Name() string {
	return s.lggr.Name()
}
