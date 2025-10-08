package fakes

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	httptypedapi "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

var _ services.Service = (*ManualHTTPTriggerService)(nil)
var _ httpserver.HTTPCapability = (*ManualHTTPTriggerService)(nil)

const HTTPTriggerServiceName = "HttpTriggerService"
const HTTPTriggerID = "http-trigger@1.0.0-alpha"

var manualHTTPTriggerInfo = capabilities.MustNewCapabilityInfo(
	HTTPTriggerID,
	capabilities.CapabilityTypeTrigger,
	"A trigger that uses an HTTP request to run periodically at fixed times, dates, or intervals.",
)

type ManualHTTPTriggerService struct {
	capabilities.CapabilityInfo
	lggr        logger.Logger
	callbackCh  map[string]chan capabilities.TriggerAndId[*httptypedapi.Payload]
	workflowIDs map[string]string // triggerID -> workflowID mapping
}

func NewManualHTTPTriggerService(parentLggr logger.Logger) *ManualHTTPTriggerService {
	lggr := logger.Named(parentLggr, "HTTPTriggerService")

	return &ManualHTTPTriggerService{
		CapabilityInfo: manualHTTPTriggerInfo,
		lggr:           lggr,
		callbackCh:     make(map[string]chan capabilities.TriggerAndId[*httptypedapi.Payload]),
		workflowIDs:    make(map[string]string),
	}
}

// HTTPCapability interface methods
func (f *ManualHTTPTriggerService) RegisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) (<-chan capabilities.TriggerAndId[*httptypedapi.Payload], error) {
	f.callbackCh[triggerID] = make(chan capabilities.TriggerAndId[*httptypedapi.Payload])
	f.workflowIDs[triggerID] = metadata.WorkflowID
	return f.callbackCh[triggerID], nil
}

func (f *ManualHTTPTriggerService) UnregisterTrigger(ctx context.Context, triggerID string, metadata capabilities.RequestMetadata, input *httptypedapi.Config) error {
	return nil
}

func (f *ManualHTTPTriggerService) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	f.lggr.Debugf("Initialising %s", HTTPTriggerServiceName)
	return f.Start(ctx)
}

// ManualTriggerCapability interface method
func (f *ManualHTTPTriggerService) ManualTrigger(ctx context.Context, triggerID string, payload *httptypedapi.Payload) error {
	triggerEvent := f.createManualTriggerEvent(payload)

	triggerEventID := triggerEvent.Id

	workflowID, exists := f.workflowIDs[triggerID]
	if !exists {
		f.lggr.Errorw("workflowID not found for triggerID", "triggerID", triggerID)
		workflowID = "unknownWorkflow"
	}

	workflowExecutionID, err := events.GenerateExecutionID(workflowID, triggerEventID)
	if err != nil {
		f.lggr.Errorw("failed to generate execution ID", "err", err)
		workflowExecutionID = ""
	}
	err = events.EmitTriggerExecutionStarted(ctx, map[string]string{}, triggerEventID, workflowExecutionID)
	if err != nil {
		f.lggr.Errorw("failed to emit trigger execution started event", "err", err)
	}

	// Run in a goroutine to avoid blocking
	go func() {
		select {
		case f.callbackCh[triggerID] <- f.createManualTriggerEvent(payload):
			// Successfully sent trigger response
		case <-ctx.Done():
			// Context cancelled, cleanup goroutine
			f.lggr.Debug("ManualTrigger goroutine cancelled due to context cancellation")
		}
	}()

	return nil
}

func (f *ManualHTTPTriggerService) createManualTriggerEvent(payload *httptypedapi.Payload) capabilities.TriggerAndId[*httptypedapi.Payload] {
	return capabilities.TriggerAndId[*httptypedapi.Payload]{
		Trigger: payload,
		Id:      "manual-http-trigger-id",
	}
}

// Service interface methods
func (f *ManualHTTPTriggerService) Start(ctx context.Context) error {
	f.lggr.Debug("Starting HTTP Trigger Capability")
	return nil
}

func (f *ManualHTTPTriggerService) Close() error {
	f.lggr.Debug("Closing HTTP Trigger Capability")
	return nil
}

func (f *ManualHTTPTriggerService) HealthReport() map[string]error {
	return map[string]error{f.Name(): nil}
}

func (f *ManualHTTPTriggerService) Name() string {
	return f.lggr.Name()
}

func (f *ManualHTTPTriggerService) Description() string {
	return "Manual HTTP Trigger Service"
}

func (f *ManualHTTPTriggerService) Ready() error {
	return nil
}
