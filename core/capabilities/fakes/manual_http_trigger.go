package fakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ services.Service = (*ManualHttpTriggerService)(nil)
var _ httpserver.HTTPCapability = (*ManualHttpTriggerService)(nil)

const HTTPTriggerServiceName = "HttpTriggerService"
const HTTPTriggerID = "http-trigger@1.0.0"

var manualHttpTriggerInfo = capabilities.MustNewCapabilityInfo(
	HTTPTriggerID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses an HTTP request to run periodically at fixed times, dates, or intervals.",
)

type ManualHttpTriggerService struct {
	capabilities.CapabilityInfo
	lggr       logger.Logger
	callbackCh chan capabilities.TriggerAndId[*httptypedapi.Payload]
}

func NewManualManualHttpTriggerService(parentLggr logger.Logger) *ManualHttpTriggerService {
	lggr := logger.Named(parentLggr, "HttpTriggerService")

	return &ManualHttpTriggerService{
		CapabilityInfo: manualHttpTriggerInfo,
		lggr:           lggr,
		callbackCh:     make(chan capabilities.TriggerAndId[*httptypedapi.Payload]),
	}
}

// HTTPCapability interface methods
func (f *ManualHttpTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) (<-chan capabilities.TriggerAndId[*httptypedapi.Payload], error) {
	return f.callbackCh, nil
}

func (f *ManualHttpTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) error {
	return nil
}

func (f *ManualHttpTriggerService) Initialise(ctx context.Context, config string,
	_ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector) error {
	f.lggr.Debugf("Initialising %s", HTTPTriggerServiceName)
	return f.Start(ctx)
}

// ManualTriggerCapability interface method
func (f *ManualHttpTriggerService) ManualTrigger(ctx context.Context, payload *httptypedapi.Payload) error {
	// Run in a goroutine to avoid blocking
	go func() {
		// Send the trigger response
		f.callbackCh <- createManualHttpTriggerResponse(payload)
	}()

	return nil
}

func createManualHttpTriggerResponse(payload *httptypedapi.Payload) capabilities.TriggerAndId[*httptypedapi.Payload] {
	return capabilities.TriggerAndId[*httptypedapi.Payload]{
		Trigger: payload,
		Id:      "manual-http-trigger-id",
	}
}

// Service interface methods
func (f *ManualHttpTriggerService) Start(ctx context.Context) error {
	f.lggr.Info("Starting ManualManualHttpTriggerService")
	return nil
}

func (f *ManualHttpTriggerService) Close() error {
	f.lggr.Info("Closing ManualManualHttpTriggerService")
	return nil
}

func (f *ManualHttpTriggerService) Ready() error {
	return nil
}

func (f *ManualHttpTriggerService) HealthReport() map[string]error {
	return map[string]error{f.Name(): nil}
}

func (f *ManualHttpTriggerService) Name() string {
	return f.lggr.Name()
}

func (f *ManualHttpTriggerService) Description() string {
	return "Manual HTTP Trigger Service"
}
